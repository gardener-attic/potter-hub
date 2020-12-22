import * as React from "react";

import { IClusterBomState } from "../../reducers/clusterbom";
import { IClusterBom } from "../../shared/types";
import {
  getDisplayConditionStatus,
} from "../../shared/utils";
import { CardGrid } from "../Card";
import { MessageAlert } from "../ErrorAlert";
import InfoCard from "../InfoCard";
import LoadingWrapper from "../LoadingWrapper";
import PageHeader from "../PageHeader";

import placeholder from "../../img/bom-placeholder.png";
import "./ClusterBomList.css";

interface IClusterBomListProps {
  clusterBom: IClusterBomState;
  fetchClusterBoms: () => void;
}

class ClusterBomList extends React.Component<IClusterBomListProps> {
  public componentDidMount() {
    this.props.fetchClusterBoms();
  }

  public render() {
    const {
      clusterBom: { isFetching, error: err },
    } = this.props;

    if (err) {
      return <MessageAlert header={err.message} />;
    }

    return (
      <>
        <PageHeader>
          <h1>Cluster BoMs</h1>
        </PageHeader>
        <LoadingWrapper loaded={!isFetching}>{this.renderContent()}</LoadingWrapper>
      </>
    );
  }

  private renderContent = () => {
    const {
      clusterBom: { items: clusterBoms },
    } = this.props;

    if (clusterBoms.length === 0) {
      return (
        <MessageAlert header="No BoMs found" level="warning">
          <div>
            <p className="margin-v-normal">
              Please check that the namespace and the secret name of the target cluster in the URL are correct.
            </p>
          </div>
        </MessageAlert>
      );
    } else {
      return <CardGrid>{clusterBoms.map(cb => this.createItemCard(cb))}</CardGrid>;
    }
  };

  private createItemCard(cb: IClusterBom) {
    const readyCondition = cb.status?.conditions.find(e => {
      return e.type.toLowerCase() === "ready";
    });
    const progress = cb?.status?.overallProgress ? cb?.status?.overallProgress : 0
    const description = <div>
        Apps: {cb.spec?.applicationConfigs.length}<br/>
        Progress: {progress}%
    </div>

    return (
      <InfoCard
        key={cb.metadata.uid}
        link={`/clusterboms/${cb.metadata.name}`}
        title={cb.metadata.name}
        icon={placeholder}
        description={description}
        info={`Last updated: ${readyCondition?.lastUpdateTime}`}
        banner={undefined}
        tag1Content={getDisplayConditionStatus(readyCondition)}
        tag1Class={readyCondition?.status.toLocaleLowerCase()}
      />
    );
  }
}

export default ClusterBomList;
