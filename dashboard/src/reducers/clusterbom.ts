import { getType } from "typesafe-actions";

import actions from "../actions";
import { ClusterBomAction } from "../actions/clusterbom";
import { IClusterBom } from "../shared/types";

export interface IClusterBomState {
  isFetching: boolean;
  items: IClusterBom[];
  error?: Error;
  updateError?: Error;
}

const initialState: IClusterBomState = {
  isFetching: false,
  items: [],
};

const clusterBomReducer = (
  state: IClusterBomState = initialState,
  action: ClusterBomAction,
): IClusterBomState => {
  switch (action.type) {
    case getType(actions.clusterBom.requestClusterBoms):
      return { ...state, isFetching: true, error: undefined };
    case getType(actions.clusterBom.receiveClusterBoms):
      return { ...state, isFetching: false, items: action.payload };
    case getType(actions.clusterBom.errorClusterBom):
      return { ...state, isFetching: false, error: action.payload };
    case getType(actions.clusterBom.errorUpdateClusterBom):
      return { ...state, updateError: action.payload };
    case getType(actions.clusterBom.clearUpdateError):
      return { ...state, updateError: undefined };
    default:
  }
  return state;
};

export default clusterBomReducer;
