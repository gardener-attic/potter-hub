import { connect } from "react-redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import { ConfigAction } from "../../actions/config";

import URLParamsProvider from "../../components/URLParamsProvider";
import { IStoreState, IURLParams } from "../../shared/types";

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, ConfigAction>) {
  return {
    setURLParams: (params: IURLParams) => dispatch(actions.config.setURLParams(params)),
  };
}

export default connect(undefined, mapDispatchToProps)(URLParamsProvider);
