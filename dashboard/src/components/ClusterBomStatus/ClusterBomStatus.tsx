import * as React from "react";
import { AlertTriangle } from "react-feather";

import Check from "../../icons/Check";
import Compass from "../../icons/Compass";
import { IClusterBom, IKubeItem } from "../../shared/types";
import { calculateClusterBomStatus } from "../../shared/utils";

import "./ClusterBomStatus.css";

interface IClusterBomStatusProps {
  clusterBom?: IKubeItem<IClusterBom>;
}

class ClusterBomStatus extends React.Component<IClusterBomStatusProps> {
  public render() {
    if (this.props.clusterBom?.isFetching) {
      return <span className="ApplicationStatus">Loading...</span>;
    }

    const state = calculateClusterBomStatus(this.props.clusterBom?.item);
    let button: any;
    switch (state?.toLocaleLowerCase()) {
      case "true":
        button = this.renderSuccessStatus();
        break;
      case "false":
        button = this.renderFailedStatus();
        break;
      case "unknown":
        button = this.renderPendingStatus();
        break;
      default:
        throw Error("bom status could not be determined");
    }

    return button;
  }

  private renderSuccessStatus() {
    return (
      <span className="ClusterBomStatus ClusterBomStatus--ok">
        <Check className="icon padding-t-tiny" /> Ok
      </span>
    );
  }

  private renderPendingStatus() {
    return (
      <span className="ClusterBomStatus ClusterBomStatus--pending">
        <Compass className="icon padding-t-tiny" /> Pending
      </span>
    );
  }

  private renderFailedStatus() {
    return (
      <span className="ClusterBomStatus ClusterBomStatus--failed">
        <AlertTriangle className="icon padding-t-tiny" /> Failed
      </span>
    );
  }
}

export default ClusterBomStatus;
