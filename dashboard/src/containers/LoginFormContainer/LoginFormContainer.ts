import { connect } from "react-redux";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import LoginForm from "../../components/LoginForm";
import { IStoreState } from "../../shared/types";

function mapStateToProps({
  auth: { authenticated, authenticating, authenticationError }, config: { urlParams: { targetClusterSecretName, targetClusterSecretNamespace } }
}: IStoreState) {
  return {
    authenticated,
    authenticating,
    authenticationError,
    targetClusterSecretName,
    targetClusterSecretNamespace
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    tryToAuthenticateWithOIDC: (validateTokenOnTargetCluster: boolean) => dispatch(actions.auth.tryToAuthenticateWithOIDC(validateTokenOnTargetCluster)),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(LoginForm);
