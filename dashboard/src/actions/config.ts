import { ThunkAction } from "redux-thunk";
import { ActionType, createAction } from "typesafe-actions";
import { IStoreState, IURLParams } from "../shared/types";

import Config, { IConfig } from "../shared/Config";

export const requestConfig = createAction("REQUEST_CONFIG");
export const receiveConfig = createAction("RECEIVE_CONFIG", resolve => {
  return (config: IConfig) => resolve(config);
});
export const errorConfig = createAction("ERROR_CONFIG", resolve => {
  return (err: Error) => resolve(err);
});
export const setURLParams = createAction("SET_URL_PARAMS", resolve => {
    return (params: IURLParams) => resolve(params)
});

export const receiveControllerAppVersion = createAction("RECEIVE_CONTROLLER_APP_VERSION", resolve => {
    return (controllerAppVersion: string) => resolve(controllerAppVersion);
});

const allActions = [requestConfig, receiveConfig, errorConfig, setURLParams, receiveControllerAppVersion];
export type ConfigAction = ActionType<typeof allActions[number]>;

export function getConfig(): ThunkAction<Promise<void>, IStoreState, null, ConfigAction> {
  return async dispatch => {
    dispatch(requestConfig());
    try {
      const config = await Config.getConfig();
      const controllerAppVersion = await Config.getControllerAppVersion();
      dispatch(receiveConfig(config));
      dispatch(receiveControllerAppVersion(controllerAppVersion))
    } catch (e) {
      dispatch(errorConfig(e));
    }
  };
}
