import * as yaml from "js-yaml";

import ResourceRef from "./ResourceRef";
import { IClusterBom, ICondition, IK8sObject, IResource } from "./types";

export function escapeRegExp(str: string) {
  return str.replace(/[\-\[\]\/\{\}\(\)\*\+\?\.\\\^\$\|]/g, "\\$&");
}

export function getValueFromEvent(
  e: React.FormEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
) {
  let value: any = e.currentTarget.value;
  switch (e.currentTarget.type) {
    case "checkbox":
      // value is a boolean
      value = value === "true";
      break;
    case "number":
      // value is a number
      value = parseInt(value, 10);
      break;
  }
  return value;
}

export function getInfoAnnotation(obj: IK8sObject<any, any, any>) {
  const infoAnnotation = "hub.k8s.sap.com/info";
  return getAnnotation(obj, infoAnnotation);
}

export function isHiddenAnnotated(obj: IK8sObject<any, any, any>): boolean {
  const hiddenAnnotation = "hub.k8s.sap.com/hidden"
  return getAnnotation(obj, hiddenAnnotation) === "true";
}

function getAnnotation(obj: IK8sObject<any, any, any>, annotationKey: string): string {
  let annotation = "";
  if (
    obj &&
    obj.metadata &&
    obj.metadata.annotations &&
    annotationKey in obj.metadata.annotations
  ) {
    annotation = obj.metadata.annotations[annotationKey];
  }
  return annotation;
}

export function getConditionFromK8sObject(
  clusterBom: IClusterBom | undefined,
  conditionName: string,
): ICondition | undefined {
  const conditions = clusterBom?.status?.conditions;
  if (!conditions) {
    return undefined;
  }

  return conditions.find(e => e.type.toLowerCase() === conditionName);
}

export function getDisplayConditionStatus(condition: ICondition | undefined): string {
  switch (condition?.status.toLowerCase()) {
    case "true":
      return "ok";
    case "false":
      return "failed";
    case "unknown":
      return "pending";
    default:
      return "";
  }
}

export function downloadObjectAsYaml(exportObj: any, exportName: string) {
  const dataStr = "data:text/yaml;charset=utf-8," + encodeURIComponent(yaml.dump(exportObj));
  const downloadAnchorNode = document.createElement("a");
  downloadAnchorNode.setAttribute("href", dataStr);
  downloadAnchorNode.setAttribute("download", exportName + ".yaml");
  document.body.appendChild(downloadAnchorNode); // required for firefox
  downloadAnchorNode.click();
  downloadAnchorNode.remove();
}

export function calculateClusterBomStatus(clusterBom?: IClusterBom): string | undefined {
  const readinessCondition = getConditionFromK8sObject(clusterBom, "ready");
  const generation = clusterBom?.metadata.generation;
  const observedGeneration = clusterBom?.status?.observedGeneration;

  if (generation !== observedGeneration || !clusterBom) {
    return "unknown";
  }
  return readinessCondition?.status;
}

export function createClusterBomResourceRef(name: string, namespace: string): ResourceRef {
  const resource = {
    apiVersion: "hub.k8s.sap.com/v1",
    kind: "ClusterBom",
    type: "crd",
    metadata: {
      name,
      namespace,
      creationTimestamp: "",
      resourceVersion: "",
      uid: "",
      selfLink: "",
    },
  } as IResource;
  const ref = new ResourceRef(resource, undefined, true);

  return ref;
}

export function isAppManagedByBom(appDescription: string | null | undefined) {
  if (!appDescription) {
    return false;
  }

  try {
    const parsedMetadata = JSON.parse(appDescription);
    if (parsedMetadata && parsedMetadata.bomName) {
      return true;
    }
  } catch {
    // do nothing
  }

  if (appDescription.toLocaleLowerCase() === "controllermanaged") {
    return true;
  }

  return false;
}
