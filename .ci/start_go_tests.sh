#!/usr/bin/env bash

set -x
set -e
set -o pipefail

export ROOT_PATH="$(readlink -f "${SOURCE_PATH}/..")"
cd $SOURCE_PATH


export GO111MODULE=on

go test ./... | tee $ROOT_PATH/$BACKEND_TEST_PATH/out