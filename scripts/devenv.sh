#!/usr/bin/env bash

# you should `source` this script before building the
# Docker container image or trying to run it. Use this
# command to source this into your current shell:
#   source devenv.sh

# docker container image name for reporting service
export REPORTING_SERVICE_IMAGE=aethanol/reporting-service

# docker container image name for our customized MySQL image
export REPORTING_MYSQL_IMAGE=aethanol/mysqlreports

# database name in which our schema will be created
export MYSQL_DATABASE=reporting

# random MySQL root password
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)