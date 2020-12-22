import { RouterAction } from "connected-react-router";
import * as React from "react";
import CreatableSelect from "react-select/creatable";

import { IAppRepositoryState } from "../../reducers/repos";
import { allRepositories } from "../../shared/Catalog";
import { isHiddenAnnotated } from "../../shared/utils";
import "./RepoSelector.css";

interface IRepoSelectorProps {
  repo: string | undefined;
  appRepoState: IAppRepositoryState;
  defaultRepo: string | undefined;
  push: (location: string) => RouterAction;
}

const customSelectStyle = {
  input: (provided: any) => ({
    ...provided,
    minHeight: "1px",
    maxHeight: "36px",
    paddingTop: "0",
    paddingBottom: "0",
    margin: "0",
  })
};

class RepoSelector extends React.Component<IRepoSelectorProps> {
  public componentDidMount() {
    if (!this.props.appRepoState.isFetching && this.props.appRepoState.repos.length > 0) {
      if (!this.props.repo) {
        this.setDefaultRepo()
      }
    }
  }

  public componentDidUpdate() {
    if (!this.props.appRepoState.isFetching && this.props.appRepoState.repos.length > 0) {
      if (!this.props.repo) {
        this.setDefaultRepo()
      }
    }
  }

  public render() {
    const popularOptions =  this.props.appRepoState.repos.filter(r => !isHiddenAnnotated(r))
    .map(r => ({
      value: r.metadata.name,
      label: r.metadata.name,
    }));
    const allRepositoriesOption = { value: allRepositories.value, label: allRepositories.label };

    // combine popularOptions and allRepositoriesOption into a nested array for grouping
    const options: Array<{ label: string; options: Array<{ value: string; label: string; }> } | {value: string, label: string}> = []

    if (popularOptions.length === 1) {
      options.push(popularOptions[0]);
    }
    if (popularOptions.length > 1) {
      options.push(allRepositoriesOption);
      options.push({label: "Popular", options: popularOptions});
    }
    
    return (
      <div className="RepoSelector">
        <label className="RepoSelector__label type-tiny">REPOSITORY</label>
        <CreatableSelect
          className="RepoSelector__select type-small"
          classNamePrefix="RepoSelector"
          value={popularOptions.concat(allRepositoriesOption).filter((r) => r.value === this.props.repo)}
          options={options}
          isValidNewOption={this.isValidNewOption}
          multi={false}
          onChange={this.handleRepoChange}
          clearable={false}
          styles={customSelectStyle}
        />
      </div>
    );
  }

  private isValidNewOption = () => {
    return false;
  };

  private handleRepoChange = (value: any) => {
    this.props.push(`/catalog/${value.value}`);
  };

  private setDefaultRepo = () => {
    let defaultRepo;
    if (this.props.defaultRepo) {
      defaultRepo = this.props.appRepoState.repos.find(repo => {
        return repo.metadata.name === this.props.defaultRepo;
      })?.metadata.name;
    }
    if (!defaultRepo) {
      if (this.props.appRepoState.repos.length > 1) {
        defaultRepo = allRepositories.value;
      } else {
        defaultRepo = this.props.appRepoState.repos[0].metadata.name;
      }
    }
    this.props.push(`/catalog/${defaultRepo}`);
  }

}

export default RepoSelector;
