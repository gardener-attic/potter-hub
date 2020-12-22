import * as React from "react";

import placeholder from "../../img/placeholder.png";
import "./ChartIcon.css";

interface IChartIconProps {
  icon?: string | null;
}

class ChartIcon extends React.Component<IChartIconProps> {
  public render() {
    const { icon } = this.props;
    const iconSrc = icon ? `api/chartsvc/${icon}` : placeholder;

    return (
      <div className="ChartIcon">
        <img className="ChartIcon__img" src={iconSrc} alt=""/>
      </div>
    );
  }
}

export default ChartIcon;
