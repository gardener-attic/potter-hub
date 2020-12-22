import { ThunkAction } from "redux-thunk";
import { ActionType, createAction } from "typesafe-actions";

import { ClusterBom } from "../shared/ClusterBom";
import { IClusterBom, IStoreState } from "../shared/types";

export const requestClusterBoms = createAction("REQUEST_CLUSTERBOMS");

export const receiveClusterBoms = createAction("RECEIVE_CLUSTERBOMS", resolve => {
  return (clusterBoms: IClusterBom[]) => resolve(clusterBoms);
});

export const errorClusterBom = createAction("ERROR_CLUSTERBOM", resolve => {
  return (err: Error) => resolve(err);
});

export const errorUpdateClusterBom = createAction("ERROR_UPDATE_CLUSTERBOM", resolve => {
  return (err: Error) => resolve(err);
});
export const clearUpdateError = createAction("CLEAR_UPDATE_ERROR", resolve => {
  return () => resolve();
});

const allActions = [
  requestClusterBoms,
  receiveClusterBoms,
  errorClusterBom,
  errorUpdateClusterBom,
  clearUpdateError,
];
export type ClusterBomAction = ActionType<typeof allActions[number]>;

export const updateClusterBom = (
  clusterBom: IClusterBom,
): ThunkAction<Promise<void>, IStoreState, null, ClusterBomAction> => {
  return async dispatch => {
    try {
      await ClusterBom.update(clusterBom);
    } catch (e) {
      dispatch(errorUpdateClusterBom(e));
    }
  };
};

export const fetchClusterBoms = (): ThunkAction<
  Promise<void>,
  IStoreState,
  null,
  ClusterBomAction
> => {
  return async (dispatch, getState) => {
    dispatch(requestClusterBoms());
    try {
      const clusterBoms = await ClusterBom.list();
      dispatch(receiveClusterBoms(clusterBoms.items));
    } catch (e) {
      dispatch(errorClusterBom(e));
    }
  };
};
