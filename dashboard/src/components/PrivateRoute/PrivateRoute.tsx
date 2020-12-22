import * as React from "react";
import { Redirect, Route, RouteComponentProps, RouteProps } from "react-router";
import hubLogo from "../../img/logo.png";
import { definedNamespaces } from "../../shared/Namespace";

type IRouteComponentPropsAndRouteProps = RouteProps & RouteComponentProps<any>;

interface IPrivateRouteProps extends IRouteComponentPropsAndRouteProps {
  authenticated: boolean;
  sessionExpired: boolean;
  namespaces: string[];
  fetchTargetClusterSecretErr?: Error;
}

class PrivateRoute extends React.Component<IPrivateRouteProps> {
  public render() {
    const { authenticated, component: Component, ...rest } = this.props;
    return <Route {...rest} render={this.renderRouteConditionally} />;
  }

  public renderRouteConditionally = (props: RouteComponentProps<any>) => {
    const { sessionExpired, authenticated, component: Component, namespaces } = this.props;
    if (authenticated && Component) {
      if (this.props.fetchTargetClusterSecretErr) {
        return this.renderFetchTargetClusterSecretErr()
      }

      const pathname = props.location.pathname;
      // looks for /ns/:namespace in URL
      const matches = pathname.match(/\/ns\/([^/]*)/);

      if (matches && namespaces.length > 0) {
        if (namespaces.includes(matches[1]) || matches[1] === definedNamespaces.all) {
          return <Component {...props} />;
        } else {
          return <Redirect to={{ pathname: "/" }} />;
        }
      }

      return <Component {...props} />;
    }
    if (sessionExpired) {
      window.location.reload();
    }
    return <Redirect to={{ pathname: "/login", state: { from: props.location } }} />;
  };

  private renderFetchTargetClusterSecretErr = () => {
    return <div className="text-c align-center margin-t-huge">
      <h3>Cannot fetch target cluster secret.<br/>Please check that the namespace and the secret name of the target cluster in the URL are correct.</h3>
      <img src={hubLogo} alt="Logo" title="Logo" />
    </div>
  }
}

export default PrivateRoute;
