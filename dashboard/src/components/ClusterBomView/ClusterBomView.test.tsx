import { shallow } from "enzyme";
import * as yaml from "js-yaml";
import * as React from "react";
import AceEditor from "react-ace";

import { IClusterBomState } from "reducers/clusterbom";
import { IKubeItem, IKubeState, IResource } from "../../shared/types";
import { createClusterBomResourceRef } from "../../shared/utils";
import ClusterBomStatus from "../ClusterBomStatus";
import { MessageAlert } from "../ErrorAlert";
import LoadingWrapper from "../LoadingWrapper";
import ClusterBomInfo from "./ClusterBomInfo/ClusterBomInfo";
import ClusterBomView from "./ClusterBomView";

const defaultProps: any = {
  clusterBomName: "ClusterBoM-1",
  clusterBomNamespace: "namespace-1",
  clusterBomState: {} as IClusterBomState,
  fetchClusterBom: jest.fn(),
  handleReconcile: jest.fn(),
  handleExport: jest.fn(),
  getAndWatchClusterBom: jest.fn(),
  closeWatch: jest.fn(),
};

it("renders successfully", () => {
  const cbKey = createClusterBomResourceRef(
    defaultProps.clusterBomName,
    defaultProps.clusterBomNamespace,
  ).getResourceURL();
  const kubeState = {
    items: {
      [cbKey]: {
        isFetching: false,
        item: {
          metadata: {
            name: "ClusterBoM-1",
          },
        } as IResource,
      } as IKubeItem<IResource>,
    },
    sockets: {},
  } as IKubeState;

  const wrapper = shallow(<ClusterBomView {...defaultProps} kubeState={kubeState} />);

  const cbInfo = wrapper.find(ClusterBomInfo);
  expect(cbInfo).toExist();
  expect(cbInfo.props().clusterBom).toEqual(kubeState.items[cbKey].item);

  const cbStatus = wrapper.find(ClusterBomStatus);
  expect(cbStatus).toExist();
  expect(cbStatus.props().clusterBom).toEqual(kubeState.items[cbKey]);

  const exportButton = wrapper.find("#cb-export-btn");
  expect(exportButton).toExist();
  expect(exportButton.text()).toEqual("Export");

  const editor = wrapper.find(AceEditor);
  expect(editor).toExist();
  expect(editor.props().value).toEqual(yaml.dump(kubeState.items[cbKey].item));
});

it("renders a loading spinner when fetching the ClusterBoM", () => {
  const cbKey = createClusterBomResourceRef(
    defaultProps.clusterBomName,
    defaultProps.clusterBomNamespace,
  ).getResourceURL();
  const kubeState = {
    items: {
      [cbKey]: {
        isFetching: true,
        item: {
          metadata: {
            name: "ClusterBoM-1",
          },
        } as IResource,
      } as IKubeItem<IResource>,
    },
    sockets: {},
  } as IKubeState;

  const wrapper = shallow(<ClusterBomView {...defaultProps} kubeState={kubeState} />);

  const loadingWrapper = wrapper.find(LoadingWrapper);
  expect(loadingWrapper).toExist();
  expect(loadingWrapper.props().loaded).toEqual(false);
});

it("renders error message if an error occured", () => {
  const err = {
    message: "test error message",
  } as Error;

  const cbKey = createClusterBomResourceRef(
    defaultProps.clusterBomName,
    defaultProps.clusterBomNamespace,
  ).getResourceURL();
  const kubeState = {
    items: {
      [cbKey]: {
        error: err,
      } as IKubeItem<IResource>,
    },
    sockets: {},
  } as IKubeState;

  const wrapper = shallow(<ClusterBomView {...defaultProps} kubeState={kubeState} />);

  const messageAlert = wrapper.find(MessageAlert);
  expect(messageAlert).toExist();
  expect(messageAlert.props().header).toEqual("test error message");
});

it("renders error message if an error occured", () => {
  const clusterBomState = {
    updateError: {
      message: "this is an superdupercool error message",
    } as Error,
  } as IClusterBomState;

  const kubeState = {
    items: {},
    sockets: {},
  } as IKubeState;

  const wrapper = shallow(
    <ClusterBomView {...defaultProps} kubeState={kubeState} clusterBomState={clusterBomState} />,
  );

  const errorDiv = wrapper.find(".alert-error");
  expect(errorDiv).toExist();
  expect(errorDiv.text()).toContain("this is an superdupercool error message");
});
