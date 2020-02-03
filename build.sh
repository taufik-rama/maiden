#!/bin/bash
# Use this to build the application so we could
# get the metadata

GO_VERSION_OUT="$(go version)"
GO_VERSION_OUT=($GO_VERSION_OUT)
GIT_LOG_OUT="$(git log -n 1 | head -n 1)"
GIT_LOG_OUT=($GIT_LOG_OUT)

MAIDEN_VERSION=$(cat VERSION)
GO_VERSION=${GO_VERSION_OUT[2]}
LATEST_COMMIT=$(echo ${GIT_LOG_OUT[1]} | cut -c 1-9)
BUILD_DATE="$(date -uR)"
BUILD_OS=${GO_VERSION_OUT[3]}

BIN=${BUILD_OS}
BIN=(${BIN//\// })

go build -ldflags "\
-X 'main.maidenVersion=${MAIDEN_VERSION}' \
-X 'main.goVersion=${GO_VERSION}' \
-X 'main.latestCommit=${LATEST_COMMIT}' \
-X 'main.buildDate=${BUILD_DATE}' \
-X 'main.buildOS=${BUILD_OS}' \
-s -w" \
-o "maiden-${BIN[0]}-${BIN[1]}" ./cmd/maiden
