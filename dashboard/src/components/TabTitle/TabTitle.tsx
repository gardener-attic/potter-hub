import React from "react"

interface IProps {
  targetClusterSecretNamespace: string;
  targetClusterSecretName: string;
  appName: string;
}

class TabTitle extends React.Component<IProps> {
  public componentDidUpdate() {
    if (this.uiCalledForTargetCluster()) {
      document.title = `${this.props.appName} Dashboard`
    } else {
      document.title = `${this.props.appName} Catalog`
    }
  }

  public render() {
    return this.props.children;
  }

  private uiCalledForTargetCluster = () => {
    return this.props.targetClusterSecretNamespace !== "" && this.props.targetClusterSecretName !== "";
  }
}

export default TabTitle