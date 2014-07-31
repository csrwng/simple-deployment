#!/bin/bash
KUBERNETES_SOURCE=../../GoogleCloudPlatform/kubernetes
KUBERNETES_SOURCE=$(readlink -f $KUBERNETES_SOURCE)
EXAMPLE_SOURCE=$(readlink -f $(dirname $0))
export GOPATH=${EXAMPLE_SOURCE}/cmd:${KUBERNETES_SOURCE}/output/go:${KUBERNETES_SOURCE}/third_party
mkdir -p $EXAMPLE_SOURCE/output
cd $EXAMPLE_SOURCE/output
go build ../cmd/simple-deploy
