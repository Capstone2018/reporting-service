#!/usr/bin/env bash
# run the server in dev environment

source ./devenv.sh

netname=$NETWORK_NAME
if [ -z "$(docker network ls --filter name=$netname --quiet)" ]; then
    docker network create $netname
fi

if [ "$(docker ps -aq --filter name=devpsql)" ]; then
    docker rm -f devpsql
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

# docker run -d \
# --name devredis \
# --network $netname \
# redis

# docker run -d \
# --name devmysql \
# --network $netname \
# -p 3306:3306 \
# -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
# -e MYSQL_DATABASE=$MYSQL_DATABASE \
# $REPORTING_MYSQL_IMAGE --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

docker run -d \
--name devpsql \
--network $netname \
-p 5432:5432 \
-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
-e POSTGRES_USER=$POSTGRES_USER \
-e POSTGRES_DB=$POSTGRES_DB \
$REPORTING_POSTGRES_IMAGE

docker run -d \
--name devreporting \
--network $netname \
-p 4040:4040 \
-v $(pwd)/../server/tls:/tls:ro \
-e ADDR=:4040 \
-e TLSKEY=/tls/privkey.pem \
-e TLSCERT=/tls/fullchain.pem \
-e POSTGRES_HOST=$POSTGRES_HOST \
-e POSTGRES_PORT=$POSTGRES_PORT \
-e POSTGRES_USER=$POSTGRES_USER \
-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
-e POSTGRES_DB=$POSTGRES_DB \
$REPORTING_SERVICE_IMAGE .

#-e SESSIONKEY=$(uuidgen) \
#-e REDISADDR=devredis:6379 \



