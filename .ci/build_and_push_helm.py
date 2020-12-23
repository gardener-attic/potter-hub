#!/usr/bin/env python3

import artifactory
import helm
import secret_server
import shutil
import tempfile

import os
import utils

if __name__ == "__main__":
    print("\n ===== Starting Helm Packaging - Python Powered =====")

    source_path = os.environ['SOURCE_PATH']
    version_path = os.environ['VERSION_PATH']

    secret_server_client = secret_server.SecretServer('hub')
    helm_client = helm.HelmClient()

    art_creds = secret_server_client.get_container_registry_creds("artifactory-readwrite-hub")
    artifactory_client = artifactory.ArtifactoryRepoClient(
        base_url=art_creds.host(),
        user=art_creds.username(),
        password=art_creds.passwd()
    )

    chart_version = utils.get_chart_version(source_path, version_path)

    chart_name = "k8s-potter-hub"
    chart_path = os.path.join(source_path, 'chart', 'hub')

    helm_client.repo_add("bitnami", "https://charts.bitnami.com/bitnami")

    # Now override the place holders in helm chart with concrete values
    # To not override our source dir copy chart directory to temp dir
    with tempfile.TemporaryDirectory(prefix="helm_chart_") as temp_out:
        chart_path_out = os.path.join(temp_out, 'hubchart')
        print(f"Rendering helm chart in {chart_path_out}")
        shutil.copytree(chart_path, chart_path_out)

        utils.replace_chart_placeholder(
            chart_path=chart_path_out,
            version_path=version_path,
            chart_version=chart_version,
            chart_name=chart_name
        )

        helm_client.dependency_update(chart_path=chart_path_out)
        helm_client.package_chart(chart_path=chart_path_out, out_path=temp_out)
        tgz_name = chart_name + "-" + chart_version + ".tgz"
        chart_tgz_path = os.path.join(temp_out, tgz_name)

        with open(file=chart_tgz_path, mode='rb') as artifact:
            helm_client.upload_chart(
                repo_client=artifactory_client,
                artifact=artifact,
                chart_name=tgz_name
            )

    print("\n ===== Finished Helm Packaging - Python Over =====")
