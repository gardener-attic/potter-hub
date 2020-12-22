import { shallow } from "enzyme";
import * as React from "react";

import { NavLink } from "react-router-dom";
import { INamespaceState } from "../../reducers/namespace";
import Header from "./Header";

const defaultProps = {
  authenticated: true,
  fetchNamespaces: jest.fn(),
  logout: jest.fn(),
  namespace: {
    current: "",
    namespaces: [],
    isFetching: false
  } as INamespaceState,
  defaultNamespace: "kubeapps-user",
  pathname: "",
  push: jest.fn(),
  setNamespace: jest.fn(),
  hideLogoutLink: false,
  createNamespace: jest.fn(),
  clearNamespaceError: jest.fn(),
  fetchTargetClusterSecret: jest.fn(),
  targetClusterSecretName: "test-cluster",
  targetClusterSecretNamespace: "test-ns"
};

it("renders the header links and titles", () => {
  const wrapper = shallow(<Header {...defaultProps} />);
  const menubar = wrapper.find(".header__nav__menu").first();
  const items = menubar.children().map(p => p.props().children.props);
  const expectedItems = [
    { children: "Applications", to: "/apps" },
    { children: "Cluster BoMs", to: "/clusterboms" },
    { children: "Catalog", to: "/catalog" },
  ];
  items.forEach((item, index) => {
    expect(item.children).toBe(expectedItems[index].children);
    expect(item.to).toBe(expectedItems[index].to);
  });
});

it("updates state when the path changes", () => {
  const wrapper = shallow(<Header {...defaultProps} />);
  wrapper.setState({ configOpen: true, mobileOpne: true });
  wrapper.setProps({ pathname: "foo" });
  expect(wrapper.state()).toMatchObject({ configOpen: false, mobileOpen: false });
});

it("renders the namespace switcher", () => {
  const wrapper = shallow(<Header {...defaultProps} />);

  const namespaceSelector = wrapper.find("NamespaceSelector");

  expect(namespaceSelector).toExist();
  expect(namespaceSelector.props()).toEqual(
    expect.objectContaining({
      defaultNamespace: defaultProps.defaultNamespace,
      namespace: defaultProps.namespace,
    }),
  );
});

it("disables the logout link when hideLogoutLink is set", () => {
  const wrapper = shallow(<Header {...defaultProps} hideLogoutLink={true} />);
  const links = wrapper.find(NavLink);
  expect(links.length).toEqual(1);
  links.children().forEach(link => {
    expect(link.text).not.toContain("Logout");
  });
});

it("doesn't render the header menu items when target cluster couldn't be fetched", () => {
  const err = {
    message: "test error thrown"
  } as Error
  const wrapper = shallow(<Header {...defaultProps} fetchTargetClusterSecretErr={err} />);
  const leftSideMenubar = wrapper.find(".header__nav__menu").first();
  const leftSideItems = leftSideMenubar.children()
  expect(leftSideItems.length).toEqual(0)

  const rightSideMenubar = wrapper.find(".header__nav-config").first();
  const rightSideItems = rightSideMenubar.children()
  expect(rightSideItems.length).toEqual(0)
});

it("doesn't render cluster specific header items if no target cluster is set", () => {
  const wrapper = shallow(<Header {...defaultProps} targetClusterSecretName="" targetClusterSecretNamespace="" />);
  expect(wrapper.find("ClusterName").exists()).toBe(false);
  expect(wrapper.find("NamespaceSelector").exists()).toBe(false);
  const menubar = wrapper.find(".header__nav__menu").first();
  const items = menubar.children().map(p => p.props().children.props);
  const expectedItems = [
    { children: "Catalog", to: "/catalog" },
  ];
  items.forEach((item, index) => {
    expect(item.children).toBe(expectedItems[index].children);
    expect(item.to).toBe(expectedItems[index].to);
  });
});

