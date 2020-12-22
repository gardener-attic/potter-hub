import { mount } from "enzyme";
import * as React from "react";
import CreatableSelect from "react-select/creatable";

import { IAppRepositoryState } from "../../reducers/repos";
import { allRepositories } from "../../shared/Catalog";
import { IAppRepository } from "../../shared/types";

import RepoSelector from "./RepoSelector";

const defaultAppRepoState = {
  isFetching: false,
  repos: [
    {
      metadata: {
        name: "stable",
      },
    } as IAppRepository,
    {
      metadata: {
        name: "incubator",
      },
    } as IAppRepository,
  ] as IAppRepository[],
  errors: {},
} as IAppRepositoryState;
const defaultProps = {
  repo: "",
  appRepoState: defaultAppRepoState,
  push: jest.fn(),
  fetchApprepositories: jest.fn(),
  defaultRepo: "",
};

it("renders successful", () => {
  const wrapper = mount(<RepoSelector {...defaultProps} repo="stable" />);

  const select = wrapper.find(CreatableSelect);
  expect(select).toExist();
  expect(select.props().options).toEqual([
    { label: allRepositories.label, value: allRepositories.value },
    {
      label: "Popular", options: [
        { label: "stable", value: "stable" },
        { label: "incubator", value: "incubator" },]
    }
  ]);
  expect(select.props().value).toEqual([{ label: "stable", value: "stable" }]);
});

it("does not contain _all entry if only 1 apprepo exists", () => {
  const apprepoState = {
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

  const wrapper = mount(<RepoSelector {...defaultProps} appRepoState={apprepoState} />);

  const select = wrapper.find(CreatableSelect);
  expect(select).toExist();
  expect(select.props().options).toEqual([{ label: "stable", value: "stable" }]);
});

it("selects _all entry by default", () => {
  const spy = jest.fn();
  const wrapper = mount(<RepoSelector {...defaultProps} push={spy} />);

  const select = wrapper.find(CreatableSelect);
  expect(select).toExist();
  expect(spy).toHaveBeenCalledTimes(1);
  expect(spy).toHaveBeenCalledWith(`/catalog/${allRepositories.value}`);
});

it("selects defaultRepo", () => {
  const spy = jest.fn();
  const defaultRepo = "stable";
  const wrapper = mount(
    <RepoSelector {...defaultProps} push={spy} defaultRepo={defaultRepo} />,
  );

  const select = wrapper.find(CreatableSelect);
  expect(select).toExist();
  expect(spy).toHaveBeenCalledTimes(1);
  expect(spy).toHaveBeenCalledWith(`/catalog/${defaultRepo}`);
});

it("selects _all if defaultRepo does not exist", () => {
  const spy = jest.fn();
  const wrapper = mount(
    <RepoSelector {...defaultProps} push={spy} defaultRepo="omg omg omg" />,
  );

  const select = wrapper.find(CreatableSelect);
  expect(select).toExist();
  expect(spy).toHaveBeenCalledTimes(1);
  expect(spy).toHaveBeenCalledWith(`/catalog/${allRepositories.value}`);
});

it("renders for empty repo list", () => {
  const apprepoState = {
    isFetching: false,
    repos: [] as IAppRepository[],
    errors: {},
  } as IAppRepositoryState;
  const spy = jest.fn();
  const wrapper = mount(
    <RepoSelector {...defaultProps} appRepoState={apprepoState} push={spy} />,
  );
  const select = wrapper.find(CreatableSelect);
  expect(select.props().options).toEqual([]);
});

it("check for empty repos", () => {
  const wrapper = mount(<RepoSelector {...defaultProps} repo={""} />);
  const select = wrapper.find(CreatableSelect);
  expect(select.props().value).toEqual([]);
});
