import { shallow } from "enzyme";
import { createMemoryHistory } from "history";
import * as React from "react";
import { Redirect, RouteComponentProps } from "react-router";

import PrivateRoute from "./PrivateRoute";

const emptyRouteComponentProps: RouteComponentProps<{}> = {
  history: createMemoryHistory(),
  location: {
    hash: "",
    pathname: "",
    search: "",
    state: "",
  },
  match: {
    isExact: false,
    params: {},
    path: "",
    url: "",
  },
};

class MockComponent extends React.Component { }

it("redirects to the /login route if not authenticated", () => {
  const wrapper = shallow(
    <PrivateRoute
      sessionExpired={false}
      authenticated={false}
      path="/test"
      namespaces={[]}
      component={MockComponent}
      {...emptyRouteComponentProps}
    />,
  );
  const RenderMethod = (wrapper.instance() as PrivateRoute).renderRouteConditionally;
  const wrapper2 = shallow(<RenderMethod {...emptyRouteComponentProps} />);
  expect(wrapper2.find(Redirect).exists()).toBe(true);
  expect(wrapper2.find(Redirect).props()).toMatchObject({
    to: { pathname: "/login" },
  } as any);
});

it("renders the given component when authenticated", () => {
  const wrapper = shallow(
    <PrivateRoute
      sessionExpired={false}
      authenticated={true}
      path="/test"
      namespaces={[]}
      component={MockComponent}
      {...emptyRouteComponentProps}
    />,
  );
  const RenderMethod = (wrapper.instance() as PrivateRoute).renderRouteConditionally;
  const wrapper2 = shallow(<RenderMethod {...emptyRouteComponentProps} />);
  expect(wrapper2.find(MockComponent).exists()).toBe(true);
});

it("renders special error if target cluster secret couldn't be fetched", () => {
  const err = {
    message: "test error thrown"
  } as Error
  const wrapper = shallow(
    <PrivateRoute
      sessionExpired={false}
      authenticated={true}
      path="/test"
      namespaces={[]}
      fetchTargetClusterSecretErr={err}
      component={MockComponent}
      {...emptyRouteComponentProps}
    />,
  );
  const RenderMethod = (wrapper.instance() as PrivateRoute).renderRouteConditionally;
  const wrapper2 = shallow(<RenderMethod {...emptyRouteComponentProps} />);
  expect(wrapper2.text()).toContain("Cannot fetch target cluster secret.")
});

it("renders a redirect to / if the requested namespace is not in the list of target cluster namespaces", () => {
  const location = {
    hash: "",
    pathname: "/ns/unknown-ns",
    search: "",
    state: "",
  };
  const match = {
    isExact: false,
    params: {},
    path: "",
    url: "",
  }

  const wrapper = shallow(
    <PrivateRoute
      sessionExpired={false}
      authenticated={true}
      path="/ns/unknown-ns"
      namespaces={["default"]}
      component={MockComponent}
      location={location}
      match={match}
      history={createMemoryHistory()}
    />,
  );
  const RenderMethod = (wrapper.instance() as PrivateRoute).renderRouteConditionally;
  const wrapper2 = shallow(<RenderMethod match={match} location={location} history={emptyRouteComponentProps.history} />);
  const redirect = wrapper2.find(Redirect)
  expect(redirect.exists()).toBe(true);
  expect(redirect.props().to).toEqual({ pathname: "/" })
});
