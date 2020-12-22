import { push } from "connected-react-router";
import { connect } from "react-redux";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";

import ServiceClassView from "../../components/ServiceClassView";
import { IStoreState } from "../../shared/types";

interface IRouteProps {
  match: {
    params: {
      brokerName: string;
      className: string;
    };
  };
}

function mapStateToProps(
  { serviceCatalog: catalog, namespace }: IStoreState,
  { match: { params } }: IRouteProps,
) {
  return {
    classes: catalog.classes,
    classname: params.className,
    createError: catalog.errors.create,
    error: catalog.errors.fetch,
    namespace: namespace.current,
    plans: catalog.plans,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    getClasses: async () => {
      dispatch(actions.serviceCatalog.getClasses());
    },
    getPlans: async () => {
      dispatch(actions.serviceCatalog.getPlans());
    },
    provision: (
      instanceName: string,
      namespace: string,
      className: string,
      planName: string,
      parameters: {},
    ) => {
      return dispatch(
        actions.serviceCatalog.provision(instanceName, namespace, className, planName, parameters),
      );
    },
    push: (location: string) => dispatch(push(location)),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceClassView);
