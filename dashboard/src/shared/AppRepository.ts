import { axiosWithAuth } from "./AxiosInstance";
import { UI_BACKEND_ROOT_URL } from "./Kube";
import { IAppRepository, IAppRepositoryList } from "./types";

export class AppRepository {
  public static async list(namespace: string) {
    const { data } = await axiosWithAuth.get<IAppRepositoryList>(
      AppRepository.getResourceLink(namespace),
    );
    return data;
  }

  public static async get(name: string, namespace: string) {
    const { data } = await axiosWithAuth.get(AppRepository.getSelfLink(name, namespace));
    return data;
  }

  public static async update(name: string, namespace: string, newApp: IAppRepository) {
    const { data } = await axiosWithAuth.put(AppRepository.getSelfLink(name, namespace), newApp);
    return data;
  }

  public static async delete(name: string, namespace: string) {
    const { data } = await axiosWithAuth.delete(AppRepository.getSelfLink(name, namespace));
    return data;
  }

  public static async create(
    name: string,
    namespace: string,
    url: string,
    auth: any,
    syncJobPodTemplate: any,
  ) {
    const { data } = await axiosWithAuth.post<IAppRepository>(
      AppRepository.getResourceLink(namespace),
      {
        apiVersion: "kubeapps.com/v1alpha1",
        kind: "AppRepository",
        metadata: {
          name,
        },
        spec: { auth, type: "helm", url, syncJobPodTemplate },
      },
    );
    return data;
  }

  private static APIBase: string = UI_BACKEND_ROOT_URL;
  private static APIEndpoint: string = `${AppRepository.APIBase}/apprepositories`;
  
  private static getResourceLink(namespace?: string): string {
    return AppRepository.APIEndpoint;
  }

  private static getSelfLink(name: string, namespace: string): string {
    return `${AppRepository.APIEndpoint}/${name}`;
  }
}
