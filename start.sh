#!/usr/bin/env bash

DEPLOYMENT="sync-sandbox"

function usage() {
  echo -e "\nUSAGE: sh start.sh [--build, --compile, --mocks]
          --build   builds from docker image
          --compile compiles protobuf files
          --mocks   generates mock go files for testing
          "
  exit 1
}

echo -e "\n--- Starting ${DEPLOYMENT} -- "

# detect flags sent in when running start.sh
for arg in "$@"
do
  if [[ "${arg}" = "--build" ]]; then
    # build from docker image
    docker build -t sync-sandbox:1.0 .
    break;
  elif [[ "${arg}" = "--compile" ]]; then
    # compile proto files
    docker-compose up --force-recreate --remove-orphans protoc
    # TODO: have protoc recursively go through and compile
    break;
  elif [[ "${arg}" = "--mocks" ]]; then
    # generate mock files
    mockgen -source="${PWD}"/data/api.go -destination="${PWD}"/data/mocks/dao.go -package=data
    # TODO: add mock files needed for future tests
    break;
  fi
  usage
done

docker-compose up --force-recreate --remove-orphans sync-sandbox-api

echo "--- Complete ---"
