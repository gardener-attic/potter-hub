import { IURLParams } from "shared/types";
import { getType } from "typesafe-actions";

import actions from "../actions";
import { ConfigAction } from "../actions/config";
import { IConfig } from "../shared/Config";

export interface IConfigState extends IConfig {
  loaded: boolean;
  urlParams: IURLParams;
  controllerAppVersion: string;
}

const initialState: IConfigState = {
  loaded: false,
  namespace: "",
  appVersion: "",
  urlParams: {
      targetClusterSecretName: "",
      targetClusterSecretNamespace: ""
  },
  appName: "",
  controllerAppVersion: ""
};

const configReducer = (state: IConfigState = initialState, action: ConfigAction): IConfigState => {
  switch (action.type) {
    case getType(actions.config.requestConfig):
      return initialState;
    case getType(actions.config.receiveConfig):
      return {
        ...state,
        loaded: true,
        ...action.payload,
      };
    case getType(actions.config.errorConfig):
      return {
        ...state,
        error: action.payload,
      };
    case getType(actions.config.setURLParams):
      return {
          ...state,
          urlParams: action.payload,
      };
    case getType(actions.config.receiveControllerAppVersion):
      return {
          ...state,
          controllerAppVersion: action.payload,
      };
    default:
  }
  return state;
};

export default configReducer;
