import artifactory

import os
import requests
import subprocess
import tarfile
import tempfile

def ensure_helm_binary(f: callable):
    def wrapper(helm_client: 'HelmClient', *args, **kwargs):
        helm_client._get_helm_binary_stream()
        return f(helm_client, *args, **kwargs)
    return wrapper

def test_helm_binary(execPath):
    try:
        command = [execPath, 'version']
        result = subprocess.run(command, capture_output=True, text=True)
        print(f"Test {command} with return code: {result.returncode}")
        return result.returncode == 0
    except OSError:
        return False

class HelmClient:
    def __init__(self):
        self.helm_route = 'https://get.helm.sh/helm-v3.10.3-linux-amd64.tar.gz'
        self.bin_path = 'helm'
        if not test_helm_binary(self.bin_path):
            tempdir = tempfile.gettempdir()
            print(f"helm not found in path, installing it to {tempdir}")
            self.bin_path = f"{tempdir}/helm"

    def _get_helm_binary_stream(self):
        if os.path.isabs(self.bin_path) and not os.path.isfile(self.bin_path):
            res = requests.get(self.helm_route, stream=True)
            with tarfile.open(fileobj=res.raw, mode='r|*') as tar:
                res.raw.seekable = False
                for member in tar:
                    if not member.name == 'linux-amd64/helm':
                        continue

                    fileobj = tar.extractfile(member)
                    with open(self.bin_path, "wb") as outfile:
                        outfile.write(fileobj.read())
                    os.chmod(self.bin_path, 744)

    @ensure_helm_binary
    def repo_add(self, repo_name: str, repo_url: str):
        command = [self.bin_path, "repo", "add", repo_name, repo_url]
        result = subprocess.run(command)
        if result.returncode != 0:
            raise RuntimeError("Could not add helm repo")

    @ensure_helm_binary
    def dependency_update(self, chart_path: str):
        self.repo_add(repo_name="bitnami", repo_url="https://charts.bitnami.com/bitnami")
        command = [self.bin_path, "dependency", "build", chart_path]
        result = subprocess.run(command)
        if result.returncode != 0:
            raise RuntimeError("Could not run helm dependency build")

    @ensure_helm_binary
    def package_chart(self, chart_path: str, out_path: str):
        self.dependency_update(chart_path=chart_path)
        command = [self.bin_path, "package", chart_path, "-d", out_path]
        result = subprocess.run(command)
        if result.returncode != 0:
            raise RuntimeError("Could not package helm chart")

    def upload_chart(self, repo_client: artifactory.ArtifactoryRepoClient, artifact, chart_name: str):
        return repo_client.upload_artifact(artifact_name=chart_name, artifact=artifact)

    @ensure_helm_binary
    def deploy_chart(self, host: str, username: str, password: str, release_name: str, chart_name: str,
                     namespace: str, version: str, kubeconfig_path: str, helm_values_path: str):
        if version:
            command = [self.bin_path, "upgrade", release_name, chart_name, "--namespace", namespace, "--kubeconfig",
                       kubeconfig_path, "--version", version, "--repo", host,
                       "--username", username, "--password", password, "-f", helm_values_path, "-i"]
        else:
            command = [self.bin_path, "upgrade", release_name, chart_name, "--namespace", namespace, "--kubeconfig",
                       kubeconfig_path, "--repo", host, "--username", username,
                       "--password", password, "-f", helm_values_path, "-i"]

        print(f"  Run helm upgrade {release_name} {chart_name} --namespace {namespace} --version {version if version else '<None>'}"
              f" --kubeconfig {kubeconfig_path} --repo {host} --username {username} --password *** -f {helm_values_path} -i\n")
        result = subprocess.run(command, capture_output=True, text=True)
        print(result.stdout)
        print(result.stderr)
        if result.returncode != 0:
            raise RuntimeError("Could not upgrade helm chart on cluster")

    @ensure_helm_binary
    def undeploy_chart(self, release_name: str, namespace: str, kubeconfig_path: str):
        command = [self.bin_path, "get", "manifest", release_name, "--namespace", namespace, "--kubeconfig",
                    kubeconfig_path]
        print(f"  Run: {' '.join(command)}\n")
        result = subprocess.run(command, capture_output=True, text=True)
        if result.returncode != 0:
            print(result.stdout)
            print(result.stderr)
            print("No hub-controller release found nothi ng to do")
            return

        command = [self.bin_path, "uninstall", release_name, "--namespace", namespace, "--kubeconfig",
                    kubeconfig_path]

        print(f"  Run: {' '.join(command)}\n")
        result = subprocess.run(command, capture_output=True, text=True)
        print(result.stdout)
        print(result.stderr)
        if result.returncode != 0:
            raise RuntimeError("Could not uninstall helm chart on cluster")

    @ensure_helm_binary
    def helm_index(self, chart_dir: str, index_yaml_file: str):
        if index_yaml_file:
            command = [self.bin_path, "repo", "index", "--merge", index_yaml_file, chart_dir]
        else:
            command = [self.bin_path, "repo", "index", chart_dir]

        print(f"  Run: {' '.join(command)}\n")
        result = subprocess.run(command, capture_output=True, text=True)
        print(result.stdout)
        print(result.stderr)
        if result.returncode != 0:
            raise RuntimeError("Failed to update index.yaml")
