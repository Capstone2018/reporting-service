#!/usr/bin/env bash

# you should `source` this script before building the
# Docker container image or trying to run it. Use this
# command to source this into your current shell:
#   source devenv.sh

export NETWORK_NAME=crapstone

# docker container image name for reporting service
export REPORTING_SERVICE_IMAGE=aethan/reporting-service

# docker container image name for our customized MySQL image
export REPORTING_POSTGRES_IMAGE=aethan/postgresreports

# database name in which our schema will be created
export POSTGRES_DB=reporting

export POSTGRES_USER=admin
# random MySQL root password 
# use $(openssl rand -base64 18) in prod
export POSTGRES_PASSWORD=supersecret 

export POSTGRES_HOST=devpsql
export POSTGRES_PORT=5432