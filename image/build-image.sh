#!/bin/bash
IMAGE_DIR=$(pwd $(dirname $0))
if [ ! -d $IMAGE_DIR/bin ]; then
  mkdir $IMAGE_DIR/bin
fi
cp $IMAGE_DIR/../../../output/go/bin/simple-deploy $IMAGE_DIR/bin
DOCKER_USER=cewong
sudo docker build -t $DOCKER_USER/simple-deploy .
sudo docker push $DOCKER_USER/simple-deploy
