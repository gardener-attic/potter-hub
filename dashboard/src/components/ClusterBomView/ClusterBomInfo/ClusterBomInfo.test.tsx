import { shallow } from "enzyme";
import * as React from "react";

import { IClusterBom } from "../../../shared/types";
import ClusterBomInfo from "./ClusterBomInfo";

const defaultBom = {
  metadata: {
    name: "ClusterBoM-1",
  },
} as IClusterBom;

it("renders successfully", () => {
  const wrapper = shallow(<ClusterBomInfo clusterBom={defaultBom} />);
  expect(wrapper).toExist();
});
