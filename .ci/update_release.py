#!/usr/bin/env python3

import pathlib
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
