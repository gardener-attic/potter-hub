import * as React from "react";
import { HashLink as Link } from "react-router-hash-link";

const LinkRenderer: React.SFC<{}> = (props: any) => {
  if (props.href.startsWith("#")) {
    return <Link to={props.href}>{props.children}</Link>;
  }
  // If it's not a hash link it's an external link since it's rendering
  // the chart README. Because of that, render it as a normal anchor
  return <a href={props.href}>{props.children}</a>;
};

export default LinkRenderer;
