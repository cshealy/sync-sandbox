#!/usr/bin/env bash

DEPLOYMENT="sync-sandbox"

function usage() {
  echo -e "\nUSAGE: sh start.sh [--build, --compile, --mocks]
          --build   builds from docker image
          --compile compiles protobuf files
          --mocks   generates mock go files for testing
          --swagger serves swagger documentation
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
    docker-compose up --force-recreate --remove-orphans gomock
    # TODO: add mock files needed for future tests
    break;
  elif [[ "${arg}" = "--swagger" ]]; then
    # serve swagger docs
    docker-compose up -d --force-recreate --remove-orphans swagger-docs
    echo "--- Serving swagger docs on localhost:8081/swagger ---"
    # TODO: serve swagger docs on API server using a mux
    break;
  fi
  usage
done

docker-compose up --force-recreate --remove-orphans sync-sandbox-api

echo "--- Complete ---"
