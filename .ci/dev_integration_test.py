#!/usr/bin/env python3

import os
import subprocess

root_path = os.getcwd()
hub_path = os.environ['SOURCE_PATH']

os.environ['HUB_PATH'] = os.environ['SOURCE_PATH']
os.environ['ROOT_PATH'] = root_path

print("Run integration test on dev")

hub_kubeconfig = os.path.join(
    root_path, hub_path,
    ".ci",
    "integration_test.py"
)

command = [hub_kubeconfig, "--namespace", "release-test"]
result = subprocess.run(command)
result.check_returncode()
