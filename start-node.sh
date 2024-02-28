#!/bin/bash

# run with 
# ./start-node.sh 1 secondary 2 8083 "8081;8083"

NODE=$1
MODE=$2
NODES=$3
PORT=$4
LOCAL_PORTS=$5

CONFIG_VERSION=1 \
  FILE_PATH_PREFIX=./data \
  NODES=$NODES \
  HOSTNAME=api-$MODE-$NODE \
  LOCAL_HOSTNAME=localhost:$PORT \
  LOCAL_PORTS=$LOCAL_PORTS \
  PORT=$PORT \
  IS_LOCAL=true \
  RABBIT_MQ_HOST=localhost \
  RABBIT_MQ_PASSWORD=nuq8W2xD7Xm3lawk \
  MODE=$MODE \
  air