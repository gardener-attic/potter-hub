import * as qs from "qs";
import { connect } from "react-redux";
import { RouteComponentProps } from "react-router";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import ServiceInstanceList from "../../components/ServiceInstanceList";
import { IStoreState } from "../../shared/types";

function mapStateToProps(
  { serviceCatalog: catalog, namespace }: IStoreState,
  { location }: RouteComponentProps<{ brokerName: string }>,
) {
  const { brokers, classes, plans, instances, isServiceCatalogInstalled } = catalog;
  return {
    brokers,
    classes,
    error: catalog.errors.fetch,
    filter: qs.parse(location.search, { ignoreQueryPrefix: true }).q || "",
    instances,
    isServiceCatalogInstalled,
    namespace: namespace.current,
    plans,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    checkCatalogInstalled: async () => {
      dispatch(actions.serviceCatalog.checkCatalogInstalled());
    },
    getBrokers: async () => {
      dispatch(actions.serviceCatalog.getBrokers());
    },
    getClasses: async () => {
      dispatch(actions.serviceCatalog.getClasses());
    },
    getInstances: async (ns: string) => {
      dispatch(actions.serviceCatalog.getInstances(ns));
    },
    pushSearchFilter: (filter: string) => dispatch(actions.shared.pushSearchFilter(filter) as any),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceInstanceList);
