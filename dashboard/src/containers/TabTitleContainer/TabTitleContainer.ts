import { connect } from "react-redux";

import TabTitle from "../../components/TabTitle";
import { IStoreState } from "../../shared/types";

function mapStateToProps({
  config: {
    urlParams: {
      targetClusterSecretName,
      targetClusterSecretNamespace
    },
    appName
  }
}: IStoreState) {
  return {
    targetClusterSecretName,
    targetClusterSecretNamespace,
    appName
  };
}

export default connect(mapStateToProps)(TabTitle);
