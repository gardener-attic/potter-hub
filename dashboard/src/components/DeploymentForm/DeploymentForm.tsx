import { RouterAction } from "connected-react-router";
import * as React from "react";

import { JSONSchema4 } from "json-schema";
import { IAppRepository, IChartState, IChartVersion } from "../../shared/types";
import { getInfoAnnotation } from "../../shared/utils";
import DeploymentFormBody from "../DeploymentFormBody/DeploymentFormBody";
import { ErrorSelector } from "../ErrorAlert";
import InfoBox from "../InfoBox";
import LoadingWrapper from "../LoadingWrapper";

import "react-tabs/style/react-tabs.css";

export interface IDeploymentFormProps {
  kubeappsNamespace: string;
  chartID: string;
  chartVersion: string;
  error: Error | undefined;
  selected: IChartState["selected"];
  deployChart: (
    version: IChartVersion,
    releaseName: string,
    namespace: string,
    values?: string,
    schema?: JSONSchema4,
  ) => Promise<boolean>;
  push: (location: string) => RouterAction;
  fetchChartVersions: (id: string) => void;
  getChartVersion: (id: string, chartVersion: string) => void;
  namespace: string;
  repos: IAppRepository[];
  fetchRepos: () => void;
  repoErrors: {
    create?: Error;
    delete?: Error;
    fetch?: Error;
    update?: Error;
  };
  generateReleaseName: () => string
}

export interface IDeploymentFormState {
  isDeploying: boolean;
  releaseName: string;
  // Name of the release that was submitted for creation
  // This is different than releaseName since it is also used in the error banner
  // and we do not want to use releaseName since it is controller by the form field.
  latestSubmittedReleaseName: string;
  appValues: string;
  valuesModified: boolean;
}

class DeploymentForm extends React.Component<IDeploymentFormProps, IDeploymentFormState> {
  public state: IDeploymentFormState = {
    releaseName: this.props.generateReleaseName(),
    appValues: this.props.selected.values || "",
    isDeploying: false,
    latestSubmittedReleaseName: "",
    valuesModified: false,
  };

  public componentDidMount() {
    this.props.fetchChartVersions(this.props.chartID);
    this.props.fetchRepos();
  }

  public componentDidUpdate(prevProps: IDeploymentFormProps) {
    if (prevProps.selected.version !== this.props.selected.version && !this.state.valuesModified) {
      this.setState({ appValues: this.props.selected.values || "" });
    }
  }

  public render() {
    const { namespace } = this.props;
    if (this.props.error) {
      return (
        <ErrorSelector
          error={this.props.error}
          namespace={namespace}
          action="create"
          resource={this.state.latestSubmittedReleaseName}
        />
      );
    }
    if (this.state.isDeploying) {
      return <LoadingWrapper />;
    }

    const repoName = this.props.chartID.split("/")[0];
    let repo = {} as IAppRepository;
    for (const r of this.props.repos) {
      if (r.metadata.name === repoName) {
        repo = r;
        break;
      }
    }

    return (
      <form className="container padding-b-bigger" onSubmit={this.handleDeploy}>
        <div className="row">
          {this.props.repoErrors.fetch && (
            <div className="col-8">
              <ErrorSelector
                error={this.props.repoErrors.fetch}
                action="get"
                resource="appRepositories"
              />
            </div>
          )}
          {!this.props.repoErrors.fetch && (
            <div className="col-8">
              <InfoBox msg={getInfoAnnotation(repo)} />
            </div>
          )}
          <div className="col-12">
            <h2>{this.props.chartID}</h2>
          </div>
          <div className="col-8">
            <div>
              <label htmlFor="releaseName">Name</label>
              <input
                id="releaseName"
                pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*"
                title="Use lower case alphanumeric characters, '-' or '.'"
                onChange={this.handleReleaseNameChange}
                value={this.state.releaseName}
                required={true}
              />
            </div>
            <DeploymentFormBody
              chartID={this.props.chartID}
              chartVersion={this.props.chartVersion}
              namespace={this.props.namespace}
              selected={this.props.selected}
              push={this.props.push}
              getChartVersion={this.props.getChartVersion}
              setValues={this.handleValuesChange}
              appValues={this.state.appValues}
              setValuesModified={this.setValuesModified}
            />
          </div>
        </div>
      </form>
    );
  }

  public handleValuesChange = (value: string) => {
    this.setState({ appValues: value });
  };

  public setValuesModified = () => {
    this.setState({ valuesModified: true });
  };

  public handleDeploy = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const { selected, deployChart, push, namespace } = this.props;
    const { releaseName, appValues } = this.state;

    this.setState({ isDeploying: true, latestSubmittedReleaseName: releaseName });
    if (selected.version) {
      const deployed = await deployChart(
        selected.version,
        releaseName,
        namespace,
        appValues,
        selected.schema,
      );
      this.setState({ isDeploying: false });
      if (deployed) {
        push(`/apps/ns/${namespace}/${releaseName}`);
      }
    }
  };

  public handleReleaseNameChange = (e: React.FormEvent<HTMLInputElement>) => {
    this.setState({ releaseName: e.currentTarget.value });
  };
}

export default DeploymentForm;
