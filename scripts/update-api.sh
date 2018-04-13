#!/usr/bin/env bash
echo "$1 $2"
docker login -u $1 -p $2
docker pull aethan/reporting-service
docker pull aethan/mysqlreports

netname=reporting

if [ -z "$(docker network ls --filter name=$netname --quiet)" ]; then
    docker network create $netname
fi

if [ "$(docker ps -aq --filter name=postgres)" ]; then
    docker rm -f postgres
fi

if [ "$(docker ps -aq --filter name=reporting)" ]; then
    docker rm -f reporting
fi

# if [ "$(docker ps -aq --filter name=redis)" ]; then
#     docker rm -f redis
# fi

# mkdir -p /tls
# openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj "/CN=localhost" -keyout /tls/privkey.pem -out /tls/fullchain.pem


# docker run -d \
# --name redis \
# --network $netname \
# redis


# docker container image name for reporting service
REPORTING_SERVICE_IMAGE=aethan/reporting-service

# docker container image name for our customized postgres image
REPORTING_POSTGRES_IMAGE=aethan/postgresreports

# database user which will be used to create schema
POSTGRES_USER=admin

# database name in which our schema will be created
POSTGRES_DB=reporting

# random MySQL root password 
# use $(openssl rand -base64 18) in prod
POSTGRES_PASSWORD=supersecret 

POSTGRES_HOST=postgres
POSTGRES_PORT=5432

# create a postgres database
docker run -d \
--name postgres \
--network $netname \
-p 5432:5432 \
-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
-e POSTGRES_USER=$POSTGRES_USER \
-e POSTGRES_DB=$POSTGRES_DB \
$REPORTING_POSTGRES_IMAGE

docker run -d \
--name reporting \
--network $netname \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSKEY=/etc/letsencrypt/live/api.snopes.io/privkey.pem \
-e TLSCERT=/etc/letsencrypt/live/api.snopes.io/fullchain.pem \
-e POSTGRES_HOST=$POSTGRES_HOST \
-e POSTGRES_PORT=$POSTGRES_PORT \
-e POSTGRES_USER=$POSTGRES_USER \
-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
-e POSTGRES_DB=$POSTGRES_DB \
$REPORTING_SERVICE_IMAGE .

docker system prune -f