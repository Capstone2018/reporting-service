#!/usr/bin/env bash
set -e
source secrets/docker.sh
domain=snopes.io

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

#ensureDomainRecord ensures that the correct domain
#record exists and that it is pointing at the right
#droplet IP address
#parameters:
# 1. domain name
# 2. record type (e.g., 'A')
# 3. record name (e.g., '@' or 'api')
# 4. droplet IP address
ensureDomainRecord () {
    #get ID of existing DOM record, if any
    echo >&2 "checking for existing domain $2 record for name $3 on domain $1..."
    recID=$(doctl compute domain records list $1 | awk '$2 == "'$2'" && $3 == "'$3'" {print $1}')
    if [ "$recID" ]; then
        currentIP=$(doctl compute domain records list $1 | awk '$2 == "'$2'" && $3 == "'$3'" {print $4}')
        if [ "$currentIP" == "$4" ]; then
            echo >&2 "$2 record for name $3 on domain $1 is up to date"
        else
            echo >&2 "updating $2 record for name $3 on domain $1 to point to IP $4..."
            doctl compute domain records update $1 \
            --record-id $recID \
            --record-type $2 \
            --record-name $3 \
            --record-data $4
        fi
    else
        echo >&2 "creating $2 record for name $3 on domain $1 to point to IP $4..."
        doctl compute domain records create $domain \
        --record-type $2 \
        --record-name $3 \
        --record-data $4
    fi
}

deployAPI () {
    docker push aethan/reporting-service
    docker push aethan/postgresreports

    dropletName=reporting-service
    dropletIP=$(ensureDroplet $dropletName)
    ensureDomainRecord $domain 'A' 'api' $dropletIP
    echo >&2 "ensuring that $dropletName server at $dropletIP is provisioned..."
    ssh -oStrictHostKeyChecking=no root@$dropletIP 'bash -s' < provision.sh "api.$domain"
    echo >&2 "updating docker containers on $dropletName server at $dropletIP..."
    ssh -oStrictHostKeyChecking=no root@$dropletIP 'bash -s' < update-api.sh $DOCKER_USERNAME $DOCKER_PASSWORD
}

if [[ "$1" == "api" ]]; then
    deployAPI
else
    echo "usage: ./deploy.sh api"
fi