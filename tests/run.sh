#!/bin/bash
cd $(dirname "$0") && go test ./... -v -failfast
