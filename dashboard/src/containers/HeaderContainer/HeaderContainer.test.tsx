import { shallow } from "enzyme";
import * as React from "react";
import { IAuthState } from "reducers/auth";
import { IConfigState } from "reducers/config";
import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";

import { IKubeState } from "../../shared/types";
import Header from "./HeaderContainer";

const mockStore = configureMockStore([thunk]);

const emptyLocation = {
  hash: "",
  pathname: "",
  search: "",
};

const makeStore = (authenticated: boolean, oidcAuthenticated: boolean) => {
  const state: IAuthState = {
    sessionExpired: false,
    authenticated,
    oidcAuthenticated,
    authenticating: false,
    defaultNamespace: "",
  };
  const kubeState: IKubeState = {
      items: {},
      sockets: {},
  }
  const configState = {
      urlParams: {
        targetClusterSecretName: "cluster.kubeconfig",
        targetClusterSecretNamespace: "ns1"
      }
  } as IConfigState
  return mockStore({ auth: state, kube: kubeState, config: configState, router: { location: emptyLocation } });
};

describe("LoginFormContainer props", () => {
  it("maps authentication redux states to props", () => {
    const store = makeStore(true, true);
    const wrapper = shallow(<Header store={store} />);
    const form = wrapper.find("Header");
    expect(form).toHaveProp({
      authenticated: true,
      hideLogoutLink: true,
    });
  });
});
