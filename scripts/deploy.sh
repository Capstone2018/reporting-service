#!/usr/bin/env bash
set -e

#ensureDroplet will ensure that the droplet exists
#parameters:
# 1. droplet name
ensureDroplet () {
    #get IP of existing droplet, if any
    echo >&2 "checking for existing droplet named $1..."
    dropletIP=$(doctl compute droplet list $1 --no-header --format PublicIPv4)

    #if no existing droplet
    if [ "$dropletIP" ]; then
        echo >&2 "droplet $1 already exists"
    else
        echo >&2 "creating new droplet named $1..."
        dropletIP=$(doctl compute droplet create $1 \
        --image docker-16-04 \
        --region sfo2 \
        --size 1gb \
        --ssh-keys 65:90:af:93:b1:5f:b2:73:64:e3:05:14:bf:02:cb:a1 \
        --format PublicIPv4 \
        --no-header \
        --wait)

        echo >&2 "created new droplet named $1 with IP address $dropletIP"

        #although the --wait flag is supposed to wait
        #until the droplet is ready, it's not always
        #ready for an ssh connection yet, so sleep for a bit
        echo >&2 "letting new droplet get ready for ssh connection..."
        sleep 30s
    fi
    echo $dropletIP
}

deployAPI () {
    docker push aethan/reporting-service
    docker push aethan/mysqlreports

    dropletName=snopes-reporting-service
    dropletIP=$(ensureDroplet $dropletName)
    echo >&2 "ensuring that $dropletName server at $dropletIP is provisioned..."
    ssh -oStrictHostKeyChecking=no root@$dropletIP 'bash -s' < provision.sh
    echo >&2 "updating docker containers on $dropletName server at $dropletIP..."
    ssh -oStrictHostKeyChecking=no root@$dropletIP 'bash -s' < update-api.sh $DOCKER_USERNAME $DOCKER_PASSWORD
}

if [[ "$1" == "api" ]]; then
    deployAPI
else
    echo "usage: ./deploy.sh api"
fi