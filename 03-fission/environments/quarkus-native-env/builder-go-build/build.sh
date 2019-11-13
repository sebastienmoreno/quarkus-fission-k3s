#!/bin/sh
#set -eou pipefail
mvn package -DskipTests=true -Pnative
cp ${SRC_PKG}/target/*-runner ${DEPLOY_PKG}
