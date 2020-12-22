import { ThunkAction } from "redux-thunk";
import { ActionType, createAction } from "typesafe-actions";

import { Kube } from "../shared/Kube";
import ResourceRef from "../shared/ResourceRef";
import { IResource, IStoreState, NotFoundError } from "../shared/types";

export const requestResource = createAction("REQUEST_RESOURCE", resolve => {
  return (resourceID: string) => resolve(resourceID);
});

export const receiveResource = createAction("RECEIVE_RESOURCE", resolve => {
  return (resource: { key: string; resource: IResource }) => resolve(resource);
});

export const receiveResourceError = createAction("RECEIVE_RESOURCE_ERROR", resolve => {
  return (resource: { key: string; error: Error }) => resolve(resource);
});

// Takes a ResourceRef to open a WebSocket for and a handler to process messages
// from the socket.
export const openWatchResource = createAction("OPEN_WATCH_RESOURCE", resolve => {
  return (
    ref: ResourceRef,
    handler: (e: MessageEvent) => void,
    onError: { closeTimer: () => void; onErrorHandler: () => void },
  ) => resolve({ ref, handler, onError });
});

export const closeWatchResource = createAction("CLOSE_WATCH_RESOURCE", resolve => {
  return (ref: ResourceRef) => resolve(ref);
});

export const errorTargetClusterSecret = createAction("ERROR_TARGET_CLUSTER_SECRET", resolve => {
  return (err: Error) => resolve(err);
});

const allActions = [
  requestResource,
  receiveResource,
  receiveResourceError,
  openWatchResource,
  closeWatchResource,
  errorTargetClusterSecret,
];

export type KubeAction = ActionType<typeof allActions[number]>;

export function getResource(
  ref: ResourceRef,
  polling?: boolean,
): ThunkAction<Promise<void>, IStoreState, null, KubeAction> {
  return async (dispatch, getState) => {
    const key = ref.getResourceURL();

    // Multiple requests for a resource may happen at the same time (e.g. when
    // loading different sections of a view). This guard attempts to prevent
    // multiple requests for a resource that is already being fetched. Since
    // this action is asynchronous, it's possible that multiple requests are let
    // through, but this is not a huge concern.
    const existing = getState().kube.items[key];
    if (existing && existing.isFetching) {
      return;
    }

    // If it's not the first request, we can skip the request REDUX event
    // to avoid the loading animation
    if (!polling) {
      dispatch(requestResource(key));
    }
    try {
      const r = await ref.getResource();
      dispatch(receiveResource({ key, resource: r }));
    } catch (e) {
      dispatch(receiveResourceError({ key, error: e }));
    }
  };
}

export function getAndWatchResource(
  ref: ResourceRef,
): ThunkAction<void, IStoreState, null, KubeAction> {
  return dispatch => {
    dispatch(getResource(ref));
    let timer: NodeJS.Timeout;
    dispatch(
      openWatchResource(
        ref,
        (e: MessageEvent) => {
          const msg = JSON.parse(e.data);
          const resource: IResource = msg.object;
          const key = ref.getResourceURL();
          switch (msg.type) {
            case "ADDED":
              dispatch(receiveResource({ key, resource }));
              break;
            case "MODIFIED":
              dispatch(receiveResource({ key, resource }));
              break;
            case "DELETED":
              const err = new NotFoundError(
                `${Kube.resourcePlural(resource.kind)}.${resource.apiVersion} "${resource.metadata.name}" not found`,
              );
              dispatch(receiveResourceError({ key, error: err }));
              dispatch(closeWatchResource(ref));
              break;
            default:
          }
        },
        {
          onErrorHandler: () => {
            // If the Socket fails, create an interval to re-request the resource
            // every 5 seconds. This interval needs to be closed calling closeTimer
            timer = setInterval(async () => {
              dispatch(getResource(ref, true));
            }, 5000);
          },
          closeTimer: () => {
            clearInterval(timer);
          },
        },
      ),
    );
  };
}

export function fetchTargetClusterSecret(
  ref: ResourceRef,
): ThunkAction<Promise<void>, IStoreState, null, KubeAction> {
  return async (dispatch) => {
    try {
      await ref.getResource();
    } catch (err) {
      dispatch((errorTargetClusterSecret(err)));
    }
  };
}
