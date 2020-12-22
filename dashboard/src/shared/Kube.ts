import { Auth } from "./Auth";
import { axiosWithAuth } from "./AxiosInstance";
import { ResourceKind, ResourceKindsWithPlurals } from "./ResourceKinds";
import { IResource } from "./types";

export const UI_BACKEND_ROOT_URL = "api/ui-backend";
export const CLUSTERBOM_URL = UI_BACKEND_ROOT_URL + "/v1/clusterboms"
export const HELM_URL = UI_BACKEND_ROOT_URL + "/helm/v1";
export const K8S_REVERSE_PROXY_URL = UI_BACKEND_ROOT_URL + "/reverse-proxy/v1";

export let WebSocketAPIBase: string;
if (window.location.protocol === "https:") {
  WebSocketAPIBase = `wss://${window.location.host}${window.location.pathname}`;
} else {
  WebSocketAPIBase = `ws://${window.location.host}${window.location.pathname}`;
}

// Kube is a lower-level class for interacting with the Kubernetes API. Use
// ResourceRef to interact with a single API resource rather than using Kube
// directly.
export class Kube {
  public static getResourceURL(
    apiVersion: string,
    resource: string,
    isOnResourceCluster: boolean,
    namespace?: string,
    name?: string,
    query?: string,
  ) {
    let url = K8S_REVERSE_PROXY_URL;
    if (isOnResourceCluster) {
      url += "/resourcecluster";
    }

    url += `/${apiVersion === "v1" ? "api/v1" : `apis/${apiVersion}`}`;
    if (namespace) {
      url += `/namespaces/${namespace}`;
    }
    url += `/${resource}`;
    if (name) {
      url += `/${name}`;
    }
    if (query) {
      url += `?${query}`;
    }
    return url;
  }

  public static watchResourceURL(
    apiVersion: string,
    resource: string,
    isOnResourceCluster: boolean,
    namespace?: string,
    name?: string,
    query?: string,
  ) {
    let url = this.getResourceURL(apiVersion, resource, isOnResourceCluster, namespace);
    url = `${WebSocketAPIBase}${url}?watch=true`;
    if (name) {
      url += `&fieldSelector=metadata.name%3D${name}`;
    }
    if (query) {
      url += `&${query}`;
    }
    return url;
  }

  public static async getResource(
    apiVersion: string,
    resource: string,
    isOnResourceCluster: boolean,
    namespace?: string,
    name?: string,
    query?: string,
  ) {
    const { data } = await axiosWithAuth.get<IResource>(
      this.getResourceURL(apiVersion, resource, isOnResourceCluster, namespace, name, query),
    );
    return data;
  }

  // Opens and returns a WebSocket for the requested resource. Note: it is
  // important that this socket be properly closed when no longer needed. The
  // returned WebSocket can be attached to an event listener to read data from
  // the socket.
  public static watchResource(
    apiVersion: string,
    resource: string,
    isOnResourceCluster: boolean,
    namespace?: string,
    name?: string,
    query?: string,
  ) {
    return new WebSocket(
      this.watchResourceURL(apiVersion, resource, isOnResourceCluster, namespace, name, query),
      Auth.wsProtocols(),
    );
  }

  // Gets the plural form of the resource Kind for use in the resource path
  public static resourcePlural(kind: ResourceKind) {
    return ResourceKindsWithPlurals[kind];
  }
}
