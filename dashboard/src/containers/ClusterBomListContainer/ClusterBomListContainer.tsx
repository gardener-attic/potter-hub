import { connect } from "react-redux";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import ClusterBomList from "../../components/ClusterBomList";
import { IStoreState } from "../../shared/types";

function mapStateToProps({ clusterBom }: IStoreState) {
  return {
    clusterBom,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    fetchClusterBoms: () => dispatch(actions.clusterBom.fetchClusterBoms()),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ClusterBomList);
