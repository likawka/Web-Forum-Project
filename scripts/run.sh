#!/bin/bash
clear
DOCKER_IMAGE_NAME="app"
HOST_PORT=8080
CONTAINER_PORT=8080
docker build . -t $DOCKER_IMAGE_NAME --load
docker run -it -p $HOST_PORT:$CONTAINER_PORT $DOCKER_IMAGE_NAME
docker ps -a
