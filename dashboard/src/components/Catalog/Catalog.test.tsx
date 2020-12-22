import InfoBox from "components/InfoBox";
import { mount, shallow } from "enzyme";
import context from "jest-plugin-context";
import * as React from "react";

import { IAppRepositoryState } from "../../reducers/repos";
import { IAppRepository, IChart, IChartState } from "../../shared/types";
import { CardGrid } from "../Card";
import { ErrorSelector, MessageAlert } from "../ErrorAlert";
import LoadingWrapper from "../LoadingWrapper";
import PageHeader from "../PageHeader";
import SearchFilter from "../SearchFilter";
import Catalog from "./Catalog";
import CatalogItem from "./CatalogItem";

const defaultChartState = {
  isFetching: false,
  selected: {} as IChartState["selected"],
  deployed: {} as IChartState["deployed"],
  items: [],
  updatesInfo: {},
} as IChartState;
const defaultAppRepoState = {
  isFetching: false,
  repos: [
    {
      metadata: {
        name: "stable",
      },
    } as IAppRepository,
  ] as IAppRepository[],
  errors: {},
} as IAppRepositoryState;
const defaultProps = {
  appRepoState: defaultAppRepoState,
  charts: defaultChartState,
  filter: "",
  fetchCharts: jest.fn(),
  pushSearchFilter: jest.fn(),
  repo: "",
  fetchApprepositories: jest.fn(),
  onToggleFeaturedCharts: jest.fn(),
  showFeaturedCharts: true,
  isFeatured: false,
  push: jest.fn(),
  defaultRepo: undefined,
  targetClusterSecretNamespace: "ns-1",
  targetClusterSecretName: "my-cluster.kubeconfig",
  staticCatalogInfo: "info msg when catalog is opened without target cluster"
};

it("propagates the filter from the props", () => {
  const wrapper = shallow(<Catalog {...defaultProps} filter="foo" />);
  expect(wrapper.state("filter")).toBe("foo");
});

it("reloads charts when the repo changes", () => {
  const spy = jest.fn();
  const wrapper = shallow(<Catalog {...defaultProps} fetchCharts={spy} />);
  wrapper.setProps({ repo: "bitnami" });
  expect(spy).toHaveBeenCalledTimes(1);
  expect(spy).toHaveBeenCalledWith("bitnami");
});

it("updates the filter from props", () => {
  const wrapper = shallow(<Catalog {...defaultProps} />);
  wrapper.setProps({ filter: "foo" });
  expect(wrapper.state("filter")).toBe("foo");
});

it("keeps the filter from the state", () => {
  const wrapper = shallow(<Catalog {...defaultProps} />);
  expect(wrapper.state("filter")).toBe("");
  wrapper.setState({ filter: "foo" });
  expect(wrapper.state("filter")).toBe("foo");
});

it("renders a LoadingWrapper when fetching charts", () => {
  const wrapper = shallow(
    <Catalog {...defaultProps} charts={{ ...defaultChartState, isFetching: true }} />,
  );
  const loadingWrapper = wrapper.find(LoadingWrapper);
  expect(loadingWrapper).toExist();
  expect(loadingWrapper.at(0).props().loaded).toEqual(true);
  expect(loadingWrapper.at(1).props().loaded).toEqual(false);
});

it("renders an error when fetching apprepos failed", () => {
  const errMessage = "test123";
  const wrapper = mount(
    <Catalog
      {...defaultProps}
      appRepoState={{
        ...defaultAppRepoState,
        errors: { fetch: { name: "", message: errMessage } },
      }}
    />,
  );
  expect(wrapper.find(ErrorSelector)).toExist();
  expect(wrapper.find(ErrorSelector).text()).toContain(errMessage);
});

it("renders a LoadingWrapper when fetching apprepos", () => {
  const wrapper = mount(
    <Catalog {...defaultProps} appRepoState={{ ...defaultAppRepoState, isFetching: true }} />,
  );
  expect(wrapper.find(LoadingWrapper)).toExist();
  expect(wrapper.find(LoadingWrapper).props().loaded).toEqual(false);
});

describe("renderization", () => {
  context("when no charts", () => {
    it("should render a distinct 'no charts' message", () => {
      const wrapper = shallow(<Catalog {...defaultProps} />);
      expect(wrapper.find(MessageAlert)).toExist();
      // expect(wrapper.find(".Catalog")).not.toExist();
      expect(wrapper.find(MessageAlert).props().header).toEqual("No charts found.");
      expect(wrapper).toMatchSnapshot();
    });
  });

  context("when charts available", () => {
    const chartState = {
      isFetching: false,
      selected: {} as IChartState["selected"],
      items: [
        { id: "foo", attributes: { description: "" } } as IChart,
        { id: "bar", attributes: { description: "" } } as IChart,
      ],
    } as IChartState;

    it("should render the list of charts", () => {
      const wrapper = shallow(<Catalog {...defaultProps} charts={chartState} />);

      expect(wrapper.find(MessageAlert)).not.toExist();
      expect(wrapper.find(PageHeader)).toExist();
      expect(wrapper.find(SearchFilter)).toExist();

      const cardGrid = wrapper.find(CardGrid);
      expect(cardGrid).toExist();
      expect(cardGrid.children().length).toBe(chartState.items.length);
      expect(
        cardGrid
          .children()
          .at(0)
          .props().chart,
      ).toEqual(chartState.items[0]);
      expect(
        cardGrid
          .children()
          .at(1)
          .props().chart,
      ).toEqual(chartState.items[1]);
      expect(wrapper).toMatchSnapshot();
    });

    it("should filter apps", () => {
      // Filter "foo" app
      const wrapper = shallow(<Catalog {...defaultProps} charts={chartState} filter="foo" />);

      const cardGrid = wrapper.find(CardGrid);
      expect(cardGrid).toExist();
      expect(cardGrid.children().length).toBe(1);
      expect(
        cardGrid
          .children()
          .at(0)
          .props().chart,
      ).toEqual(chartState.items[0]);
    });
  });
});

describe("featured charts", () => {
  const charts: IChartState = {
    isFetching: false,
    selected: {} as IChartState["selected"],
    items: [
      { id: "foo", attributes: { description: "" } } as IChart,
      { id: "bar", attributes: { description: "" } } as IChart,
      { id: "flupp", attributes: { description: "" } } as IChart,
    ],
    deployed: {} as IChartState["deployed"],
  };

  it("shows featured catalogItems", () => {
    const featuredChartIds: string[] = ["foo", "bar"];
    const wrapper = shallow(
      <Catalog {...defaultProps} charts={charts} featuredChartIds={featuredChartIds} />,
    );
    expect(
      wrapper
        .find(CatalogItem)
        .at(0)
        .props().isFeatured,
    ).toBeTruthy();
    expect(
      wrapper
        .find(CatalogItem)
        .at(1)
        .props().isFeatured,
    ).toBeTruthy();
    expect(
      !wrapper
        .find(CatalogItem)
        .at(2)
        .props().isFeatured,
    ).toBeTruthy();
  });

  it("not render featured charts if no data", () => {
    const featuredChartIds: string[] = [];
    const wrapper = shallow(
      <Catalog {...defaultProps} charts={charts} featuredChartIds={featuredChartIds} />,
    );
    expect(
      !wrapper
        .find(CatalogItem)
        .at(0)
        .props().isFeatured,
    ).toBeTruthy();
    expect(
      !wrapper
        .find(CatalogItem)
        .at(1)
        .props().isFeatured,
    ).toBeTruthy();
    expect(
      !wrapper
        .find(CatalogItem)
        .at(2)
        .props().isFeatured,
    ).toBeTruthy();
  });
});

describe("staticCatalogInfo message", () => {
  it("doesn't show info message when ui opened for target cluster", () => {
    const wrapper = shallow(<Catalog {...defaultProps} />);
    expect(wrapper.find(InfoBox).length).toEqual(1);
    expect(wrapper.find(InfoBox).at(0).prop("msg")).toEqual("");
  });

  it("shows info message when ui opened without target cluster", () => {
    const wrapper = shallow(<Catalog {...defaultProps} targetClusterSecretName="" targetClusterSecretNamespace="" />);
    expect(wrapper.find(InfoBox).length).toEqual(2);
    expect(wrapper.find(InfoBox).at(0).prop("msg")).toEqual("info msg when catalog is opened without target cluster");
  });
})

