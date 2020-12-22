import { shallow } from "enzyme";
import * as React from "react";

import ChartIcon from "../ChartIcon";
import ChartDeployButton from "./ChartDeployButton";
import ChartHeader from "./ChartHeader";

const testProps: any = {
  description: "A Test Chart",
  id: "testrepo/test",
  repo: "testrepo",
  showDeployButton: true
};

it("renders a header for the chart", () => {
  const wrapper = shallow(<ChartHeader {...testProps} />);
  expect(wrapper.text()).toContain("testrepo/test");
  expect(wrapper.text()).toContain("A Test Chart");
  expect(wrapper.find(ChartIcon).exists()).toBe(true);
  expect(wrapper.find(ChartDeployButton).exists()).toBe(true)
  expect(wrapper).toMatchSnapshot();
});

it("uses the icon", () => {
  const wrapper = shallow(<ChartHeader {...testProps} icon="test.jpg" />);
  const icon = wrapper.find(ChartIcon);
  expect(icon.exists()).toBe(true);
  expect(icon.props()).toMatchObject({ icon: "test.jpg" });
});

it("hides the deploy button", () => {
    const wrapper = shallow(<ChartHeader {...testProps} showDeployButton={false} />);
    expect(wrapper.find(ChartDeployButton).exists()).toBe(false)
  });
