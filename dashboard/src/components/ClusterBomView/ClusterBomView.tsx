import * as yaml from "js-yaml";
import * as React from "react";
import AceEditor from "react-ace";
import { AlertOctagon } from "react-feather";
import { IClusterBomState } from "reducers/clusterbom";

import ResourceRef from "../../shared/ResourceRef";
import { IClusterBom, IKubeItem, IKubeState } from "../../shared/types";
import { createClusterBomResourceRef } from "../../shared/utils";
import ClusterBomStatus from "../ClusterBomStatus";
import { MessageAlert } from "../ErrorAlert";
import LoadingWrapper from "../LoadingWrapper";
import ClusterBomInfo from "./ClusterBomInfo/ClusterBomInfo";

interface IClusterBomViewProps {
  clusterBomName: string;
  clusterBomNamespace: string;
  kubeState: IKubeState;
  clusterBomState: IClusterBomState;

  clearUpdateError: () => void;
  closeWatch: (ref: ResourceRef) => void;
  getAndWatchClusterBom: (ref: ResourceRef) => void;
  handleReconcile: (clusterBom: IClusterBom) => void;
  handleExport: (clusterBom: IClusterBom) => void;
}

class ClusterBomView extends React.Component<IClusterBomViewProps> {
  public componentWillUnmount() {
    const clusterBomResourceRef = createClusterBomResourceRef(
      this.props.clusterBomName,
      this.props.clusterBomNamespace,
    );
    this.props.closeWatch(clusterBomResourceRef);
  }

  public componentDidMount() {
    const clusterBomResourceRef = createClusterBomResourceRef(
      this.props.clusterBomName,
      this.props.clusterBomNamespace,
    );
    this.props.getAndWatchClusterBom(clusterBomResourceRef);
  }

  public render() {
    const clusterBomResourceRef = createClusterBomResourceRef(
      this.props.clusterBomName,
      this.props.clusterBomNamespace,
    );
    const clusterBom = this.props.kubeState.items[
      clusterBomResourceRef.getResourceURL()
    ] as IKubeItem<IClusterBom>;
    const err = clusterBom?.error;
    const isFetching = clusterBom?.isFetching;
    const item = clusterBom?.item;

    const updateErr = this.props.clusterBomState.updateError;

    if (err) {
      return <MessageAlert header={err.message} />;
    }

    return (
      <>
        <LoadingWrapper loaded={!isFetching}>
          <div className="row collapse-b-tablet">
            <div className="col-3">
              <ClusterBomInfo clusterBom={item} />
            </div>
            <div className="col-9">
              {updateErr && (
                <div className="row padding-t-big">
                  <div className="col-12">
                    <div className="alert alert-error AlertBox" role="alert">
                      <div className="AlertBoxIconContainer">
                        <AlertOctagon />
                      </div>
                      <span>{updateErr.message}</span>
                      <button onClick={this.clearUpdateError} className="alert__close">
                        &times;
                      </button>
                    </div>
                  </div>
                </div>
              )}
              <div className="row padding-t-big">
                <div className="col-4">
                  <ClusterBomStatus clusterBom={clusterBom} />
                </div>
                <div className="col-8 text-r">
                  <span style={{ marginRight: "1em" }}>
                    Last reconcile: {item?.status?.overallTime}
                  </span>
                  <button className="button" onClick={this.handleReconcileClick}>
                    Reconcile
                  </button>
                  <button id="cb-export-btn" className="button" onClick={this.handleExportClick}>
                    Export
                  </button>
                </div>
              </div>
              <div className="padding-b-big">
                <AceEditor
                  mode="yaml"
                  theme="xcode"
                  name="values"
                  width="100%"
                  maxLines={40}
                  setOptions={{ showPrintMargin: false }}
                  editorProps={{ $blockScrolling: Infinity }}
                  value={yaml.dump(item)}
                  readOnly={true}
                />
              </div>
            </div>
          </div>
        </LoadingWrapper>
      </>
    );
  }

  private handleExportClick = () => {
    const clusterBomResourceRef = createClusterBomResourceRef(
      this.props.clusterBomName,
      this.props.clusterBomNamespace,
    );
    const clusterBom = this.props.kubeState.items[
      clusterBomResourceRef.getResourceURL()
    ] as IKubeItem<IClusterBom>;
    if (clusterBom.item) {
        this.props.handleExport(clusterBom.item);
    }
  };

  private handleReconcileClick = () => {
    const clusterBomResourceRef = createClusterBomResourceRef(
      this.props.clusterBomName,
      this.props.clusterBomNamespace,
    );
    const clusterBom = this.props.kubeState.items[
      clusterBomResourceRef.getResourceURL()
    ] as IKubeItem<IClusterBom>;
    if (clusterBom.item) {
      clusterBom.item.metadata.annotations = {
        ...clusterBom.item.metadata.annotations,
        "hub.k8s.sap.com/reconcile": "reconcile",
      };
    }
    if (clusterBom.item) {
        this.props.handleReconcile(clusterBom.item);
    }
  };
  private clearUpdateError = () => {
    this.props.clearUpdateError();
  };
}

export default ClusterBomView;
