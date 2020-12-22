import { shallow } from "enzyme";
import * as React from "react";

import { InfoBox } from "./InfoBox";

it("renders the info msg", () => {
  const info = 'testinfo.<br/><a href="www.test123.com">link to doc</a>.';
  const wrapper = shallow(<InfoBox msg={info} />);
  expect(wrapper.contains(info));
});

it("renders null if the info msg is empty", () => {
  const wrapper = shallow(<InfoBox msg={""} />);
  expect(wrapper.type()).toEqual(null);
});
