#!/usr/bin/env bash

DEPLOYMENT="sync-sandbox"

function usage() {
  echo -e "\nUSAGE: sh start.sh [--compile]
          --compile compiles protobuf files
          "
  exit
}
echo "--- Starting ${sync-sandbox} -- "

# detect flags sent in when running start.sh
for arg in "$@"
do
  if [[ "${arg}" = "--compile" ]]
  then
    docker-compose up --force-recreate --remove-orphans -d protoc
    break;
  fi
  usage
done

echo "--- Complete ---"
