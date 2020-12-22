import * as React from "react";

import placeholder from "../../img/placeholder.png";
import { IAppOverview } from "../../shared/types";
import InfoCard from "../InfoCard";
import "./AppListItem.css";

interface IAppListItemProps {
  app: IAppOverview;
}

class AppListItem extends React.Component<IAppListItemProps> {
  public render() {
    const { app } = this.props;
    const icon = app.icon ? app.icon : placeholder;
    const banner =
      app.updateInfo && !app.updateInfo.error && !app.updateInfo.upToDate
        ? "Update available"
        : undefined;
    const info = <span className="AppListItemInfo">
        {app.chart}<br/>
        Chart:&nbsp;{app.chartMetadata.version}<br/>
        App:&nbsp;{app.chartMetadata.appVersion || "-"}
    </span>
    return (
      <InfoCard
        key={app.releaseName}
        link={`/apps/ns/${app.namespace}/${app.releaseName}`}
        title={app.releaseName}
        icon={icon}
        info={info}
        banner={banner}
        tag1Content={app.namespace}
        tag2Content={app.status.toLocaleLowerCase()}
        tag2Class={app.status.toLocaleLowerCase()}
        tag3Content={app.description}
      />
    );
  }
}

export default AppListItem;
