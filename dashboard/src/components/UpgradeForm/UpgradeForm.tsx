import { RouterAction } from "connected-react-router";
import * as jsonpatch from "fast-json-patch";
import { JSONSchema4 } from "json-schema";
import * as React from "react";
import * as YAML from "yaml";

import { deleteValue, setValue } from "../../shared/schema";
import { IChartState, IChartVersion } from "../../shared/types";
import DeploymentFormBody from "../DeploymentFormBody/DeploymentFormBody";
import { ErrorSelector } from "../ErrorAlert";
import LoadingWrapper from "../LoadingWrapper";

export interface IUpgradeFormProps {
  appCurrentVersion: string;
  appCurrentValues?: string;
  chartName: string;
  namespace: string;
  releaseName: string;
  repo: string;
  error: Error | undefined;
  selected: IChartState["selected"];
  deployed: IChartState["deployed"];
  upgradeApp: (
    version: IChartVersion,
    releaseName: string,
    namespace: string,
    values?: string,
    schema?: JSONSchema4,
  ) => Promise<boolean>;
  push: (location: string) => RouterAction;
  goBack: () => RouterAction;
  fetchChartVersions: (id: string) => Promise<IChartVersion[] | undefined>;
  getChartVersion: (id: string, chartVersion: string) => void;
}

interface IUpgradeFormState {
  appValues: string;
  valuesModified: boolean;
  isDeploying: boolean;
  modifications?: jsonpatch.Operation[];
  deployedValues?: string;
}

class UpgradeForm extends React.Component<IUpgradeFormProps, IUpgradeFormState> {
  public state: IUpgradeFormState = {
    appValues: this.props.appCurrentValues || "",
    isDeploying: false,
    valuesModified: false,
  };

  public componentDidMount() {
    const chartID = `${this.props.repo}/${this.props.chartName}`;
    this.props.fetchChartVersions(chartID);
  }

  public componentDidUpdate = (prevProps: IUpgradeFormProps) => {
    let modifications = this.state.modifications;
    // applyModifications is an expensive operatior, that's why it's only defined within
    // the componentDidUpdate scope
    const applyModifications = (mods: jsonpatch.Operation[], appValues: string) => {
      // And we add any possible change made to the original version
      if (mods.length) {
        mods.forEach(modification => {
          if (modification.op === "remove") {
            appValues = deleteValue(appValues, modification.path);
          } else {
            // Transform the modification as a ReplaceOperation to read its value
            const value = (modification as jsonpatch.ReplaceOperation<any>).value;
            appValues = setValue(appValues, modification.path, value);
          }
        });
      }
      return appValues;
    };

    if (this.props.deployed.values && !modifications) {
      // Calculate modifications from the default values
      const defaultValuesObj = YAML.parse(this.props.deployed.values);
      const deployedValuesObj = YAML.parse(this.props.appCurrentValues || "");
      modifications = jsonpatch.compare(defaultValuesObj, deployedValuesObj);
      const values = applyModifications(modifications, this.props.deployed.values);
      this.setState({ modifications });
      this.setState({ appValues: values, deployedValues: values });
    }

    if (prevProps.selected.version !== this.props.selected.version && !this.state.valuesModified) {
      // Apply modifications to the new selected version
      const appValues = modifications
        ? applyModifications(modifications, this.props.selected.values || "")
        : this.props.selected.values || "";
      this.setState({ appValues });
    }
  };

  public render() {
    const { namespace, releaseName, error, selected } = this.props;
    if (error) {
      return (
        <ErrorSelector error={error} namespace={namespace} action="update" resource={releaseName} />
      );
    }
    if (selected.versions.length === 0) {
      return <LoadingWrapper />;
    }
    const chartID = `${this.props.repo}/${this.props.chartName}`;
    return (
      <form className="container padding-b-bigger" onSubmit={this.handleDeploy}>
        <div className="row">
          <div className="col-12">
            <h2>{`${this.props.releaseName} (${chartID})`}</h2>
          </div>
          <div className="col-8">
            <DeploymentFormBody
              chartID={chartID}
              chartVersion={this.props.appCurrentVersion}
              deployedValues={this.state.deployedValues}
              namespace={this.props.namespace}
              releaseVersion={this.props.appCurrentVersion}
              selected={this.props.selected}
              push={this.props.push}
              goBack={this.props.goBack}
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

  public setValuesModified = () => {
    this.setState({ valuesModified: true });
  };

  public handleValuesChange = (value: string) => {
    this.setState({ appValues: value });
  };

  public handleDeploy = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const { selected, push, upgradeApp, releaseName, namespace } = this.props;
    const { appValues } = this.state;

    this.setState({ isDeploying: true });
    if (selected.version) {
      const deployed = await upgradeApp(
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
}

export default UpgradeForm;
