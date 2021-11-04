#!/usr/bin/env bash

DOCKER_COMPOSE="ci/assets/docker-compose.yml"

docker-compose -f "$DOCKER_COMPOSE" -p "$GITHUB_SHA" down
