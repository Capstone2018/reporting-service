#!/usr/bin/env bash

# you should `source` this script before building the
# Docker container image or trying to run it. Use this
# command to source this into your current shell:
#   source devenv.sh

export NETWORK_NAME=crapstone

# docker container image name for reporting service
export REPORTING_SERVICE_IMAGE=aethanol/reporting-service

# docker container image name for our customized MySQL image
export REPORTING_MYSQL_IMAGE=aethanol/mysqlreports

# database name in which our schema will be created
export MYSQL_DATABASE=reporting

# random MySQL root password 
# use $(openssl rand -base64 18) in prod
export MYSQL_ROOT_PASSWORD=supersecret 

export MYSQL_ADDR=devmysql:3306