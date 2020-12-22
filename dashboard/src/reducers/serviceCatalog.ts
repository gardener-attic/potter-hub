import { LOCATION_CHANGE, LocationChangeAction } from "connected-react-router";
import { getType } from "typesafe-actions";

import actions from "../actions";
import { NamespaceAction } from "../actions/namespace";
import { ServiceCatalogAction } from "../actions/serviceCatalog";
import { IClusterServiceClass } from "../shared/ClusterServiceClass";
import { IServiceBindingWithSecret } from "../shared/ServiceBinding";
import { IServiceBroker, IServicePlan } from "../shared/ServiceCatalog";
import { IServiceInstance } from "../shared/ServiceInstance";

export interface IServiceCatalogState {
  bindingsWithSecrets: {
    isFetching: boolean;
    list: IServiceBindingWithSecret[];
  };
  brokers: {
    isFetching: boolean;
    list: IServiceBroker[];
  };
  classes: {
    isFetching: boolean;
    list: IClusterServiceClass[];
  };
  errors: {
    create?: Error;
    fetch?: Error;
    delete?: Error;
    deprovision?: Error;
    update?: Error;
  };
  instances: {
    isFetching: boolean;
    list: IServiceInstance[];
  };
  isChecking: boolean;
  isServiceCatalogInstalled: boolean;
  plans: {
    isFetching: boolean;
    list: IServicePlan[];
  };
}

const initialState: IServiceCatalogState = {
  bindingsWithSecrets: { isFetching: false, list: [] },
  brokers: { isFetching: false, list: [] },
  classes: { isFetching: false, list: [] },
  errors: {},
  instances: { isFetching: false, list: [] },
  isChecking: true,
  isServiceCatalogInstalled: false,
  plans: { isFetching: false, list: [] },
};

const serviceCatalogReducer = (
  state: IServiceCatalogState = initialState,
  action: ServiceCatalogAction | LocationChangeAction | NamespaceAction,
): IServiceCatalogState => {
  const { serviceCatalog } = actions;
  let list: any = [];
  switch (action.type) {
    case getType(serviceCatalog.installed):
      return { ...state, isChecking: false, isServiceCatalogInstalled: true };
    case getType(serviceCatalog.notInstalled):
      return { ...state, isChecking: false, isServiceCatalogInstalled: false };
    case getType(serviceCatalog.checkCatalogInstall):
      return { ...state, isChecking: true };
    case getType(serviceCatalog.requestBrokers):
      list = state.brokers.list;
      return { ...state, brokers: { isFetching: true, list } };
    case getType(serviceCatalog.receiveBrokers):
      return { ...state, brokers: { isFetching: false, list: action.payload } };
    case getType(serviceCatalog.requestBindingsWithSecrets):
      list = state.bindingsWithSecrets.list;
      return { ...state, bindingsWithSecrets: { isFetching: true, list } };
    case getType(serviceCatalog.receiveBindingsWithSecrets):
      return { ...state, bindingsWithSecrets: { isFetching: false, list: action.payload } };
    case getType(serviceCatalog.requestClasses):
      list = state.classes.list;
      return { ...state, classes: { isFetching: true, list } };
    case getType(serviceCatalog.receiveClasses):
      return { ...state, classes: { isFetching: false, list: action.payload } };
    case getType(serviceCatalog.requestInstances):
      list = state.instances.list;
      return { ...state, instances: { isFetching: true, list } };
    case getType(serviceCatalog.receiveInstances):
      return { ...state, instances: { isFetching: false, list: action.payload } };
    case getType(serviceCatalog.requestPlans):
      list = state.plans.list;
      return { ...state, plans: { isFetching: true, list } };
    case getType(serviceCatalog.receivePlans):
      return { ...state, plans: { isFetching: false, list: action.payload } };
    case getType(serviceCatalog.errorCatalog):
      const brokers = { ...state.brokers, isFetching: false };
      const bindingsWithSecrets = { ...state.bindingsWithSecrets, isFetching: false };
      const classes = { ...state.classes, isFetching: false };
      const instances = { ...state.instances, isFetching: false };
      const plans = { ...state.plans, isFetching: false };
      return {
        ...state,
        brokers,
        bindingsWithSecrets,
        classes,
        instances,
        plans,
        errors: { [action.payload.op]: action.payload.err },
      };
    case LOCATION_CHANGE:
      return { ...state, errors: {} };
    case getType(actions.namespace.setNamespace):
      return { ...state, errors: {} };
    default:
      return { ...state };
  }
};

export default serviceCatalogReducer;
