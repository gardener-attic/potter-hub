import React from "react"
import { IURLParams } from "shared/types";

interface IProps {
  setURLParams: (params: IURLParams) => void;
}

class URLParamsProvider extends React.Component<IProps> {
  public componentDidMount() {
    const params: IURLParams = {
      targetClusterSecretName: getClusterNameFromURL(),
      targetClusterSecretNamespace: getNamespaceFromURL()
    }
    this.props.setURLParams(params)
  }

  public render() {
    return this.props.children;
  }
}

function getClusterNameFromURL(): string {
  const urlPath = window.location.pathname.split("/");
  const secretName = decodeURIComponent(urlPath[2]);
  return secretName ? secretName : "";
}

function getNamespaceFromURL(): string {
  const urlPath = window.location.pathname.split("/");
  const namespace = decodeURIComponent(urlPath[1]);
  return namespace;
}

export default URLParamsProvider