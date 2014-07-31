#!/bin/bash
EXAMPLE_SOURCE=$(readlink -f $(dirname $0))
IMAGE_DIR=$EXAMPLE_SOURCE/image
mkdir -p $IMAGE_DIR/bin
cp $EXAMPLE_SOURCE/output/simple-deploy $IMAGE_DIR/bin
cd $IMAGE_DIR
DOCKER_USER=cewong
sudo docker build -t $DOCKER_USER/simple-deploy .
