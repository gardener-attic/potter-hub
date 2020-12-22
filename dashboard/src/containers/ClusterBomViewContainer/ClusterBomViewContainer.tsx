import { connect } from "react-redux";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import ClusterBomView from "../../components/ClusterBomView";
import ResourceRef from "../../shared/ResourceRef";
import { IClusterBom, IStoreState } from "../../shared/types";
import { downloadObjectAsYaml } from "../../shared/utils";

interface IRouteProps {
  match: {
    params: {
      clusterBomName: string;
    };
  };
}

function mapStateToProps({
  kube,
  clusterBom,
  config: { urlParams: { targetClusterSecretNamespace } }
}: IStoreState, {
  match: { params }
}: IRouteProps) {
  return {
    kubeState: kube,
    clusterBomName: params.clusterBomName,
    clusterBomNamespace: targetClusterSecretNamespace,
    clusterBomState: clusterBom,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    getAndWatchClusterBom: (ref: ResourceRef) => dispatch(actions.kube.getAndWatchResource(ref)),
    closeWatch: (ref: ResourceRef) => dispatch(actions.kube.closeWatchResource(ref)),
    handleExport: (clusterBom: IClusterBom) =>
      downloadObjectAsYaml(clusterBom, clusterBom.metadata.name),
    handleReconcile: (clusterBom: IClusterBom) =>
      dispatch(actions.clusterBom.updateClusterBom(clusterBom)),
    clearUpdateError: () => dispatch(actions.clusterBom.clearUpdateError()),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ClusterBomView);
