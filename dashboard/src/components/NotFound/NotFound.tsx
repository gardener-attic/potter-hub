import * as React from "react";
import hubLogo from "../../img/logo.png";

class NotFound extends React.Component {
  public render() {
    return (
      <div className="text-c align-center margin-t-huge">
        <h3>The page you are looking for can't be found.</h3>
        <img src={hubLogo} alt="Logo" title="Logo" />
      </div>
    );
  }
}

export default NotFound;
