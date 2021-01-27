#!/usr/bin/env python3

import pathlib
import re
import util
import os

from github.util import GitHubRepositoryHelper

OUTPUT_FILE_NAME = 'out'
VERSION_FILE_NAME = 'VERSION'

repo_owner_and_name = util.check_env('SOURCE_GITHUB_REPO_OWNER_AND_NAME')
repo_dir = util.check_env('MAIN_REPO_DIR')

lint_path = util.check_env('LINT_PATH')
backend_test_path = util.check_env('BACKEND_TEST_PATH')
frontend_test_path = util.check_env('FRONTEND_TEST_PATH')
helm_chart_path = os.environ.get('HELM_CHART_PATH')

lint_path = pathlib.Path(lint_path).resolve()
backend_test_path = pathlib.Path(backend_test_path).resolve()
frontend_test_path = pathlib.Path(frontend_test_path).resolve()

repo_owner, repo_name = repo_owner_and_name.split('/')

repo_path = pathlib.Path(repo_dir).resolve()

lint_path = lint_path / OUTPUT_FILE_NAME
backend_test_path = backend_test_path / OUTPUT_FILE_NAME
frontend_test_path = frontend_test_path / OUTPUT_FILE_NAME

version_file_path = repo_path / VERSION_FILE_NAME

version_file_contents = version_file_path.read_text()

cfg_factory = util.ctx().cfg_factory()
github_cfg = cfg_factory.github('github_com')

github_repo_helper = GitHubRepositoryHelper(
    owner=repo_owner,
    name=repo_name,
    github_cfg=github_cfg,
)

gh_release = github_repo_helper.repository.release_from_tag(version_file_contents)

gh_release.upload_asset(
    content_type='text/plain',
    name=f'linting-result-{version_file_contents}.txt',
    asset=lint_path.open(mode='rb'),
)
gh_release.upload_asset(
    content_type='text/plain',
    name=f'backend-test-result-{version_file_contents}.txt',
    asset=backend_test_path.open(mode='rb'),
)
gh_release.upload_asset(
    content_type='text/plain',
    name=f'frontend-test-result-{version_file_contents}.txt',
    asset=frontend_test_path.open(mode='rb'),
)
try:
    os.environ['INTEGRATION_TEST_PATH']
except KeyError:
    print("No integration test output path found. Output will not be added to release")
else:
    integration_test_path = util.check_env('INTEGRATION_TEST_PATH')
    integration_test_path = pathlib.Path(integration_test_path).resolve()
    integration_test_path = integration_test_path / OUTPUT_FILE_NAME
    gh_release.upload_asset(
        content_type='text/plain',
        name=f'integration-test-result-{version_file_contents}.txt',
        asset=integration_test_path.open(mode='rb'),
    )

if helm_chart_path:
    files = os.listdir(helm_chart_path)
    print(f"Found helm_chart_path with content: {files}")
    helm_chart_name = f"k8s-potter-hub-{version_file_contents}.tgz"
    helm_tgz_path = os.path.join(helm_chart_path, helm_chart_name)
    if os.path.exists(helm_tgz_path):
        with open(helm_tgz_path, mode='rb') as asset:
            gh_release.upload_asset(
                content_type='application/x-tar',
                name=helm_chart_name,
                asset=asset,
            )
    else:
        print(f"Warning: Helm Chart not found: {helm_tgz_path}")

# Update description of the release notes
release_notes = gh_release.body
if not release_notes:
    release_notes = "n/a"

description_md = f"""# Release {version_file_contents} of the potter-hub
This is release {version_file_contents} of the Gardener Potter-Hub project.
The artifacts of this release are used in deployments of the
[potter-hub](https://github.com/gardener/potter-hub) project. Potter is distributed
and installed via Helm. A helm chart is available [here](https://storage.googleapis.com/potter-charts/k8s-potter-hub-{version_file_contents}.tgz).
You can add the corresponding helam chart repository to helm with the following command:

```
 helm repo add potter https://storage.googleapis.com/potter-charts

```

To list the content use

```
helm search repo potter

```

To get chart information use:

```
helm show chart potter/k8s-potter-hub


```

For detailed installation instructions, visit the
Helm Chart's [README.md](https://github.com/gardener/potter-hub/blob/{version_file_contents}/chart/hub/README.md).
To make the most out of a Potter Deployment, you should also take a
look at the [README.md](https://github.com/gardener/potter-controller/blob/main/chart/hub-controller/Readme.md)
of the corresponding Potter-Controller.

```
helm show readme potter/k8s-potter-hub

```

## Release Notes

{release_notes}
"""

# Github behaves a bit strange with respect to line endings, so replace all '\n' with spaces but
# preserve double '\n'
description_md = description_md.replace("\n\n", " $$$###$$$ ")
description_md = description_md.replace("\n", " ")
description_md = description_md.replace(" $$$###$$$ ", "\n\n")


gh_release.edit(body=description_md)