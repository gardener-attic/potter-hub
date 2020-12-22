import Axios, { AxiosResponse } from "axios";
import * as jwt from "jsonwebtoken";
import { K8S_REVERSE_PROXY_URL } from "./Kube";

const AuthTokenKey = "kubeapps_auth_token";
const AuthTokenOIDCKey = "kubeapps_auth_token_oidc";

export const DEFAULT_NAMESPACE = "_all";

export class Auth {
  public static getAuthToken() {
    return localStorage.getItem(AuthTokenKey);
  }

  public static setAuthToken(token: string, oidc: boolean) {
    localStorage.setItem(AuthTokenOIDCKey, oidc.toString());
    return localStorage.setItem(AuthTokenKey, token);
  }

  public static unsetAuthToken() {
    return localStorage.removeItem(AuthTokenKey);
  }

  public static usingOIDCToken() {
    return localStorage.getItem(AuthTokenOIDCKey) === "true";
  }

  public static wsProtocols() {
    const token = this.getAuthToken();
    if (!token) {
      return [];
    }
    return ["base64url.bearer.authorization.k8s.io." + token, "binary.k8s.io"];
  }

  public static fetchOptions(): RequestInit {
    const headers = new Headers();
    headers.append("Authorization", `Bearer ${this.getAuthToken()}`);
    return {
      headers,
    };
  }

  // Throws an error if the token is invalid
  public static async validateToken(token: string) {
    try {
      await Axios.get(K8S_REVERSE_PROXY_URL + "/", {
        headers: { Authorization: `Bearer ${token}` },
      });
    } catch (e) {
      const res = e.response as AxiosResponse;
      if (res.status === 401) {
        throw new Error("invalid token");
      } else if (res.status === 403) {
        throw new Error("token not authorised for cluster");
      }
      // A 403 authorization error only occurs if the token resulted in
      // successful authentication. We don't make any assumptions over RBAC
      // for the root "/" nonResourceURL or other required authz permissions
      // until operations on those resources are attempted (though we may
      // want to revisit this in the future).
      // In the hub project we want to reject everyone with a 403 error
      // since he can't do anything in the ui anyway. We only have
      // admin or unauthorised users.
      if (res.status !== 403) {
        throw new Error(`${res.status}: ${res.data}`);
      }
    }
  }

  // fetchOIDCToken does a HEAD request to collect the Bearer token
  // from the authorization header if exists
  public static async fetchOIDCToken(): Promise<string | null> {
    try {
      const { headers } = await Axios.head("");
      if (headers && headers.authorization) {
        const tokenMatch = (headers.authorization as string).match(/Bearer\s(.*)/);
        if (tokenMatch) {
          return tokenMatch[1];
        }
      }
    } catch (e) {
      // Unable to retrieve token
    }
    return null;
  }

  // defaultNamespaceFromToken decodes a jwt token to return the k8s service
  // account namespace.
  // TODO(mnelson): until we call jwt.verify on the token during validateToken above
  // we use a default namespace for both invalid tokens and tokens without the expected
  // key.
  public static defaultNamespaceFromToken(token: string) {
    const payload = jwt.decode(token);
    const namespaceKey = "kubernetes.io/serviceaccount/namespace";
    if (!payload || !payload[namespaceKey]) {
      return DEFAULT_NAMESPACE;
    }
    return payload[namespaceKey];
  }
}
