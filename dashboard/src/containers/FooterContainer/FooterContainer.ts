import { connect } from "react-redux";

import Footer from "../../components/Footer";
import { IStoreState } from "../../shared/types";

function mapStateToProps({ config }: IStoreState) {
  return {
    config: config.footer
  };
}

export default connect(mapStateToProps)(Footer);
