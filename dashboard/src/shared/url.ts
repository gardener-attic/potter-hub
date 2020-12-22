import { K8S_REVERSE_PROXY_URL } from "./Kube";
import { definedNamespaces } from "./Namespace";
import { IServiceBroker } from "./ServiceCatalog";
import { IChartVersion } from "./types";

export const app = {
  charts: {
    version: (cv: IChartVersion) =>
      `/charts/${cv.relationships.chart.data.repo.name}/${cv.relationships.chart.data.name}/versions/${cv.attributes.version}`,
  },
};

export const api = {
  apprepostories: {
    base: `${K8S_REVERSE_PROXY_URL}/apis/kubeapps.com/v1alpha1`,
    create: (namespace = definedNamespaces.default) =>
      `${api.apprepostories.base}/namespaces/${namespace}/apprepositories`,
  },

  charts: {
    base: "api/chartsvc/v1",
    get: (id: string) => `${api.charts.base}/charts/${id}`,
    getReadme: (id: string, version: string) =>
      `${api.charts.base}/assets/${id}/versions/${encodeURIComponent(version)}/README.md`,
    getValues: (id: string, version: string) =>
      `${api.charts.base}/assets/${id}/versions/${encodeURIComponent(version)}/values.yaml`,
    getSchema: (id: string, version: string) =>
      `${api.charts.base}/assets/${id}/versions/${encodeURIComponent(version)}/values.schema.json`,
    getVersion: (id: string, version: string) =>
      `${api.charts.base}/charts/${id}/versions/${encodeURIComponent(version)}`,
    list: (repo?: string) => `${api.charts.base}/charts${repo ? `/${repo}` : ""}`,
    listVersions: (id: string) => `${api.charts.get(id)}/versions`,
  },

  serviceinstances: {
    base: `${K8S_REVERSE_PROXY_URL}/apis/servicecatalog.k8s.io/v1beta1`,
    create: (namespace = definedNamespaces.default) =>
      `${api.serviceinstances.base}/namespaces/${namespace}/serviceinstances`,
  },

  clusterservicebrokers: {
    base: `${K8S_REVERSE_PROXY_URL}/apis/servicecatalog.k8s.io/v1beta1`,
    sync: (broker: IServiceBroker) =>
      `${api.clusterservicebrokers.base}/clusterservicebrokers/${broker.metadata.name}`,
  },
};
