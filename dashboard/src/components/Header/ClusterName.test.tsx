import { shallow } from "enzyme";
import * as React from "react";

import ClusterName from "./ClusterName";

it("renders the cluster name", () => {
  const name = "mycluster"
  const wrapper = shallow(<ClusterName clustername={name} />);
  const label = wrapper.find(".NamespaceSelector__select").first();
  
  const expectedName = "mycluster"
  expect(label.text()).toEqual(expectedName)
});

it("renders the prettified cluster name for names that follow the schema <name>.kubeconfig", () => {
  const name = "mycluster.kubeconfig"
  const wrapper = shallow(<ClusterName clustername={name} />);
  const label = wrapper.find(".NamespaceSelector__select").first();

  const expectedName = "mycluster"
  expect(label.text()).toEqual(expectedName)
});
