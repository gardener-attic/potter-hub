---
title: User Interface
type: docs
---

The AppHub UI can be opened in two different modes, with or without a specific target cluster in the URL.

## Using the AppHub UI with a target cluster

To manually deploy Helm charts into a cluster, the Application Hub provides a *Hub UI*. To access it on the Gardener Dashboard, choose \[YOUR-PROJECT\] \> *CLUSTERS* \> *External Tools* \> K8s Applications and Service Hub*.

The UI is divided into three main sections called *Applications*, *Cluster BoMs*, and *Catalog*:

* The **Applications section** displays the deployments that are currently available in the Cluster. Note that this view depends on the Namespace selection, which is located in the top-right corner of the UI. The default selection is *All Namespaces*. Deployments that are managed via Cluster BoMs are marked with a special orange badge. By clicking on this badge, it is possible to directly jump to the corresponding Cluster BoM in the UI.

* The **Cluster BoMs section** displays the list of all Cluster BoMs that currently exist for the selected cluster. By clicking on a single Cluster BoM in the list, the Cluster BoM details view is opening. In this view, the detailed state of the Cluster BoM is displayed and upgraded live via WebSockets. Users can also manually trigger a reconcile or download the Cluster BoM YAML via the UI.

* The **Catalog view** displays all Helm charts of all Repositories connected to the Application Hub in a landscape. To deploy a Helm chart, choose the respective card within the catalog, and then *Deploy* in the upper-right corner. A view is displayed, showing the values that are passed to the Helm chart during deployment. The values can be edited in this view. After clicking *Submit*, the Helm deployment starts. The *Status View* displays the current state of the deployment.

Note that deployments triggered in the UI are directly deployed to the Cluster, and therefore aren’t handled using Cluster-BoMs or added to an existing Cluster-BoM. This behaviour will change in the very near future.

## Using the AppHub UI without a target cluster

If you open the AppHub UI without a target cluster in the URL, it only shows the Hub Catalog. This allows you to browse all Helm charts of all repositories connected to the Application Hub, without being in the context of a target cluster. In this mode, you cannot perform any actions that depend on a target cluster, like deploying/updating/deleting applications. The Hub Catalog can be accessed via `<hub-base-url>/#/catalog`.
