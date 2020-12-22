import { RouterAction } from "connected-react-router";
import * as React from "react";

import { IAppRepositoryState } from "../../reducers/repos";
import { allRepositories } from "../../shared/Catalog";
import { IAppRepository, IChart, IChartState } from "../../shared/types";
import { escapeRegExp, getInfoAnnotation } from "../../shared/utils";
import { CardGrid } from "../Card";
import { ErrorSelector, MessageAlert } from "../ErrorAlert";
import InfoBox from "../InfoBox";
import LoadingWrapper from "../LoadingWrapper";
import PageHeader from "../PageHeader";
import RepoSelector from "../RepoSelector";
import SearchFilter from "../SearchFilter";
import CatalogItem from "./CatalogItem";

interface ICatalogProps {
  appRepoState: IAppRepositoryState;
  charts: IChartState;
  featuredChartIds?: string[];
  filter: string;
  generalRepoInfo?: string;
  repo: string | undefined;
  fetchApprepositories: () => void;
  fetchCharts: (repo: string) => void;
  pushSearchFilter: (filter: string) => RouterAction;
  push: (location: string) => RouterAction;
  defaultRepo: string | undefined;
  targetClusterSecretNamespace: string;
  targetClusterSecretName: string;
  staticCatalogInfo?: string;
}

interface ICatalogState {
  filter: string;
}

class Catalog extends React.Component<ICatalogProps, ICatalogState> {
  public state: ICatalogState = {
    filter: "",
  };

  public componentDidMount() {
    const { filter } = this.props;
    this.setState({
      filter,
    });
    this.props.fetchApprepositories();
    this.fetchCharts();
  }

  public componentDidUpdate(prevProps: ICatalogProps) {
    if (this.props.filter !== prevProps.filter) {
      this.setState({ filter: this.props.filter });
    }

    if (this.props.repo !== prevProps.repo) {
      this.fetchCharts();
    }
  }

  public render() {
    const {
      charts: {
        items: allItems,
        selected: { error: chartError },
      },
      featuredChartIds = [],
      appRepoState: {
        isFetching: areReposFetching,
        errors: { fetch: repoFetchError },
        repos
      },
      charts: { isFetching: areChartsFetching },
      pushSearchFilter,
      repo,
      generalRepoInfo,
    } = this.props;
    if (repoFetchError) {
      return <ErrorSelector error={repoFetchError} action="list" resource="appRepositories" />;
    }

    let selectedApprepo: IAppRepository | undefined
    if (!areReposFetching && repos.length > 0) {
      selectedApprepo = repos.find(r => {
        return r.metadata.name === repo;
      });

      if (repo && !selectedApprepo && repo !== allRepositories.value) {
        return <MessageAlert header={"Cannot find repository. Please check the repository name in the URL."} />;
      }
    }

    let infoMsg = "";
    if (selectedApprepo) {
      infoMsg = getInfoAnnotation(selectedApprepo);
    } else if (this.props.repo === allRepositories.value && generalRepoInfo) {
      infoMsg = generalRepoInfo;
    }

    return (
      <section className="Catalog">
        <LoadingWrapper loaded={!areReposFetching}>
          <PageHeader>{this.renderHeader(pushSearchFilter)}</PageHeader>

          <InfoBox msg={infoMsg} />

          <LoadingWrapper loaded={!areChartsFetching}>
            {this.renderContent(chartError, allItems, featuredChartIds)}
          </LoadingWrapper>
        </LoadingWrapper>
      </section>
    );
  }

  private renderHeader(pushSearchFilter: (filter: string) => RouterAction) {
    const staticCatalogInfo = this.props.staticCatalogInfo ? this.props.staticCatalogInfo : ""
    return (
      <>
        <div style={{ display: "flex", alignItems: "center", width: "100%" }}>
          <h1>Catalog</h1>
          <RepoSelector
            repo={this.props.repo}
            defaultRepo={this.props.defaultRepo}
            push={this.props.push}
            appRepoState={this.props.appRepoState}
          />
          <SearchFilter
            className="margin-l-big"
            placeholder="search charts..."
            onChange={this.handleFilterQueryChange}
            value={this.state.filter}
            onSubmit={pushSearchFilter}
          />
        </div>
        {!this.uiCalledForTargetCluster() && (
          <InfoBox
            cssStyles={{ width: "100%" }}
            className="margin-b-small margin-t-small"
            msg={staticCatalogInfo} />
        )}
      </>
    );
  }

  private renderContent(
    chartError: Error | undefined,
    allItems: IChart[],
    featuredChartIds: string[],
  ) {
    if (chartError) {
      return <MessageAlert header={chartError.message} />;
    }

    const items = this.filteredCharts(allItems, this.state.filter);

    if (items.length === 0) {
      return <MessageAlert header={"No charts found."} />;
    }
    const featuredCharts = [] as React.ReactElement[];
    const unfeaturedCharts = [] as React.ReactElement[];

    items.forEach(item => {
      if (featuredChartIds.includes(item.id)) {
        featuredCharts.push(<CatalogItem isFeatured={true} key={item.id} chart={item} />);
      } else {
        unfeaturedCharts.push(<CatalogItem isFeatured={false} key={item.id} chart={item} />);
      }
    });

    const allCharts = [] as React.ReactElement[];

    allCharts.push(...featuredCharts);
    allCharts.push(...unfeaturedCharts);

    return <CardGrid>{allCharts}</CardGrid>;
  }

  private filteredCharts(charts: IChart[], filter: string) {
    return charts.filter(c => new RegExp(escapeRegExp(filter), "i").test(c.id));
  }

  private handleFilterQueryChange = (filter: string) => {
    this.setState({
      filter,
    });
  };

  private fetchCharts = () => {
    if (this.props.repo) {
      if (this.props.repo === allRepositories.value) {
        this.props.fetchCharts("");
      } else {
        this.props.fetchCharts(this.props.repo);
      }
    }
  }

  private uiCalledForTargetCluster = () => {
    return this.props.targetClusterSecretNamespace !== "" && this.props.targetClusterSecretName !== "";
  }

}

export default Catalog;
