import * as React from "react";
import { Info } from "react-feather";

import "./InfoBox.css";

export interface IInfoBoxProps {
  msg: string;
  className?: string;
  cssStyles?: React.CSSProperties | undefined;
}

export class InfoBox extends React.Component<IInfoBoxProps, {}> {
  public render() {
    const { msg, cssStyles } = this.props;

    if (msg !== "") {
      return (
        <div className={`alert alert-warning MessageContainer ${this.props.className}`} role="alert" style={cssStyles}>
          <span className="error__icon margin-r-small">
            <Info />
          </span>
          <div dangerouslySetInnerHTML={{ __html: msg }} />
        </div>
      );
    } else {
      return null;
    }
  }
}

export default InfoBox;
