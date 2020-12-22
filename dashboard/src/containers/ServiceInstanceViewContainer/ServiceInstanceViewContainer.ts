import { connect } from "react-redux";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import ServiceInstanceView from "../../components/ServiceInstanceView";
import { IServiceInstance } from "../../shared/ServiceInstance";
import { IStoreState } from "../../shared/types";

interface IRouteProps {
  match: {
    params: {
      brokerName: string;
      instanceName: string;
      namespace: string;
    };
  };
}

function mapStateToProps(
  { serviceCatalog: catalog }: IStoreState,
  { match: { params } }: IRouteProps,
) {
  const { instanceName, namespace } = params;
  const { bindingsWithSecrets, instances, classes, plans } = catalog;
  return {
    bindingsWithSecrets,
    errors: catalog.errors,
    instances,
    name: instanceName,
    namespace,
    classes,
    plans,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    addBinding: (bindingName: string, instanceName: string, namespace: string, parameters: {}) =>
      dispatch(actions.serviceCatalog.addBinding(bindingName, instanceName, namespace, parameters)),
    deprovision: (instance: IServiceInstance) =>
      dispatch(actions.serviceCatalog.deprovision(instance)),
    getPlans: async () => {
      dispatch(actions.serviceCatalog.getPlans());
    },
    getClasses: async () => {
      dispatch(actions.serviceCatalog.getClasses());
    },
    getInstances: async (ns: string) => {
      dispatch(actions.serviceCatalog.getInstances(ns));
    },
    getBindings: async (ns: string) => {
      dispatch(actions.serviceCatalog.getBindings(ns));
    },
    removeBinding: (name: string, ns: string) =>
      dispatch(actions.serviceCatalog.removeBinding(name, ns)),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceInstanceView);
