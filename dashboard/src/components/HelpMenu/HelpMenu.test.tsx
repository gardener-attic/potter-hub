import { mount } from "enzyme";
import * as React from "react";

import HelpMenu from "./HelpMenu";

const defaultProps = {
    appName: "Potter",
    appVersion: "1.0",
    controllerAppVersion: "0.9"
}

it("opens and closes the menu popup when the help button is clicked", () => {
  const wrapper = mount(<HelpMenu {...defaultProps} />);

  expect(wrapper.state("isOpen")).toEqual(false);
  expect(wrapper.find("button.icon-button").length).toBe(1);
  expect(wrapper.find("section.help-menu").length).toBe(0);

  // open the menu
  wrapper.find("button.icon-button").simulate("click");

  expect(wrapper.state("isOpen")).toEqual(true);
  expect(wrapper.find("button.icon-button").length).toBe(1);
  expect(wrapper.find("section.help-menu").length).toBe(1);
  expect(wrapper.find("section.help-menu").is(":focus")).toBe(true);

  // close the menu
  wrapper.find("button.icon-button").simulate("click");

  expect(wrapper.state("isOpen")).toEqual(false);
  expect(wrapper.find("button.icon-button").length).toBe(1);
  expect(wrapper.find("section.help-menu").length).toBe(0);
});

it("closes the menu popup when the menu loses focus", () => {
  jest.useFakeTimers();
  const wrapper = mount(<HelpMenu {...defaultProps} />);

  // open the menu
  wrapper.find("button.icon-button").simulate("click");

  expect(wrapper.state("isOpen")).toEqual(true);
  expect(wrapper.find("button.icon-button").length).toBe(1);
  expect(wrapper.find("section.help-menu").length).toBe(1);
  expect(wrapper.find("section.help-menu").is(":focus")).toBe(true);

  // close the menu
  wrapper.find("section.help-menu").simulate("blur", { relatedTarget: null });

  jest.runAllTimers();
  wrapper.update();

  expect(wrapper.state("isOpen")).toEqual(false);
  expect(wrapper.find("button.icon-button").length).toBe(1);
  expect(wrapper.find("section.help-menu").length).toBe(0);
});

it("renders the help menu", () => {
  const wrapper = mount(<HelpMenu {...defaultProps} />);
  // open the menu
  wrapper.find("button.icon-button").simulate("click");
  expect(wrapper).toMatchSnapshot();
});
