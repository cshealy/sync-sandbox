#!/usr/bin/env bash

DEPLOYMENT="sync-sandbox"

function usage() {
  echo -e "\nUSAGE: sh start.sh [--compile, --mocks]
          --compile compiles protobuf files
          --mocks   generates mock go files for testing
          "
  exit 1
}

echo -e "\n--- Starting ${sync-sandbox} -- "

# check for no args
if [ $# -eq 0 ]; then
    usage
    exit 1
fi

# detect flags sent in when running start.sh
for arg in "$@"
do
  if [[ "${arg}" = "--compile" ]]; then
    # compile proto files
    docker-compose up --force-recreate --remove-orphans -d protoc
    # TODO: have protoc recusively go through and compile
    break;
  elif [[ "${arg}" = "--mocks" ]]; then
    # generate mock files
    mockgen -source="${PWD}"/data/api.go -destination="${PWD}"/data/mocks/dao.go -package=data
    # TODO: add mock files needed for future tests
    break;
  fi
  usage
done

echo "--- Complete ---"
