import * as React from "react";
import { NavLink } from "react-router-dom";
import logo from "../../img/logo.png";

import HelpMenu from "../../containers/HelpMenuContainer";
import { INamespaceState } from "../../reducers/namespace";
import ResourceRef from "../../shared/ResourceRef";
import { IResource } from "../../shared/types";
import ClusterName from "./ClusterName";
import HeaderLink from "./HeaderLink";
import NamespaceSelector from "./NamespaceSelector";

import "./Header.css";

interface IHeaderProps {
  authenticated: boolean;
  fetchNamespaces: () => void;
  logout: () => void;
  namespace: INamespaceState;
  defaultNamespace: string;
  pathname: string;
  push: (path: string) => void;
  setNamespace: (ns: string) => void;
  hideLogoutLink: boolean;
  createNamespace: (ns: string) => void;
  clearNamespaceError: () => void;
  fetchTargetClusterSecretErr?: Error;
  fetchTargetClusterSecret: (ref: ResourceRef) => void;
  targetClusterSecretNamespace: string;
  targetClusterSecretName: string;
}

interface IHeaderState {
  configOpen: boolean;
  mobileOpen: boolean;
  initiatedTargetClusterSecretCheck: boolean;
}

class Header extends React.Component<IHeaderProps, IHeaderState> {
  constructor(props: any) {
    super(props);

    this.state = {
      configOpen: false,
      mobileOpen: false,
      initiatedTargetClusterSecretCheck: false
    };
  }

  public componentDidUpdate(prevProps: IHeaderProps) {
    if (prevProps.pathname !== this.props.pathname) {
      this.setState({
        configOpen: false,
        mobileOpen: false,
      });
    }

    if (this.uiCalledForTargetCluster()) {
      if (!this.state.initiatedTargetClusterSecretCheck) {
        if (this.props.authenticated) {
          this.props.fetchTargetClusterSecret(this.createTargetClusterSecretRef(
            this.props.targetClusterSecretNamespace,
            this.props.targetClusterSecretName
          ));
          this.setState({
            initiatedTargetClusterSecretCheck: true
          });
        }
      }
    }
  }

  public render() {
    const {
      clearNamespaceError,
      createNamespace,
      fetchNamespaces,
      namespace,
      defaultNamespace,
      authenticated,
      fetchTargetClusterSecretErr,
      targetClusterSecretName
    } = this.props;
    const header = `row header ${this.state.mobileOpen ? "header-open" : ""}`;
    const showNav = authenticated && !fetchTargetClusterSecretErr

    return (
      <section className="gradient-135-brand type-color-reverse type-color-reverse-anchor-reset">
        <div className="container">
          <header className={header}>
            <div className="header__logo">
              <NavLink to="/">
                <img src={logo} alt="Logo" title="Logo" />
              </NavLink>
            </div>
            {showNav && (
              <nav className="header__nav">
                <button
                  className="header__nav__hamburguer"
                  aria-label="Menu"
                  aria-haspopup="true"
                  aria-expanded="false"
                  onClick={this.toggleMobile}
                >
                  <div />
                  <div />
                  <div />
                </button>
                <ul className="header__nav__menu" role="menubar">
                  {this.uiCalledForTargetCluster() &&
                    <>
                      <li key="/apps">
                        <HeaderLink children="Applications" exact={true} namespaced={true} to="/apps" currentNamespace={namespace.current} />
                      </li>
                      <li key="/clusterboms">
                        <HeaderLink children="Cluster BoMs" to="/clusterboms" currentNamespace={namespace.current} />
                      </li>
                    </>
                  }
                  <li key="/catalog">
                    <HeaderLink children="Catalog" to="/catalog" currentNamespace={namespace.current} />
                  </li>
                </ul>
              </nav>
            )}
            <div className="header__nav header__nav-config">
              {showNav && (
                <>
                  {this.uiCalledForTargetCluster() && (
                    <>
                      <ClusterName clustername={targetClusterSecretName} />
                      <NamespaceSelector
                        namespace={namespace}
                        defaultNamespace={defaultNamespace}
                        onChange={this.handleNamespaceChange}
                        fetchNamespaces={fetchNamespaces}
                        createNamespace={createNamespace}
                        clearNamespaceError={clearNamespaceError}
                      />
                    </>
                  )}
                  < HelpMenu />
                </>
              )}
            </div>
          </header>
        </div>
      </section>
    );
  }

  private toggleMobile = () => {
    this.setState({ mobileOpen: !this.state.mobileOpen });
  };

  private handleNamespaceChange = (ns: string) => {
    const { pathname, push, setNamespace } = this.props;
    const to = pathname.replace(/\/ns\/[^/]*/, `/ns/${ns}`);
    if (to === pathname) {
      setNamespace(ns);
    } else {
      push(to);
    }
  };

  private createTargetClusterSecretRef(namespace: string, name: string) {
    const resource = {
      apiVersion: "v1",
      kind: "Secret",
      metadata: {
        name,
        namespace,
        creationTimestamp: "",
        resourceVersion: "",
        uid: "",
        selfLink: "",
      },
    } as IResource;
    const ref = new ResourceRef(resource, undefined, true);

    return ref;
  }

  private uiCalledForTargetCluster = () => {
    return this.props.targetClusterSecretNamespace !== "" && this.props.targetClusterSecretName !== "";
  }

}

export default Header;
