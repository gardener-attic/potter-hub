import { shallow } from "enzyme";
import * as React from "react";

import { IClusterBomState } from "../../reducers/clusterbom";
import { IClusterBom } from "../../shared/types";
import { MessageAlert } from "../ErrorAlert";
import InfoCard from "../InfoCard";
import LoadingWrapper from "../LoadingWrapper";
import ClusterBomList from "./ClusterBomList";

const defaultClusterBomState = {
  isFetching: false,
  items: [
    {
      metadata: {
        name: "ClusterBoM-1",
      },
    } as IClusterBom,
    {
      metadata: {
        name: "ClusterBoM-2",
      },
    } as IClusterBom,
  ] as IClusterBom[],
} as IClusterBomState;

it("renders successfully", () => {
  const wrapper = shallow(
    <ClusterBomList clusterBom={{ ...defaultClusterBomState }} fetchClusterBoms={jest.fn()} />,
  );

  const infoCards = wrapper.find(InfoCard);
  expect(infoCards.length).toEqual(2);
  expect(infoCards.at(0).props().title).toEqual("ClusterBoM-1");
  expect(infoCards.at(1).props().title).toEqual("ClusterBoM-2");
});

it("renders a loading spinner when fetching ClusterBoMs", () => {
  const wrapper = shallow(
    <ClusterBomList
      clusterBom={{ ...defaultClusterBomState, isFetching: true }}
      fetchClusterBoms={jest.fn()}
    />,
  );

  const loadingWrapper = wrapper.find(LoadingWrapper);
  expect(loadingWrapper).toExist();
  expect(loadingWrapper.props().loaded).toEqual(false);
});

it("renders an alert message if ClusterBoMs are empty", () => {
  const wrapper = shallow(
    <ClusterBomList
      clusterBom={{ ...defaultClusterBomState, items: [] }}
      fetchClusterBoms={jest.fn()}
    />,
  );

  const messageAlert = wrapper.find(MessageAlert);
  expect(messageAlert).toExist();
  expect(messageAlert.props().header).toEqual("No BoMs found");
});

it("renders error message if an error occured", () => {
  const err = {
    message: "test error message",
  } as Error;

  const wrapper = shallow(
    <ClusterBomList
      clusterBom={{ ...defaultClusterBomState, error: err }}
      fetchClusterBoms={jest.fn()}
    />,
  );

  const messageAlert = wrapper.find(MessageAlert);
  expect(messageAlert).toExist();
  expect(messageAlert.props().header).toEqual("test error message");
});
