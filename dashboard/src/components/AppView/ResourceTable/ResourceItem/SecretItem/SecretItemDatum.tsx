import * as React from "react";
import { Copy, Eye, EyeOff } from "react-feather";

interface ISecretItemDatumProps {
  name: string;
  value: string;
}

interface ISecretItemDatumState {
  hidden: boolean;
}

class SecretItemDatum extends React.PureComponent<ISecretItemDatumProps, ISecretItemDatumState> {
  // Secret datum is hidden by default
  public state: ISecretItemDatumState = {
    hidden: true,
  };

  public render() {
    const { name, value } = this.props;
    const { hidden } = this.state;
    const decodedValue = atob(value);
    return (
      <span className="flex">
        <a id="togglesecret" onClick={this.toggleDisplay} title="Toggle secret">
          {hidden ? <Eye /> : <EyeOff />}
        </a>
        <span className="flex margin-l-normal">
          <a id="copysecret" onClick={this.copySecret} title="Copy plain secret to clipboard">
            <Copy />
          </a>
        </span>
        <span className="flex margin-l-normal">
          <span>{name}:</span>
          {hidden ? (
            <span className="margin-l-small">{decodedValue.length} bytes</span>
          ) : (
            <pre className="SecretContainer">
              <code className="SecretContent">{this.shortenSecret()}</code>
            </pre>
          )}
        </span>
      </span>
    );
  }

  private copySecret = () => {
    const value = this.props.value;
    const decodedValue = atob(value);
    navigator.clipboard.writeText(decodedValue);
  };

  private shortenSecret = () => {
    const maxLength = 128;
    const value = this.props.value;
    const decodedValue = atob(value);
    if (decodedValue.length > maxLength) {
      return decodedValue.substring(0, maxLength - 3) + "...";
    }
    return decodedValue;
  };

  private toggleDisplay = () => {
    this.setState({
      hidden: !this.state.hidden,
    });
  };
}

export default SecretItemDatum;
