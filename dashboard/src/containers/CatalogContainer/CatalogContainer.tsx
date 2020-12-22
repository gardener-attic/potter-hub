import { push } from "connected-react-router";
import * as qs from "qs";
import { connect } from "react-redux";
import { RouteComponentProps } from "react-router";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";
import actions from "../../actions";
import Catalog from "../../components/Catalog";
import { IStoreState } from "../../shared/types";

function mapStateToProps(
  {
    charts,
    config: {
      featuredChartIds,
      generalRepoInfo,
      defaultRepo,
      staticCatalogInfo,
      urlParams: { targetClusterSecretName, targetClusterSecretNamespace } },
    repos,
  }: IStoreState,
  { match: { params }, location }: RouteComponentProps<{ repo: string }>,
) {
  return {
    appRepoState: repos,
    charts,
    featuredChartIds,
    filter: qs.parse(location.search, { ignoreQueryPrefix: true }).q || "",
    repo: params.repo,
    generalRepoInfo,
    defaultRepo,
    targetClusterSecretName,
    targetClusterSecretNamespace,
    staticCatalogInfo
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    fetchApprepositories: () => dispatch(actions.repos.fetchRepos()),
    fetchCharts: (repo: string) => dispatch(actions.charts.fetchCharts(repo)),
    pushSearchFilter: (filter: string) => dispatch(actions.shared.pushSearchFilter(filter)),
    push: (location: string) => dispatch(push(location)),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Catalog);
