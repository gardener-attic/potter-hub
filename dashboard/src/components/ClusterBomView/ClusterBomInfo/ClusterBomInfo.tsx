import * as React from "react";

import placeholder from "../../../img/bom-placeholder.png";
import { IClusterBom } from "../../../shared/types";
import Card, { CardContent, CardGrid, CardIcon } from "../../Card";

import "./ClusterBomInfo.css";

interface IClusterBomInfoProps {
  clusterBom?: IClusterBom;
}

class ClusterBomInfo extends React.Component<IClusterBomInfoProps> {
  public render() {
    const { clusterBom } = this.props;

    const readyCondition = clusterBom?.status?.conditions.find(e => {
      return e.type.toLowerCase() === "ready";
    });

    let ok = 0;
    let failed = 0;
    let unknown = 0;
    clusterBom?.status?.applicationStates?.forEach(value => {
      switch (value.state.toLowerCase()) {
        case "ok":
          ok++;
          break;
        case "failed":
          failed++;
          break;
        case "pending":
          unknown++;
          break;
        case "unknown":
          unknown++;
          break;
      }
    });
    const total = clusterBom?.status?.applicationStates
      ? clusterBom.status.applicationStates.length
      : 0;
    const progress = clusterBom?.status?.overallProgress ? clusterBom?.status?.overallProgress : 0

    return (
      <CardGrid className="ClusterBomInfo">
        <Card>
          <CardIcon icon={placeholder} />
          <CardContent>
            <div className="ListItem__content">
              <div>
                <h3 className="ListItem__content__title type-big">{clusterBom?.metadata.name}</h3>
              </div>
              <div className="ListItem__content__info">
                <div>
                  <p className="TotalApps margin-b-reset">
                    Apps:
                    <span className="ListItem__content__info_tag type-small type-color-white padding-t-tiny padding-h-normal ListItem__content__info_tag-1">
                      {total}
                    </span>
                  </p>
                  <p className="AppStates margin-b-reset">
                    App States:
                    <div>
                      {ok > 0 && (
                        <span className="ListItem__content__info_tag type-small type-color-white padding-t-tiny padding-h-normal margin-b-normal true">
                          {ok} ok
                        </span>
                      )}
                      {unknown > 0 && (
                        <span className="ListItem__content__info_tag type-small type-color-white padding-t-tiny padding-h-normal unknown">
                          {unknown} unknown
                        </span>
                      )}
                      {failed > 0 && (
                        <span className="ListItem__content__info_tag type-small type-color-white padding-t-tiny padding-h-normal false">
                          {failed} failed
                        </span>
                      )}
                    </div>
                  </p>
                  <p className="TotalApps margin-b-reset">
                    Progress: {progress}%
                  </p>
                </div>
              </div>
              <hr className="separator-small" />
              <p className="margin-reset type-small padding-t-tiny type-color-light-blue">{`Last reconcile: ${readyCondition?.lastUpdateTime}`}</p>
            </div>
          </CardContent>
        </Card>
      </CardGrid>
    );
  }
}

export default ClusterBomInfo;
