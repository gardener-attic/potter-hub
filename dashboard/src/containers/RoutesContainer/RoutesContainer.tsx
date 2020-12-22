import { connect } from "react-redux";
import { withRouter } from "react-router";

import { IStoreState } from "../../shared/types";
import Routes from "./Routes";

function mapStateToProps({ auth, namespace, config: { urlParams: { targetClusterSecretName, targetClusterSecretNamespace } } }: IStoreState) {
  return {
    namespace: namespace.current || auth.defaultNamespace,
    authenticated: auth.authenticated,
    namespaces: namespace.namespaces,
    targetClusterSecretName,
    targetClusterSecretNamespace
  };
}

export default withRouter(connect(mapStateToProps)(Routes));
