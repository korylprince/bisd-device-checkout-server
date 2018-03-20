#!/bin/bash

version=$1

cred=$(cat ~/.git-credentials)

docker build --no-cache --build-arg "CREDENTIALS=$cred" --build-arg "VERSION=$version" --tag "korylprince/bisd-device-checkout-server:$version" .

docker push "korylprince/bisd-device-checkout-server:$version"
