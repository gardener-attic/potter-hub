import { push } from "connected-react-router";
import { connect } from "react-redux";
import { RouteComponentProps } from "react-router";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import Header from "../../components/Header";
import ResourceRef from "../../shared/ResourceRef";
import { IStoreState } from "../../shared/types";

interface IState extends IStoreState {
  router: RouteComponentProps<{}>;
}

function mapStateToProps({
  auth: { authenticated, oidcAuthenticated, defaultNamespace },
  namespace,
  router: {
    location: { pathname },
  },
  kube: {
      fetchTargetClusterSecretErr
  },
  config: {
      urlParams: {
          targetClusterSecretName,
          targetClusterSecretNamespace
      }
  }
}: IState) {
  return {
    authenticated,
    namespace,
    defaultNamespace,
    pathname,
    // If oidcAuthenticated it's not yet supported to logout
    // Some IdP like Keycloak allows to hit an endpoint to logout:
    // https://www.keycloak.org/docs/latest/securing_apps/index.html#logout-endpoint
    hideLogoutLink: oidcAuthenticated,
    fetchTargetClusterSecretErr,
    targetClusterSecretName,
    targetClusterSecretNamespace,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    fetchNamespaces: () => dispatch(actions.namespace.fetchNamespaces()),
    logout: () => dispatch(actions.auth.logout()),
    push: (path: string) => dispatch(push(path)),
    setNamespace: (ns: string) => dispatch(actions.namespace.setNamespace(ns)),
    createNamespace: (ns: string) => dispatch(actions.namespace.createNamespace(ns)),
    clearNamespaceError: () => dispatch(actions.namespace.clearNamespaceError()),
    fetchTargetClusterSecret: (ref: ResourceRef) => dispatch(actions.kube.fetchTargetClusterSecret(ref)),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Header);
