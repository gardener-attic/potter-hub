import { shallow } from "enzyme";
import { has } from "lodash";
import * as React from "react";
import ReactTooltip from "react-tooltip";

import { IK8sList, IKubeItem, IResource } from "shared/types";
import ApplicationStatus from "./ApplicationStatus";

const defaultProps = {
  watchWorkloads: jest.fn(),
  closeWatches: jest.fn(),
  deployments: [],
  statefulsets: [],
  daemonsets: [],
};

const consoleError = global.console.error;
beforeEach(() => {
  // Mute console.error since we are getting a lot of error for rendering the PieChart component
  // more info here: https://github.com/toomuchdesign/react-minimal-pie-chart/issues/131
  global.console.error = jest.fn();
});
afterEach(() => {
  jest.resetAllMocks();
  global.console.error = consoleError;
});

describe("componentDidMount", () => {
  it("calls watchWorkloads", () => {
    const mock = jest.fn();
    shallow(<ApplicationStatus {...defaultProps} watchWorkloads={mock} />);
    expect(mock).toHaveBeenCalled();
  });
});

describe("componentWillUnmount", () => {
  it("calls closeWatches", () => {
    const mock = jest.fn();
    const wrapper = shallow(<ApplicationStatus {...defaultProps} closeWatches={mock} />);
    wrapper.unmount();
    expect(mock).toHaveBeenCalled();
  });
});

it("renders a loading status", () => {
  const deployments = [
    {
      isFetching: true,
    },
  ];
  const wrapper = shallow(<ApplicationStatus {...defaultProps} deployments={deployments} />);
  expect(wrapper.text()).toContain("Loading");
  expect(wrapper).toMatchSnapshot();
});

it("renders a deleting status", () => {
  const deployments = [
    {
      isFetching: false,
    },
  ];
  const wrapper = shallow(
    <ApplicationStatus {...defaultProps} deployments={deployments} info={{ deleted: {} }} />,
  );
  expect(wrapper.text()).toContain("Deleted");
  expect(wrapper).toMatchSnapshot();
});

it("renders a failed status", () => {
  const deployments = [
    {
      isFetching: false,
    },
  ];
  const wrapper = shallow(
    <ApplicationStatus {...defaultProps} deployments={deployments} info={{ status: "failed" }} />,
  );
  expect(wrapper.text()).toContain("failed");
  expect(wrapper).toMatchSnapshot();
});

describe("isFetching", () => {
  const tests: Array<{
    title: string;
    deployments: Array<IKubeItem<IResource | IK8sList<IResource, {}>>>;
    statefulsets: Array<IKubeItem<IResource | IK8sList<IResource, {}>>>;
    daemonsets: Array<IKubeItem<IResource | IK8sList<IResource, {}>>>;
    deployed: boolean;
    totalPods: number;
    readyPods: number;
  }> = [
    {
      title: "shows a warning if no workloads are present",
      deployments: [],
      statefulsets: [],
      daemonsets: [],
      deployed: false,
      totalPods: 0,
      readyPods: 0,
    },
    {
      title: "shows a deploying status if there is a non deployed deployment",
      deployments: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo" },
            spec: {
              replicas: 1,
            },
            status: {
              availableReplicas: 0,
            },
          } as IResource,
        },
      ],
      statefulsets: [],
      daemonsets: [],
      deployed: false,
      totalPods: 1,
      readyPods: 0,
    },
    {
      title: "shows a deploying status if there is a non deployed statefulset",
      statefulsets: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo" },
            spec: {
              replicas: 1,
            },
            status: {
              readyReplicas: 0,
            },
          } as IResource,
        },
      ],
      deployments: [],
      daemonsets: [],
      deployed: false,
      totalPods: 1,
      readyPods: 0,
    },
    {
      title: "shows a deploying status if there is a non deployed daemonset",
      daemonsets: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo" },
            status: {
              currentNumberScheduled: 1,
              numberReady: 0,
            },
          } as IResource,
        },
      ],
      deployments: [],
      statefulsets: [],
      deployed: false,
      totalPods: 1,
      readyPods: 0,
    },
    {
      title: "shows a deployed status if it has a daemonset, deployment and statefulset deployed",
      daemonsets: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo" },
            status: {
              currentNumberScheduled: 1,
              numberReady: 1,
            },
          } as IResource,
        },
      ],
      deployments: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo" },
            spec: {
              replicas: 1,
            },
            status: {
              availableReplicas: 1,
            },
          } as IResource,
        },
      ],
      statefulsets: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo" },
            spec: {
              replicas: 1,
            },
            status: {
              readyReplicas: 1,
            },
          } as IResource,
        },
      ],
      deployed: true,
      totalPods: 3,
      readyPods: 3,
    },
    {
      title:
        "shows a deploying status if it has a daemonset, deployment (deployed) and statefulset (not deployed)",
      daemonsets: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo-ds" },
            status: {
              currentNumberScheduled: 1,
              numberReady: 1,
            },
          } as IResource,
        },
      ],
      deployments: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo-dp" },
            spec: {
              replicas: 1,
            },
            status: {
              availableReplicas: 1,
            },
          } as IResource,
        },
      ],
      statefulsets: [
        {
          isFetching: false,
          item: {
            metadata: { name: "foo-ss" },
            spec: {
              replicas: 1,
            },
            status: {
              readyReplicas: 0,
            },
          } as IResource,
        },
      ],
      deployed: true,
      totalPods: 3,
      readyPods: 2,
    },
    {
      title:
        "shows a deploying status if it has a daemonset, deployment (deployed) and statefulset (not deployed) with lists",
      daemonsets: [
        {
          isFetching: false,
          item: {
            items: [
              {
                metadata: { name: "foo-ds" },
                status: {
                  currentNumberScheduled: 1,
                  numberReady: 1,
                },
              } as IResource,
            ],
          } as IK8sList<IResource, {}>,
        },
      ],
      deployments: [
        {
          isFetching: false,
          item: {
            items: [
              {
                metadata: { name: "foo-dp" },
                spec: {
                  replicas: 1,
                },
                status: {
                  availableReplicas: 1,
                },
              } as IResource,
            ],
          } as IK8sList<IResource, {}>,
        },
      ],
      statefulsets: [
        {
          isFetching: false,
          item: {
            items: [
              {
                metadata: { name: "foo-ss" },
                spec: {
                  replicas: 1,
                },
                status: {
                  readyReplicas: 0,
                },
              } as IResource,
            ],
          } as IK8sList<IResource, {}>,
        },
      ],
      deployed: true,
      totalPods: 3,
      readyPods: 2,
    },
  ];
  tests.forEach(t => {
    it(t.title, () => {
      const wrapper = shallow(<ApplicationStatus {...defaultProps} />);
      wrapper.setProps({
        deployments: t.deployments,
        statefulsets: t.statefulsets,
        daemonsets: t.daemonsets,
      });
      const getItem = (i?: IResource | IK8sList<IResource, {}>): IResource => {
        return has(i, "items") ? (i as IK8sList<IResource, {}>).items[0] : (i as IResource);
      };
      if (!t.deployments.length && !t.statefulsets.length && !t.daemonsets.length) {
        expect(wrapper.text()).toContain("No workload found");
        return;
      }
      expect(wrapper.text()).toContain(t.deployed ? "Ready" : "Not Ready");
      expect(wrapper.state()).toMatchObject({ totalPods: t.totalPods, readyPods: t.readyPods });
      // Check tooltip text
      const tooltipText = wrapper
        .find(ReactTooltip)
        .dive()
        .text();
      t.deployments.forEach(d => {
        const item = getItem(d.item);
        expect(tooltipText).toContain(
          `${item.status.availableReplicas}/${item.spec.replicas}${item.metadata.name}`,
        );
      });
      t.statefulsets.forEach(d => {
        const item = getItem(d.item);
        expect(tooltipText).toContain(
          `${item.status.readyReplicas}/${item.spec.replicas}${item.metadata.name}`,
        );
      });
      t.daemonsets.forEach(d => {
        const item = getItem(d.item);
        expect(tooltipText).toContain(
          `${item.status.numberReady}/${item.status.currentNumberScheduled}${item.metadata.name}`,
        );
      });
    });
  });
});
