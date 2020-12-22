import { Location } from "history";
import * as React from "react";
import { Lock } from "react-feather";
import { Redirect } from "react-router";

import LoadingWrapper from "../../components/LoadingWrapper";
import "./LoginForm.css";

interface ILoginFormProps {
  authenticated: boolean;
  authenticating: boolean;
  authenticationError: string | undefined;
  tryToAuthenticateWithOIDC: (validateTokenOnTargetCluster: boolean) => void;
  location: Location;
  targetClusterSecretNamespace: string;
  targetClusterSecretName: string;
}

interface ILoginFormState {
  token: string;
}

class LoginForm extends React.Component<ILoginFormProps, ILoginFormState> {
  public state: ILoginFormState = { token: "" };

  public componentDidMount() {
    const validateTokenOnTargetCluster = this.uiCalledForTargetCluster()
    this.props.tryToAuthenticateWithOIDC(validateTokenOnTargetCluster);
  }

  public render() {
    if (this.props.authenticating) {
      return <LoadingWrapper />;
    }
    if (this.props.authenticated) {
      const { from } = this.props.location.state || { from: { pathname: "/" } };
      return <Redirect to={from} />;
    }
    return (
      <section className="LoginForm">
        <div className="LoginForm__container padding-v-bigger bg-skew">
          <div className="container container-tiny">
            {this.props.authenticationError && (
              <div className="alert alert-error margin-c" role="alert">
                An error occured during authentication. Troubleshooting:
                <ul>
                    <li>Check that the namespace and secret name of the target cluster in the URL are correct.</li>
                    <li>Check that the target cluster is reachable.</li>
                    <li>Check that your token is valid on the target cluster.</li>
                </ul>
              </div>
            )}
          </div>
          <div className="bg-skew__pattern bg-skew__pattern-dark type-color-reverse">
            <div className="container">
              <h2>
                <Lock /> Login
              </h2>
              <p>
                Authentication failed.
              </p>
            </div>
          </div>
        </div>
      </section>
    );
  }

  private uiCalledForTargetCluster = () => {
    return this.props.targetClusterSecretNamespace !== "" && this.props.targetClusterSecretName !== "";
  }

}

export default LoginForm;
