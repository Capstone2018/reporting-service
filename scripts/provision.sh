#!/usr/bin/env bash
set -e

#ensure firewall is on
fwStatus=$(ufw status | awk '$1=="Status:" {print $2}')
if [ "$fwStatus" != "active" ]; then
    echo >&2 "enabling firewall..."
    ufw enable
fi

#ensure ports are open on the firewall
# portAction=$(ufw status | awk '$1 == "80" && $2 != "(v6)" {print $2}')
# if [ "$portAction" != "ALLOW" ]; then
#     echo >&2 "opening port 80..."
#     ufw allow 80
# fi
portAction=$(ufw status | awk '$1=="443" && $2 != "(v6)" {print $2}')
if [ "$portAction" != "ALLOW" ]; then
    echo >&2 "opening port 443..."
    ufw allow 443
fi

#ensure lets encrypt cert/key exist
if [ -d "/etc/letsencrypt/live/$1" ]; then
    echo >&2 "Let's Encrypt cert and key already exist"
else
    echo >&2 "need to get TLS certificate and key...waiting a bit for domain records to update..."
    sleep 30s

    echo >&2 "ensuring letsencrypt is installed..."
    apt-get update && apt install -y letsencrypt

    echo >&2 "getting TLS certificate and key from Let's Encrypt..."

    letsencrypt certonly \
    --standalone \
    -n \
    --agree-tos \
    --email operations@teamsnopes.com \
    -d $1
fi