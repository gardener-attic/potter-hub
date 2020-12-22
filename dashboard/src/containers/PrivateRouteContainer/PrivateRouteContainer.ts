import { connect } from "react-redux";
import { withRouter } from "react-router";

import PrivateRoute from "../../components/PrivateRoute";
import { IStoreState } from "../../shared/types";

function mapStateToProps({
  auth: { authenticated, oidcAuthenticated, sessionExpired },
  namespace: { namespaces },
  kube: {
      fetchTargetClusterSecretErr
  }
}: IStoreState) {
  return {
    sessionExpired,
    authenticated,
    oidcAuthenticated,
    namespaces,
    fetchTargetClusterSecretErr
  };
}

export default withRouter(connect(mapStateToProps)(PrivateRoute));
