import { connect } from "react-redux";

import HelpMenu from "../../components/HelpMenu";
import { IStoreState } from "../../shared/types";

function mapStateToProps({ config }: IStoreState) {
  return {
    appVersion: config.appVersion,
    appName: config.appName,
    config: config.header?.helpMenu,
    controllerAppVersion: config.controllerAppVersion
  };
}

export default connect(mapStateToProps)(HelpMenu);
