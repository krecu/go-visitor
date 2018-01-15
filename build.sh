#!/usr/bin/env bash

APP_ENV=$1
APP_ROOT_PATH=/app
APP_DB_PATH=/databases
APP_CPU=4
APP_AEROSPIKE_NODE_1=192.168.0.2
APP_AEROSPIKE_NODE_2=192.168.0.5
APP_GRAYLOG=192.168.0.48
APP_WEB_ADDR=:8090
APP_GRPC_ADDR=:8091
APP_DEBUG_LEVEL=7
APP_CONTAINER_NAME=videonow/visitor:$APP_ENV

rm -rf ./build/visitor
rm -rf ./build/config.yaml

GOOS=linux GOARCH=amd64 go build -o ./build/visitor ./app
sed -e "s&{APP_ROOT_PATH}&$APP_ROOT_PATH&g" \
    -e "s&{APP_DB_PATH}&$APP_DB_PATH&g" \
    -e "s&{APP_ENV}&$APP_ENV&g" \
    -e "s&{APP_AEROSPIKE_NODE_1}&$APP_AEROSPIKE_NODE_1&g" \
    -e "s&{APP_AEROSPIKE_NODE_2}&$APP_AEROSPIKE_NODE_2&g" \
    -e "s&{APP_GRAYLOG}&$APP_GRAYLOG&g" \
    -e "s&{APP_WEB_ADDR}&$APP_WEB_ADDR&g" \
    -e "s&{APP_GRPC_ADDR}&$APP_GRPC_ADDR&g" \
    -e "s&{APP_DEBUG_LEVEL}&$APP_DEBUG_LEVEL&g" \
    -e "s&{APP_CPU}&$APP_CPU&g" \
    ./app/config.default.yaml >> ./build/config.yaml
docker build -t $APP_CONTAINER_NAME .