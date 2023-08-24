#!/usr/bin/env bash

docker compose --env-file ./.env.local down -v
docker rmi $(docker images --format '{{.Repository}}:{{.Tag}}' | grep 'transactions-proxy-backend')
docker compose --env-file ./.env.local up -d
