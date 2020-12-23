#!/usr/bin/env python3

import os
from subprocess import Popen, PIPE, STDOUT, run
import sys
import utils
import yaml

from util import ctx

print("Starting integration test")
hub_path = os.environ['HUB_PATH']
root_path = os.environ["ROOT_PATH"]
landscape = os.environ["LANDSCAPE"]
namespace = os.environ["NAMESPACE"]

try:
    # env var is implicitly set by the output dir in case of a release job
    integration_test_path = os.environ["INTEGRATION_TEST_PATH"]
except KeyError:
    print("Integration test env var not set")

factory = ctx().cfg_factory()
landscape_kubeconfig = factory.kubernetes("hub-" + landscape)
landscape_test_kubeconfig = factory.kubernetes("hub-" + landscape + "-test")

landscape_kubeconfig_name = "landscape_kubeconfig"
landscape_kubeconfig_path = os.path.join(root_path, hub_path,
                                         "integration-test",
                                         landscape_kubeconfig_name)

landscape_test_kubeconfig_name = "landscape_test_kubeconfig"
landscape_test_kubeconfig_path = os.path.join(root_path, hub_path,
                                              "integration-test",
                                              landscape_test_kubeconfig_name)

utils.write_data(landscape_kubeconfig_path, yaml.dump(
                landscape_kubeconfig.kubeconfig()))
utils.write_data(landscape_test_kubeconfig_path, yaml.dump(
                landscape_test_kubeconfig.kubeconfig()))

landscape_config = utils.get_landscape_config("hub-" + landscape)
int_test_config = landscape_config.raw["int-test"]["config"]
token = int_test_config["auth"]["token"]

os.chdir(os.path.join(root_path, hub_path, "integration-test"))

command = ["go", "run", "main.go",
           "--kubeconfig", landscape_kubeconfig_name,
           "--target-kubeconfig", landscape_test_kubeconfig_name,
           "--token", token, "--namespace", namespace]

try:
    # check if path var is set
    integration_test_path
except NameError:
    run = run(command)
else:
    output_path = os.path.join(root_path, integration_test_path, "out")

    with Popen(command, stdout=PIPE, stderr=STDOUT, bufsize=1, universal_newlines=True) as run, open(output_path, 'w') as file:
        for line in run.stdout:
            sys.stdout.write(line)
            file.write(line)

if run.returncode != 0:
    raise EnvironmentError("Integration test exited with errors")