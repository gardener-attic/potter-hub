import { shallow } from "enzyme";
import * as React from "react";

import { UnexpectedErrorAlert } from ".";
import { IRBACRole } from "../../shared/types";
import ErrorAlertHeader from "./ErrorAlertHeader";
import PermissionsErrorAlert from "./PermissionsErrorAlert";
import PermissionsListItem from "./PermissionsListItem";
import { genericMessage } from "./UnexpectedErrorAlert";

it("renders an error message for the action", () => {
  const roles: IRBACRole[] = [];
  const action = "unit-test";
  const wrapper = shallow(<PermissionsErrorAlert roles={roles} action={action} namespace="test" />);
  const header = wrapper
    .find(UnexpectedErrorAlert)
    .shallow()
    .find(ErrorAlertHeader);
  expect(header).toExist();
  expect(header.shallow().text()).toContain(`You don't have sufficient permissions to ${action}`);
  expect(wrapper.html()).toContain("Ask your administrator for the following RBAC roles:");
  expect(wrapper).toMatchSnapshot();
});

it("renders PermissionsListItem for each RBAC role", () => {
  const roles: IRBACRole[] = [
    {
      apiGroup: "test.kubeapps.com",
      resource: "tests",
      verbs: ["get", "create"],
    },
    {
      apiGroup: "apps",
      namespace: "test",
      resource: "deployments",
      verbs: ["list", "watch"],
    },
  ];
  const wrapper = shallow(<PermissionsErrorAlert roles={roles} action="test" namespace="test" />);
  expect(wrapper.find(PermissionsListItem)).toHaveLength(2);
});

it("renders a information about access control", () => {
  const roles: IRBACRole[] = [];
  const wrapper = shallow(<PermissionsErrorAlert roles={roles} action="test" namespace="test" />);
  expect(wrapper.html()).toMatch(
    /Ask you Gardener project manager to add you as a member, so that you have access to the cluster./,
  );
  expect(wrapper.html()).not.toContain(shallow(genericMessage).html());
});
