#!/bin/ash

set -x
set -e
set -o pipefail

export ROOT_PATH="$(readlink -f "${SOURCE_PATH}/..")"

yarn --cwd $SOURCE_PATH/dashboard/ install --frozen-lockfile 2>&1 | tee $ROOT_PATH/$FRONTEND_TEST_PATH/out
yarn --cwd=$SOURCE_PATH/dashboard run lint 2>&1 | tee -a $ROOT_PATH/$FRONTEND_TEST_PATH/out
CI=true yarn --cwd $SOURCE_PATH/dashboard/ run test 2>&1 | tee -a $ROOT_PATH/$FRONTEND_TEST_PATH/out