import * as React from "react";
import { Lock } from "react-feather";

import { UnexpectedErrorAlert } from ".";
import { IRBACRole } from "../../shared/types";
import { namespaceText } from "./helpers";
import PermissionsListItem from "./PermissionsListItem";

interface IPermissionsErrorPage {
  action: string;
  roles: IRBACRole[];
  namespace: string;
}

class PermissionsErrorPage extends React.Component<IPermissionsErrorPage> {
  public render() {
    const { action, namespace, roles } = this.props;
    return (
      <UnexpectedErrorAlert
        title={
          <span>
            You don't have sufficient permissions to {action} in {namespaceText(namespace)}
          </span>
        }
        icon={Lock}
        showGenericMessage={false}
      >
        <div>
          <p>Ask your administrator for the following RBAC roles:</p>
          <ul className="error__permisions-list">
            {roles.map((r, i) => (
              <PermissionsListItem key={i} namespace={namespace} role={r} />
            ))}
          </ul>
          <p>
            Ask you Gardener project manager to add you as a member, so that you have access to the
            cluster.
          </p>
        </div>
      </UnexpectedErrorAlert>
    );
  }
}

export default PermissionsErrorPage;
