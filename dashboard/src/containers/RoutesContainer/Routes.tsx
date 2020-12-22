import * as React from "react";
import { Redirect, Route, RouteComponentProps, RouteProps, Switch } from "react-router";

import NotFound from "../../components/NotFound";
import AppListContainer from "../../containers/AppListContainer";
import AppNewContainer from "../../containers/AppNewContainer";
import AppUpgradeContainer from "../../containers/AppUpgradeContainer";
import AppViewContainer from "../../containers/AppViewContainer";
import CatalogContainer from "../../containers/CatalogContainer";
import ChartViewContainer from "../../containers/ChartViewContainer";
import ClusterBomListContainer from "../../containers/ClusterBomListContainer";
import ClusterBomViewContainer from "../../containers/ClusterBomViewContainer";
import LoginFormContainer from "../../containers/LoginFormContainer";
import PrivateRouteContainer from "../../containers/PrivateRouteContainer";

type IRouteComponentPropsAndRouteProps = RouteProps & RouteComponentProps<any>;

const privateRoutes = {
  "/apps/ns/:namespace": AppListContainer,
  "/apps/ns/:namespace/:releaseName": AppViewContainer,
  "/apps/ns/:namespace/new/:repo/:id/versions/:version": AppNewContainer,
  "/apps/ns/:namespace/upgrade/:releaseName": AppUpgradeContainer,
  "/catalog/:repo?": CatalogContainer,
  "/charts/:repo/:id": ChartViewContainer,
  "/charts/:repo/:id/versions/:version": ChartViewContainer,
  "/clusterboms": ClusterBomListContainer,
  "/clusterboms/:clusterBomName": ClusterBomViewContainer,
} as const;

// Public routes that don't require authentication
const routes = {
  "/login": LoginFormContainer,
} as const;

interface IRoutesProps extends IRouteComponentPropsAndRouteProps {
  namespace: string;
  authenticated: boolean;
  namespaces: string[];
  targetClusterSecretName: string;
  targetClusterSecretNamespace: string;
}

class Routes extends React.Component<IRoutesProps> {
  public render() {
    return (
      <Switch>
        <Route exact={true} path="/" render={this.rootNamespacedRedirect} />
        {Object.entries(routes).map(([route, component]) => (
          <Route key={route} exact={true} path={route} component={component} />
        ))}
        {Object.entries(privateRoutes).map(([route, component]) => (
          <PrivateRouteContainer key={route} exact={true} path={route} component={component} />
        ))}
        {/* If the route doesn't match any expected path redirect to a 404 page  */}
        <Route component={NotFound} />
      </Switch>
    );
  }

  private rootNamespacedRedirect = () => {
    if (this.props.authenticated) {
        if (this.uiCalledForTargetCluster()) {
            if (this.props.namespace) {
              if (this.props.namespaces.length > 0 && !this.props.namespaces.includes(this.props.namespace)) {
                const currentURL = window.location.href
                const fragmentIndex = currentURL.indexOf("#")
                const redirectURL = currentURL.substring(0, fragmentIndex)
                window.location.replace(redirectURL);
              } else {
                return <Redirect to={`/apps/ns/${this.props.namespace}`} />;
              }
            }
        } else {
            return <Redirect to="/catalog" />;
        }
    }
    return <Redirect to={"/login"} />;
  };

  private uiCalledForTargetCluster = () => {
    return this.props.targetClusterSecretNamespace !== "" && this.props.targetClusterSecretName !== "";
  }
}

export default Routes;
