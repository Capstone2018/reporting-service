#!/usr/bin/env bash
echo "$1 $2"
docker login -u $1 -p $2
docker pull aethan/reporting-service
docker pull aethan/mysqlreports

netname=reporting

if [ -z "$(docker network ls --filter name=$netname --quiet)" ]; then
    docker network create $netname
fi

if [ "$(docker ps -aq --filter name=mysql)" ]; then
    docker rm -f mysql
fi

if [ "$(docker ps -aq --filter name=reporting-service)" ]; then
    docker rm -f reporting-service
fi

if [ "$(docker ps -aq --filter name=redis)" ]; then
    docker rm -f redis
fi

mkdir -p /tls
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj "/CN=localhost" -keyout /tls/privkey.pem -out /tls/fullchain.pem


docker run -d \
--name redis \
--network $netname \
redis


# docker container image name for reporting service
REPORTING_SERVICE_IMAGE=aethan/reporting-service

# docker container image name for our customized MySQL image
REPORTING_MYSQL_IMAGE=aethan/mysqlreports

# database name in which our schema will be created
MYSQL_DATABASE=reporting

# random MySQL root password 
# use $(openssl rand -base64 18) in prod
MYSQL_ROOT_PASSWORD=supersecret 

MYSQL_ADDR=mysql:3306

docker run -d \
--name mysql \
--network $netname \
-p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
$REPORTING_MYSQL_IMAGE --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

docker run -d \
--name reporting-service \
--network $netname \
-p 443:443 \
--restart always \
-v /tls:/tls:ro \
-e ADDR=:443 \
-e TLSKEY=/tls/privkey.pem \
-e TLSCERT=/tls/fullchain.pem \
-e SESSIONKEY=$(uuidgen) \
-e REDISADDR=redis:6379 \
-e MYSQL_ADDR=mysql:3306 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=reporting \
$REPORTING_SERVICE_IMAGE .

docker system prune -f