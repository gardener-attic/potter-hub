import * as React from "react";

import placeholder from "../../img/placeholder.png";
import { IChart } from "../../shared/types";
import InfoCard from "../InfoCard";

import "./CatalogItem.css";

interface ICatalogItemProps {
  chart: IChart;
  isFeatured?: boolean;
}

// 3 lines description max
const MAX_DESC_LENGTH = 90;

function trimDescription(desc: string): string {
  if (desc.length > MAX_DESC_LENGTH) {
    // Trim to the last word under the max length
    return desc.substr(0, desc.lastIndexOf(" ", MAX_DESC_LENGTH)).concat("...");
  }
  return desc;
}

const CatalogItem: React.SFC<ICatalogItemProps> = props => {
  const { chart, isFeatured } = props;
  const { icon, name, repo } = chart.attributes;
  const iconSrc = icon ? `api/chartsvc/${icon}` : placeholder;
  const latestChartVersion = chart.relationships.latestChartVersion.data.version
  const latestAppVersion = chart.relationships.latestChartVersion.data.app_version;
  const repoTag = <span className="ListItem__content__info_tag_link">{repo.name}</span>;
  const description = (
    <div className="ListItem__content__description">
      {trimDescription(chart.attributes.description)}
    </div>
  );
  const info = <span>
      Chart:&nbsp;{latestChartVersion}<br/>
      App:&nbsp;{latestAppVersion || "-"}
  </span>
  return (
    <InfoCard
      key={`${repo}/${name}`}
      title={name}
      link={`/charts/${chart.id}`}
      info={info}
      icon={iconSrc}
      description={description}
      tag1Content={repoTag}
      tag1Class={repo.name}
      highlight={isFeatured}
      tagContainerCSS={{ display: "flex", alignItems: "center" }}
    />
  );
};

export default CatalogItem;
