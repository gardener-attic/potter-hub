import { mount } from "enzyme";
import { createMemoryHistory } from "history";
import * as React from "react";
import { Route, StaticRouter } from "react-router";
import { RouteComponentProps } from "react-router-dom";

import NotFound from "../../components/NotFound";
import Routes from "./Routes";

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

it("invalid path should show a 404 error", () => {
  const wrapper = mount(
    <StaticRouter location="/random" context={{}}>
      <Routes {...emptyRouteComponentProps} namespace={"default"} authenticated={true} namespaces={[]} />
    </StaticRouter>,
  );
  expect(wrapper.find(NotFound)).toExist();
  expect(wrapper.text()).toContain("The page you are looking for can't be found.");
});

it("should render a redirect to the default namespace", () => {
  const wrapper = mount(
    <StaticRouter location="/" context={{}}>
      <Routes {...emptyRouteComponentProps} namespace={"default"} authenticated={true} namespaces={[]}/>
    </StaticRouter>,
  );
  expect(wrapper.find(NotFound)).not.toExist();
  expect(
    wrapper
      .find(Route)
      .props()
      .render().props.to,
  ).toEqual("/apps/ns/default");
});

it("should render a redirect to the login page", () => {
  const wrapper = mount(
    <StaticRouter location="/" context={{}}>
      <Routes {...emptyRouteComponentProps} namespace={""} authenticated={true} namespaces={[]}/>
    </StaticRouter>,
  );
  expect(wrapper.find(NotFound)).not.toExist();
  expect(
    wrapper
      .find(Route)
      .props()
      .render().props.to,
  ).toEqual("/login");
});

it("should render a redirect to the login page (when not authenticated)", () => {
  const wrapper = mount(
    <StaticRouter location="/" context={{}}>
      <Routes {...emptyRouteComponentProps} namespace={"default"} authenticated={false} namespaces={[]}/>
    </StaticRouter>,
  );
  expect(wrapper.find(NotFound)).not.toExist();
  expect(
    wrapper
      .find(Route)
      .props()
      .render().props.to,
  ).toEqual("/login");
});
