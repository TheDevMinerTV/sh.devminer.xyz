---
description: Simple wrapper around the compose.yml files for different environments
tags: [docker]
---
#!/usr/bin/env bash

# Simple wrapper around the compose.yml files for different environments
#
# Usage: ./compose.sh <environment> [compose command...]
#
# Author: DevMiner <devminer@devminer.xyz>

VALID_ENVS=("dev" "prod")

if [ -z $1 ]; then
  echo "Environments:" >&2
  for ENV in ${VALID_ENVS[@]}; do
    echo "- $ENV" >&2
  done
  exit 1
fi

ENV="$1"
shift

if [[ ! " ${VALID_ENVS[*]} " =~ [[:space:]]${ENV}[[:space:]] ]]; then
  echo "ERROR: Invalid environment specified: \"${ENV}\"" >&2
  exit 1
fi

FILES=("-f compose.yml" "-f compose.$ENV.yml")

set -x
docker compose ${FILES[*]} $@
