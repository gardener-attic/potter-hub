import { axiosWithAuth } from "./AxiosInstance";
import { K8S_REVERSE_PROXY_URL } from "./Kube";

import { IK8sList, IResource } from "./types";

export default class Namespace {
  public static async list() {
    const { data } = await axiosWithAuth.get<IK8sList<IResource, {}>>(`${Namespace.APIEndpoint}`);
    return data;
  }

  public static async create(namespace: string) {
    const body = {
      apiVersion: "v1",
      kind: "Namespace",
      metadata: {
        name: namespace,
      },
    } as IResource;
    const { data } = await axiosWithAuth.post<IResource>(`${Namespace.APIEndpoint}`, body);
    return data;
  }

  private static APIBase: string = K8S_REVERSE_PROXY_URL;
  private static APIEndpoint: string = `${Namespace.APIBase}/api/v1/namespaces`;
}

// Set of namespaces used accross the applications as default and "all ns" placeholders
export const definedNamespaces = {
  default: "default",
  all: "_all",
};
