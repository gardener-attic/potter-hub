import { shallow } from "enzyme";
import { Location } from "history";
import context from "jest-plugin-context";
import * as React from "react";
import { Redirect } from "react-router";
import itBehavesLike from "../../shared/specs";

import LoginForm from "./LoginForm";

const emptyLocation: Location = {
  hash: "",
  pathname: "",
  search: "",
  state: "",
};

const defaultProps = {
  authenticate: jest.fn(),
  authenticated: false,
  authenticating: false,
  authenticationError: undefined,
  location: emptyLocation,
  tryToAuthenticateWithOIDC: jest.fn(),
};

const authenticationError = "it's a trap";

describe("componentDidMount", () => {
  it("should call tryToAuthenticateWithOIDC", () => {
    const tryToAuthenticateWithOIDC = jest.fn();
    shallow(<LoginForm {...defaultProps} tryToAuthenticateWithOIDC={tryToAuthenticateWithOIDC} />);
    expect(tryToAuthenticateWithOIDC).toHaveBeenCalled();
  });
});

context("while authenticating", () => {
  itBehavesLike("aLoadingComponent", {
    component: LoginForm,
    props: { ...defaultProps, authenticating: true },
  });
});

it("does not render the token login form", () => {
  const wrapper = shallow(<LoginForm {...defaultProps} />);
  expect(wrapper.find("input#token").exists()).toBe(false);
  expect(wrapper.find(Redirect).exists()).toBe(false);
  expect(wrapper).toMatchSnapshot();
});

it("renders a hint for access control", () => {
  const wrapper = shallow(<LoginForm {...defaultProps} />);
  expect(wrapper.html()).toMatch(
    /Authentication failed./,
  );
});

describe("redirect if authenticated", () => {
  it("redirects to / if no current location", () => {
    const wrapper = shallow(<LoginForm {...defaultProps} authenticated={true} />);
    const redirect = wrapper.find(Redirect);
    expect(redirect.props()).toEqual({ to: { pathname: "/" } });
  });

  it("redirects to previous location", () => {
    const location = Object.assign({}, emptyLocation);
    location.state = { from: "/test" };
    const wrapper = shallow(
      <LoginForm {...defaultProps} authenticated={true} location={location} />,
    );
    const redirect = wrapper.find(Redirect);
    expect(redirect.props()).toEqual({ to: "/test" });
  });
});

it("displays an error if the authentication error is passed", () => {
  const wrapper = shallow(
    <LoginForm {...defaultProps} authenticationError={authenticationError} />,
  );

  expect(wrapper.find(".alert-error").exists()).toBe(true);
  expect(wrapper).toMatchSnapshot();
});
