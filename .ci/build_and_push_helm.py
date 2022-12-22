#!/usr/bin/env python3

import artifactory
import helm
import secret_server
import shutil
import tempfile

import os
import subprocess
import utils

from google.cloud import storage
from ci.util import ctx

def uploadFolder(packagedChartsDir, bucketname):
    print("==== Upload directory", packagedChartsDir, "to GCS ====")
    files = [f for f in os.listdir(packagedChartsDir) if os.path.isfile(os.path.join(packagedChartsDir, f))]

    client = storage.Client()
    for f in files:
        print(f"Uploading file {f}")
        bucket = client.get_bucket(bucketname)
        testfile = bucket.blob(f)
        testfile.upload_from_filename(os.path.join(packagedChartsDir, f))
    
    print(f"Uploaded to bucket {bucketname}")

def runHelmIndex(packagedChartsDir, bucketname):
    """ Build index.yaml for the directory *packagedChartsDir* and write the index file into this directory. """
    print("==== Building index.yaml for directory", packagedChartsDir, "====")
    indexFile = packagedChartsDir + os.sep + "index.yaml"

    client = storage.Client()
    bucket = client.get_bucket(bucketname)
    blob = bucket.blob("index.yaml")

    if blob.exists():
        blob.download_to_filename(indexFile)
    else:
        indexFile = None

    helm_client.helm_index(packagedChartsDir, indexFile)

def prepare_gcs_credentials(gcs_credentials_cfg_name, credentials_file):
    print("==== Reading GCS credentials from cc-config ====")
    cfg_factory = ctx().cfg_factory()
    cfg = cfg_factory._cfg_element(cfg_type_name="gcs", cfg_name=gcs_credentials_cfg_name)
    gcs_credentials = cfg.raw["password"]
    credentials_file.write(gcs_credentials)
    os.environ["GOOGLE_APPLICATION_CREDENTIALS"] = os.path.abspath(credentials_file.name)

def packageChartRepo(chart_dir, bucketname):
    print("== Packaging Repo Directory", chart_dir, "==")
    runHelmIndex(chart_dir, bucketname)
    uploadFolder(chart_dir, bucketname)


if __name__ == "__main__":
    print("\n ===== Starting Helm Packaging - Python Powered =====")

    bucketname = "potter_charts"
    credentials_file_name = "gcs_credentials.json.tmp"

    source_path = os.environ['SOURCE_PATH']
    version_path = os.environ['VERSION_PATH']
    if os.environ.get('HELM_CHART_PATH'):
        pipeline_out_path =  os.environ['HELM_CHART_PATH']
    else:
        print(f"Environment: {os.environ}")
        pipeline_out_path = None

    secret_server_client = secret_server.SecretServer('hub')
    helm_client = helm.HelmClient()

    chart_version = utils.get_chart_version(source_path, version_path)

    chart_name = "k8s-potter-hub"
    chart_path = os.path.join(source_path, 'chart', 'hub')

    helm_client.repo_add("potter-charts", "https://potter-charts.storage.googleapis.com")

    # Now override the place holders in helm chart with concrete values
    # To not override our source dir copy chart directory to temp dir

    # Note use this for debugging if you want to preserve the temp
    # directory: temp_out = mkdtemp(prefix="helm_chart_")
    with tempfile.TemporaryDirectory(prefix="helm_chart_") as temp_out:
        chart_path_out = os.path.join(temp_out, 'hubchart')
        print(f"Rendering helm chart in {chart_path_out}")
        shutil.copytree(chart_path, chart_path_out)

        # get the image version from file
        image_version_file = version_path + "/version"
        with open(image_version_file) as image_file:
            image_version = image_file.read()

        utils.replace_chart_placeholder(
            chart_path=chart_path_out,
            image_version=image_version,
            chart_version=chart_version,
            chart_name=chart_name
        )

        # helm_client.dependency_update(chart_path=chart_path_out)
        print("package chart")
        helm_client.package_chart(chart_path=chart_path_out, out_path=temp_out)
        tgz_name = chart_name + "-" + chart_version + ".tgz"
        chart_tgz_path = os.path.join(temp_out, tgz_name)

        with utils.TempFileAuto(prefix="gcs_credentials.json_") as cred_file:
            prepare_gcs_credentials("hub-chart-pipeline-stable", cred_file)
            cred_file.switch()

            packageChartRepo(temp_out, "potter-charts")

        if pipeline_out_path:
            print(f"Copying helm chart {chart_tgz_path} to pipeline-out dir {pipeline_out_path}")
            shutil.copy(chart_tgz_path, pipeline_out_path)

    print("\n ===== Finished Helm Packaging - Python Over =====")
