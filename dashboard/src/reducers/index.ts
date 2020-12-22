import { combineReducers } from "redux";

import { IStoreState } from "../shared/types";
import appsReducer from "./apps";
import authReducer from "./auth";
import chartsReducer from "./charts";
import clusterBomReducer from "./clusterbom";
import configReducer from "./config";
import kubeReducer from "./kube";
import namespaceReducer from "./namespace";
import reposReducer from "./repos";
import serviceCatalogReducer from "./serviceCatalog";

const rootReducer = combineReducers<IStoreState>({
  apps: appsReducer,
  auth: authReducer,
  serviceCatalog: serviceCatalogReducer,
  charts: chartsReducer,
  clusterBom: clusterBomReducer,
  config: configReducer,
  kube: kubeReducer,
  namespace: namespaceReducer,
  repos: reposReducer,
});

export default rootReducer;
