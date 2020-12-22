import * as React from "react";

import "./ClusterName.css";
import "./NamespaceSelector.css";

interface IProps {
  clustername: string
}

class ClusterName extends React.Component<IProps> {
  public render() {
    const {
      clustername
    } = this.props

    if (clustername !== "") {
      return (
        <div className="ClusterNameSelector margin-r-normal">
          <label className="NamespaceSelector__label type-tiny">CLUSTER</label>
          <br />
          <label className="NamespaceSelector__select type-large">{this.prettifyClustername(clustername)}</label>
        </div>
      );
    }

    return "";
  }

  private prettifyClustername(clustername: string) {
    const prettyName = clustername.replace(/.kubeconfig$/, "");
    return clustername ? prettyName : "";
  }

}

export default ClusterName;
