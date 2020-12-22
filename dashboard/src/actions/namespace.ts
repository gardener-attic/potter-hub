import { ThunkAction } from "redux-thunk";

import { ActionType, createAction } from "typesafe-actions";

import Namespace from "../shared/Namespace";
import { IResource, IStoreState } from "../shared/types";


export const setNamespace = createAction("SET_NAMESPACE", resolve => {
  return (namespace: string) => resolve(namespace);
});

export const receiveNamespaces = createAction("RECEIVE_NAMESPACES", resolve => {
  return (namespaces: string[]) => resolve(namespaces);
});

export const errorNamespaces = createAction("ERROR_NAMESPACES", resolve => {
  return (err: Error, op: "list" | "create") => resolve({ err, op });
});

export const clearNamespaces = createAction("CLEAR_NAMESPACES");

export const clearNamespaceError = createAction("CLEAR_NAMESPACE_ERROR")

export const requestNamespaces = createAction("REQUEST_NAMESPACES")

export const createNamespaceCompleted = createAction("CREATE_NAMESPACE_COMPLETED")

const allActions = [setNamespace, receiveNamespaces, errorNamespaces, clearNamespaces, clearNamespaceError, requestNamespaces, createNamespaceCompleted];
export type NamespaceAction = ActionType<typeof allActions[number]>;

export function fetchNamespaces(): ThunkAction<Promise<void>, IStoreState, null, NamespaceAction> {
  return async dispatch => {
    dispatch(requestNamespaces())
    try {
      const namespaces = await Namespace.list();
      const namespaceStrings = namespaces.items.map((n: IResource) => n.metadata.name);
      dispatch(receiveNamespaces(namespaceStrings));
    } catch (e) {
      dispatch(errorNamespaces(e, "list"));

      // Reset namespace if the namespace couldn't be fetched in order to prevent injection attacks 
      dispatch(setNamespace(""))
      return;
    }
  };
}

export function createNamespace(namespace: string): ThunkAction<Promise<void>, IStoreState, null, NamespaceAction> {
  return async dispatch => {
    dispatch(requestNamespaces())
    try {
      await Namespace.create(namespace);
      dispatch(createNamespaceCompleted())
      dispatch(fetchNamespaces())
    } catch (e) {
      dispatch(errorNamespaces(e, "create"));
      return;
    }
  };
}
