import { push } from "connected-react-router";
import { connect } from "react-redux";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import { JSONSchema4 } from "json-schema";
import actions from "../../actions";
import DeploymentForm from "../../components/DeploymentForm";
import { getRandomName } from "../../shared/namegenerator";
import { IChartVersion, IStoreState } from "../../shared/types";

interface IRouteProps {
  match: {
    params: {
      repo: string;
      id: string;
      version: string;
    };
  };
}

function mapStateToProps(
  { apps, charts, config, namespace, repos }: IStoreState,
  { match: { params } }: IRouteProps,
) {
  return {
    chartID: `${params.repo}/${params.id}`,
    chartVersion: params.version,
    error: apps.error,
    kubeappsNamespace: config.namespace,
    namespace: namespace.current,
    selected: charts.selected,
    repos: repos.repos,
    repoErrors: repos.errors,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    deployChart: (
      version: IChartVersion,
      releaseName: string,
      namespace: string,
      values?: string,
      schema?: JSONSchema4,
    ) => dispatch(actions.apps.deployChart(version, releaseName, namespace, values, schema)),
    fetchChartVersions: (id: string) => dispatch(actions.charts.fetchChartVersions(id)),
    getChartVersion: (id: string, version: string) =>
      dispatch(actions.charts.getChartVersion(id, version)),
    push: (location: string) => dispatch(push(location)),
    fetchRepos: async () => {
      return dispatch(actions.repos.fetchRepos());
    },
    generateReleaseName: getRandomName
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(DeploymentForm);
