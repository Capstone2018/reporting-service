if [ -z "$1" ]; then
    echo "usage: ./self-signed.sh ./path/to/directory"
    exit 1
fi
mkdir -p $1
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj "/CN=localhost" -keyout $1/privkey.pem -out $1/fullchain.pem
