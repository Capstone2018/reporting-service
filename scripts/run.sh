#!/usr/bin/env bash
# run the server in dev environment

source ./devenv.sh

netname=$NETWORK_NAME
if [ -z "$(docker network ls --filter name=$netname --quiet)" ]; then
    docker network create $netname
fi

if [ "$(docker ps -aq --filter name=devmysql)" ]; then
    docker rm -f devmysql
fi

if [ "$(docker ps -aq --filter name=devreporting)" ]; then
    docker rm -f devreporting
fi

if [ "$(docker ps -aq --filter name=devredis)" ]; then
    docker rm -f devredis
fi

# ensure that the tls certs are generated
if [ ! -e $(pwd)/../server/tls ]; then
    ./self-signed.sh $(pwd)/../server/tls
fi

docker run -d \
--name devredis \
--network $netname \
redis

docker run -d \
--name devmysql \
--network $netname \
-p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
$REPORTING_MYSQL_IMAGE --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

docker run -d \
--name devreporting \
--network $netname \
-p 4040:4040 \
-v $(pwd)/../server/tls:/tls:ro \
-e ADDR=:4040 \
-e TLSKEY=/tls/privkey.pem \
-e TLSCERT=/tls/fullchain.pem \
-e SESSIONKEY=$(uuidgen) \
-e REDISADDR=devredis:6379 \
-e MYSQL_ADDR=$MYSQL_ADDR \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
$REPORTING_SERVICE_IMAGE .



