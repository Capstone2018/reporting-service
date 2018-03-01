#!/usr/bin/env bash
set -e

echo "building reporting service"
cd ../server
GOOS=linux go build

echo "building reporting service docker image"
docker build -t aethan/reporting-service .

echo "cleaning up reporting service.."
go clean
cd -

echo "building mysql database"
cd ./sql
docker build -t aethan/postgresreports .
cd -

echo "pruning.."
docker system prune -f