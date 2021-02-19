import kube.ctx
from util import ctx

import base64
from enum import Enum
import json
import os
import subprocess
import tarfile
import tempfile
import urllib.request
import yaml


class EfkConsts(Enum):
    efk_namepace = 'efk'
    efk_release_name = 'efk'


class NginxConsts(Enum):
    chart_name = 'ingress-nginx'
    chart_version = '2.11.1'
    namespace = 'ingress-nginx'
    release_name = 'ingress-nginx'


class HubConsts(Enum):
    chart_name = "k8s-potter-hub"
    namespace = 'hub'
    release_name = 'hub'


class ChartRepoCredentials:
    def __init__(self, host: str, username: str, password: str):
        self.host = host
        self.username = username
        self.password = password


def get_helm3():
    helm_url = 'https://get.helm.sh/helm-v3.1.2-linux-amd64.tar.gz'
    file_name = helm_url.split('/')[-1]
    urllib.request.urlretrieve(helm_url, file_name)

    if file_name.endswith("tar.gz"):
        with tarfile.open(file_name) as tar:
            member = tar.getmember('linux-amd64/helm')
            tar.fileobj.seek(member.offset_data)

            with open("helm", "wb") as outfile:
                outfile.write(tar.fileobj.read(member.size))

    os.chmod("helm", 744)


def read_data(path: str):
    with open(path, "r") as file:
        return file.read()


def write_data(path: str, data: str):
    with open(path, "w") as file:
        file.write(data)


def helm_repo_add(repo_name: str, repo_url: str):
    command = ["./helm", "repo", "add", repo_name, repo_url]
    result = subprocess.run(command)
    if result.returncode != 0:
        raise Exception("Could not add helm repo")


def helm_dependency_update(chart_path: str):
    command = ["./helm", "dependency", "build", chart_path]
    result = subprocess.run(command)
    if result.returncode != 0:
        raise Exception("Could not run helm dependency build")


def get_kubeconfig(cfg_name: str):
    factory = ctx().cfg_factory().cfg_set("hub")
    return yaml.dump(factory.kubernetes(cfg_name).kubeconfig())


def get_landscape_config(cfg_name: str):
    factory = ctx().cfg_factory().cfg_set("hub")
    return factory.hub(cfg_name)


def get_container_registry_creds(cfg_name: str):
    factory = ctx().cfg_factory().cfg_set("hub")
    return factory.container_registry(cfg_name).credentials()


def get_chart_version(source_path, version_path):
    if is_release_job():
        version_file = os.path.join(version_path, "version")
        with open(version_file) as f:
            chart_version = f.readline().rstrip()
        return chart_version
    else:
        chart_version_file = os.path.join(source_path, "VERSION")
        with open(chart_version_file) as f:
            chart_version = f.readline().rstrip()
        return chart_version


def replace_placeholder_in_yaml(file_name, placeholder, value):
    with open(file_name) as f:
        modified_file = f.read().replace(placeholder, value)
    with open(file_name, "w") as f:
        f.write(modified_file)


def is_release_job():
    job_name = os.environ['BUILD_JOB_NAME']
    if "release" in job_name:
        return True
    else:
        return False


def replace_chart_placeholder(chart_path: str, image_version: str, chart_version: str, chart_name: str):

    values_yaml = chart_path + "/values.yaml"
    chart_yaml = chart_path + "/Chart.yaml"

    image_registry = "eu.gcr.io"
    image_repo = "sap-gcp-cp-k8s-stable-hub/potter/"
    image_repo_3rd = "sap-gcp-cp-k8s-stable-hub/3rd/"

    replace_placeholder_in_yaml(values_yaml, "<UI_BACKEND_REGISTRY>", image_registry)
    replace_placeholder_in_yaml(values_yaml, "<UI_BACKEND_REPO>", image_repo + "ui-backend")
    replace_placeholder_in_yaml(values_yaml, "<UI_BACKEND_TAG>", image_version)

    replace_placeholder_in_yaml(values_yaml, "<APPREPOSITORY_REGISTRY>", image_registry)
    replace_placeholder_in_yaml(values_yaml, "<APPREPOSITORY_CONTROLLER_REPO>", image_repo + "apprepo")
    replace_placeholder_in_yaml(values_yaml, "<APPREPOSITORY_CONTROLLER_TAG>", image_version)

    replace_placeholder_in_yaml(values_yaml, "<DASHBOARD_REGISTRY>", image_registry)
    replace_placeholder_in_yaml(values_yaml, "<DASHBOARD_REPO>", image_repo + "dashboard")
    replace_placeholder_in_yaml(values_yaml, "<DASHBOARD_TAG>", image_version)
    
    replace_placeholder_in_yaml(values_yaml, "<HUB_TOKEN_BUTLER_REGISTRY>", image_registry)
    replace_placeholder_in_yaml(values_yaml, "<HUB_TOKEN_BUTLER_REPOSITORY>", image_repo + "hubtokenbutler")
    replace_placeholder_in_yaml(values_yaml, "<HUB_TOKEN_BUTLER_TAG>", image_version)

    replace_placeholder_in_yaml(values_yaml, "<KUBECTL_REGISTRY>", image_repo + "kubectl")
    replace_placeholder_in_yaml(values_yaml, "<KUBECTL_REPOSITORY>", image_repo + "kubectl")
    replace_placeholder_in_yaml(values_yaml, "<KUBECTL_TAG>", image_version)

    # 3rd party images
    replace_placeholder_in_yaml(values_yaml, "<AUTH_PROXY_REGISTRY>", image_registry)
    replace_placeholder_in_yaml(values_yaml, "<AUTH_PROXY_REPOSITORY>", image_repo_3rd + "keycloak-gatekeeper")
    replace_placeholder_in_yaml(values_yaml, "<AUTH_PROXY_TAG>", "2.3.0-scratch-r4")

    replace_placeholder_in_yaml(values_yaml, "<CHART_SERVICE_REGISTRY>", image_registry)
    replace_placeholder_in_yaml(values_yaml, "<CHART_SERVICE_REPOSITORY>", image_repo_3rd + "kubeapps-chartsvc")
    replace_placeholder_in_yaml(values_yaml, "<CHART_SERVICE_TAG>", "1.9.0-r0") 
    
    # for public repo check latest version at: https://github.com/bitnami/charts/blob/master/bitnami/mongodb/values.yaml
    replace_placeholder_in_yaml(values_yaml, "<MONGODB_REGISTRY>", "docker.io")
    replace_placeholder_in_yaml(values_yaml, "<MONGODB_REPOSITORY>", "bitnami/mongodb")
    replace_placeholder_in_yaml(values_yaml, "<MONGODB_TAG>", "4.4.3-debian-10-r21")

    replace_placeholder_in_yaml(values_yaml, "<MONGODB_EXPORTER_REGISTRY>", "eu.gcr.io")
    replace_placeholder_in_yaml(values_yaml, "<MONGODB_EXPORTER_REPOSITORY>", image_repo_3rd + "mongodb-exporter")
    replace_placeholder_in_yaml(values_yaml, "<MONGODB_EXPORTER_TAG>", "0.10.0-debian-10-r69")

    replace_placeholder_in_yaml(values_yaml, "<MINIDEB_REGISTRY>", "eu.gcr.io")
    replace_placeholder_in_yaml(values_yaml, "<MINIDEB_REPOSITORY>", image_repo_3rd + "minideb")
    replace_placeholder_in_yaml(values_yaml, "<MINIDEB_TAG>", "buster")

    replace_placeholder_in_yaml(values_yaml, "<NGINX_REGISTRY>", "eu.gcr.io")
    replace_placeholder_in_yaml(values_yaml, "<NGINX_REPOSITORY>", image_repo_3rd + "nginx")
    replace_placeholder_in_yaml(values_yaml, "<NGINX_TAG>", "1.19.2-alpine")

    replace_placeholder_in_yaml(values_yaml, "<IMAGE_PULL_POLICY>", "IfNotPresent")

    replace_placeholder_in_yaml(chart_yaml, "<CHART_VERSION>", chart_version)
    replace_placeholder_in_yaml(chart_yaml, "<CHART_NAME>", chart_name)

class TempFileAuto(object):
    def __init__(self, prefix=None, mode='w+', suffix=".yaml"):
        self.file_obj = tempfile.NamedTemporaryFile(mode=mode, prefix=prefix, suffix=suffix, delete=False)
        self.name = self.file_obj.name
    def __enter__(self):
        return self
    def write(self, b):
        self.file_obj.write(b)
    def writelines(self, lines):
        self.file_obj.writelines(lines)
    def switch(self):
        self.file_obj.close()
        return self.file_obj.name
    def __exit__(self, type, value, traceback):
        if not self.file_obj.closed:
            self.file_obj.close()
        os.remove(self.file_obj.name)
        return False