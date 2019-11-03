#!/bin/sh
set -eou pipefail
mvn clean package
cp ${SRC_PKG}/target/*-runner.jar ${DEPLOY_PKG}
