#!/usr/bin/env bash

set -x
set -e
set -o pipefail

curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.27.0

mv bin/golangci-lint $SOURCE_PATH

export ROOT_PATH="$(readlink -f "${SOURCE_PATH}/..")"

cd $SOURCE_PATH

./golangci-lint run -c .golangci.yml --timeout 10m0s | tee $ROOT_PATH/$LINT_PATH/out
if [[ $? == 1 ]]; then
  echo "Linter exited with errors" | tee -a $ROOT_PATH/$LINT_PATH/out
else
  echo "Linter exited without errors" | tee -a $ROOT_PATH/$LINT_PATH/out
fi