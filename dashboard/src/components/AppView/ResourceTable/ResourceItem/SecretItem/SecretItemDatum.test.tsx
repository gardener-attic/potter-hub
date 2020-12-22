import { shallow } from "enzyme";
import * as React from "react";
import SecretItemDatum from "./SecretItemDatum";

const testProps = {
  name: "foo",
  value: "YmFy", // foo
};

// First 125 characters
const shortenedSecret =
  "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna ali";
const testLongProps = {
  name: "foo",
  // The original secret is 300 characters long
  value:
    "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNldGV0dXIgc2FkaXBzY2luZyBlbGl0ciwgc2VkIGRpYW0gbm9udW15IGVpcm1vZCB0ZW1wb3IgaW52aWR1bnQgdXQgbGFib3JlIGV0IGRvbG9yZSBtYWduYSBhbGlxdXlhbSBlcmF0LCBzZWQgZGlhbSB2b2x1cHR1YS4gQXQgdmVybyBlb3MgZXQgYWNjdXNhbSBldCBqdXN0byBkdW8gZG9sb3JlcyBldCBlYSByZWJ1bS4gU3RldCBjbGl0YSBrYXNkIGd1YmVyZ3Jlbiwgbm8gc2VhIHRha2ltYXRhIHNhbmN0dXMgZXN0IExvcmVtIGlwc3VtIGRvbG9yIHNpdCBhbWV0LiBMb3Jl",
};

it("renders the secret datum (hidden by default)", () => {
  const wrapper = shallow(<SecretItemDatum {...testProps} />);
  expect(wrapper.state()).toMatchObject({ hidden: true });
  expect(wrapper).toMatchSnapshot();
});

it("displays the secret datum value when clicking on the icon", () => {
  const wrapper = shallow(<SecretItemDatum {...testProps} />);
  expect(wrapper.text()).toContain("foo:3 bytes");
  const icon = wrapper.find("a#togglesecret");
  expect(icon).toExist();
  icon.simulate("click");
  expect(wrapper.state()).toMatchObject({ hidden: false });
  expect(wrapper.text()).toContain("foo:bar");
});

it("displays the secret datum shortened when clicking on the icon", () => {
  const wrapper = shallow(<SecretItemDatum {...testLongProps} />);
  expect(wrapper.text()).toContain("foo:300 bytes");
  const icon = wrapper.find("a#togglesecret");
  expect(icon).toExist();
  icon.simulate("click");
  expect(wrapper.state()).toMatchObject({ hidden: false });
  expect(wrapper.text()).toContain("foo:" + shortenedSecret);
});
