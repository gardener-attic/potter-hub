# Potter-Hub

[Potter-Hub](https://github.com/gardener/potter-hub) is a web-based UI for centrally managing applications in Kubernetes clusters. It was specifically (but not exclusively) adapted for the [Gardener](https://gardener.cloud/) environment.

In order to run Potter-Hub properly, you need two Kubernetes clusters. The first cluster runs the actual Potter-Hub installation. The second cluster holds all the access information (kubeconfigs) of the clusters which should be managed via Potter-Hub (resource cluster). This cluster must also be configured with a working OIDC provider. It is the central place for the user management and the store for all critical access information. In a Gardener landscape, the second cluster is equivalent to the "Garden" cluster.

When installing a new helm chart in a target cluster, Potter-Hub will try to retrieve the kubeconfig of the target cluster from the resource cluster using the current user's OIDC token as authentication. When the OIDC token is not valid or the requested kubeconfig does not exist on the resource cluster, the installation can't be completed.

## Prerequisites

- A K8s cluster with administrative access and a working installation of the [Potter-Controller](https://github.com/gardener/potter-controller) (needed for Cluster-BoM support). This includes a running Ingress controller.
- Helm 3 CLI


## Installation Guide

**1. Configure OIDC with Auth0**
> Optional: If you deploy on a garden landscape, this step is not necessary, since the garden cluster acts as the OIDC provider.

> We use Auth0 as an example, other OIDC provider may also work. The configuration steps should also be applicable on those providers.

1.1. Create an account at Auth0.

1.2. Visit the Auth0 dashboard and create a new "Regular Web Application". Please ensure that the grant type "Authorization Code" is set (in Settings > Advanced Settings (bottom of page) > Grant Types). After configuring the ingress in a later step, you have to set the correct "Allowed Callback URL" on the setting page.

1.3. All target clusters require a configuration with the OIDC provider. We will do it in a later step when creating secrets for those clusters.

1.4. Create all required users. The email adress has to be set as verified.

1.5. This user requires the ClusterRole cluster-admin on the **resource cluster and target cluster**. Modify and apply the following snippet on the resource and target cluster: 
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cluster-admin-<username>
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: <user@user.user>
```
1.6. Additionally, the target kubernetes cluster has to be configured to accept OIDC tokens. For a gardener-managed kubernetes cluster, see the corresponding [Gardener documentation](https://gardener.cloud/documentation/tutorials/oidc-login/#configure-the-shoot-cluster). The issueUrl can be found in the Auth0 settings page. Do not forget to add the trailing "/". 

**2. Pull the Helm Chart**
```
helm repo add <open source helm chart repo> potter
helm pull potter/potter-hub
```

**3.1 Configure the Ingress**

To enable ingress integration, set `ingress.enabled` to `true`.

> **Note**: if you want to expose the UI, you need to set a valid entry in `ingress.hosts`. A host entry follows the following format:
```
  hosts:
  - names:
    - <arbitrary>.<ingress domain of hub cluster>
    path: /
```
You can get the cluster domain from the kubeconfig by removing the api subdomain from the API server URL.

Most likely, you will only want to have one hostname that maps to this potter-hub installation, however, it is possible to have more than one host. To facilitate this, the `ingress.hosts` object is an array.

To set ingress-specific annotations, use the `ingress.annotations` field.

**3.2 Configure TLS for ingress**

TLS can be configured by setting the `ingress.hosts[].tls` of the corresponding hostname to `true`, then you can choose the TLS secret name setting `ingress.hosts[].tlsSecret`. Please see [this example](https://github.com/kubernetes/contrib/tree/master/ingress/controllers/nginx/examples/tls) for more information.

You can provide your own certificates using the `ingress.secrets` object. If your cluster has a [cert-manager](https://github.com/jetstack/cert-manager) add-on to automate the management and issuance of TLS certificates, set `ingress.hosts[].certManager` boolean to true to enable the corresponding annotations for cert-manager. For a full list of configuration parameters related to configuring TLS can see the [values.yaml](values.yaml) file.

If you are using a [Gardener](https://gardener.cloud/) provided cluster you can enable the `ingress.gardenerCertManager`. This will set the necessary annotation at the ingress and gardener will automatically provide the necessary certificates.

**4. Configure Auth Proxy with OIDC values**

The `authProxy` values in the values.yaml have to be replaced according to the OIDC provider that was configured (either the gardener OIDC provider or Auth0 from step 1)

The following values have to be set for a working configuration, for more details, please refer to the values.yaml:
- `authProxy.oidcClusterURL`: URL of the cluster containing the target cluster kubeconfig secrets. Format: `https://api.<domain>`. This url can be found in the kubeconfig for the cluster. 
- `authProxy.oidcClusterCA`: The CA data from the kubeconfig for the cluster containing the target cluster secrets.
- `authProxy.discoveryURL`
- `authProxy.clientID`
- `authProxy.clientSecret`

With the `authProxy.additionalFlags` field, it is possible to disable secure-cookies when the website should be served over http. Therefore, adding `--secure-cookie=false` will allow the browser to store cookies for non-encrypted domains.

> **Note**: Here are some additional flags which might be useful:
>
> - --scopes=<your-scope>
> - --oauth-uri=<your-uri>

In this step, you must also set the allowed callback url in your OIDC provider configuration from step 1.2 to the ingress url. Use http or https as protocol depending on your configuration. Append `/oauth/callback` to the ingress url to provide the correct callback url.

**5. Configure initial Apprepos**

The Helm Chart repositories that should be connected to a Potter-Hub installation can be configured via the parameter `apprepository.initialRepos`. By default, only the repository [Service Catalog](https://svc-catalog-charts.storage.googleapis.com) will be connected.

For a later addition or modification, simply add a new apprepo CR to the namespace of the potter-hub installation (default 'hub').

**6. Installing the Chart**

Install the chart by executing the following shell commands:

```console
helm install potter-hub .
```

> **Caveat**: Only one potter-hub installation is supported per namespace

For a full list of configuration parameters of this chart, see the values.yaml file.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install` or use a second `values.yaml` and provide it with the command `helm install -f </path/to/your/values.yaml>`.,

```console
$ helm install --name potter-hub \
  --set frontend.replicas=1
```

The above command sets the replica count for the frontend deployment to one (default is two).

## Optional, but recommended to configure

Parameter | Description | Default | Type
--- | --- | --- | ---
`ingress.gardenerCertManager` | enable Gardener cert managing  | `false` | bool
`ingress.gardenerDNS` | enables Gardener DNS management | `false` | bool
`uiBackend.hubsec` | configure auto-mounted imagepullsecret in target cluster when added in installed helm chart values | `nil` | map
`apprepository.initialRepos` | set of repos which will be crawled  | `stable, incubator, svc-cat` | map
`ingress.health.user` | username for health check endpoints | `admin` | string
`ingress.health.password` | password for health check endpoints | `admin` | string
