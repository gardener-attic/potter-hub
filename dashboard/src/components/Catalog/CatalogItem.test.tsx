import { shallow } from "enzyme";
import context from "jest-plugin-context";
import { cloneDeep } from "lodash";
import * as React from "react";

import { IChart, IRepo } from "../../shared/types";
import { CardIcon } from "../Card";
import InfoCard from "../InfoCard";
import CatalogItem from "./CatalogItem";

jest.mock("../../img/placeholder.png", () => "img/placeholder.png");

const defaultChart = {
  id: "foo",
  attributes: {
    description: "",
    keywords: [""],
    maintainers: [{ name: "" }],
    sources: [""],
    icon: "icon.png",
    name: "foo",
    repo: {} as IRepo,
  },
  relationships: {
    latestChartVersion: {
      data: {
        version: "4.5.6",
        app_version: "1.0.0",
      },
    },
  },
} as IChart;

it("should render an item", () => {
  const wrapper = shallow(<CatalogItem isFeatured={false} chart={defaultChart} />);
  expect(wrapper).toMatchSnapshot();
});

it("should use the default placeholder for the icon if it doesn't exist", () => {
  const chartWithoutIcon = cloneDeep(defaultChart);
  chartWithoutIcon.attributes.icon = undefined;
  const wrapper = shallow(<CatalogItem isFeatured={false} chart={chartWithoutIcon} />);
  // Importing an image returns "undefined"
  expect(
    wrapper
      .find(InfoCard)
      .shallow()
      .find(CardIcon)
      .prop("src"),
  ).toBe(undefined);
});

it("should place a dash if the version is not avaliable", () => {
  const chartWithoutVersion = cloneDeep(defaultChart);
  chartWithoutVersion.relationships.latestChartVersion.data.app_version = "";
  const wrapper = shallow(<CatalogItem isFeatured={false} chart={chartWithoutVersion} />);
  expect(
    wrapper
      .find(InfoCard)
      .shallow()
      .find(".type-color-light-blue")
      .text(),
  ).toBe("Chart: 4.5.6App: -");
});

it("show the chart description", () => {
  const chartWithDescription = cloneDeep(defaultChart);
  chartWithDescription.attributes.description = "This is a description";
  const wrapper = shallow(<CatalogItem isFeatured={false} chart={chartWithDescription} />);
  expect(
    wrapper
      .find(InfoCard)
      .shallow()
      .find(".ListItem__content__description")
      .text(),
  ).toBe(chartWithDescription.attributes.description);
});

context("when the description is too long", () => {
  it("trims the description", () => {
    const chartWithDescription = cloneDeep(defaultChart);
    chartWithDescription.attributes.description =
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum ultrices velit leo, quis pharetra mi vestibulum quis.";
    const wrapper = shallow(<CatalogItem isFeatured={false} chart={chartWithDescription} />);
    expect(
      wrapper
        .find(InfoCard)
        .shallow()
        .find(".ListItem__content__description")
        .text(),
    ).toMatch(/\.\.\.$/);
  });
});
