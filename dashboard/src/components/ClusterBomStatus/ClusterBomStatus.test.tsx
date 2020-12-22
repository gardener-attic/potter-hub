import { shallow } from "enzyme";
import * as React from "react";

import { IClusterBom, IKubeItem } from "../../shared/types";
import ClusterBomStatus from "./ClusterBomStatus";

it("renders bom status in ok state", () => {
  const kubeItem = {
    item: {
      metadata: {
        generation: 1,
      },
      status: {
        observedGeneration: 1,
        conditions: [
          {
            type: "ready",
            status: "True",
          },
        ],
      },
    } as IClusterBom,
  } as IKubeItem<IClusterBom>;

  const wrapper = shallow(<ClusterBomStatus clusterBom={kubeItem} />);
  const status = wrapper.find(".ClusterBomStatus--ok");
  expect(status).toExist();
  expect(status.text()).toContain("Ok");
});

it("renders bom status in pending state", () => {
  const kubeItem = {
    item: {
      metadata: {
        generation: 2,
      },
      status: {
        observedGeneration: 1,
      },
    } as IClusterBom,
  } as IKubeItem<IClusterBom>;

  const wrapper = shallow(<ClusterBomStatus clusterBom={kubeItem} />);
  const status = wrapper.find(".ClusterBomStatus--pending");
  expect(status).toExist();
  expect(status.text()).toContain("Pending");
});

it("renders bom status in failed state", () => {
  const kubeItem = {
    item: {
      metadata: {
        generation: 1,
      },
      status: {
        observedGeneration: 1,
        conditions: [
          {
            type: "ready",
            status: "False",
          },
        ],
      },
    } as IClusterBom,
  } as IKubeItem<IClusterBom>;

  const wrapper = shallow(<ClusterBomStatus clusterBom={kubeItem} />);
  const status = wrapper.find(".ClusterBomStatus--failed");
  expect(status).toExist();
  expect(status.text()).toContain("Failed");
});
