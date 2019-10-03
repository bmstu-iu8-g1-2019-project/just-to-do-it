#!/bin/sh

if [ "$TRAVIS_BRANCH" = "master" ]; then
    TAG="latest"
else
    TAG="$TRAVIS_BRANCH"
fi

echo "$DOCKER_PASS" | docker login -u $DOCKER_USER --password-stdin

export DOCKER_IMAGE=cicd:$TAG

docker build --tag=$DOCKER_USER/$DOCKER_IMAGE .
docker push $DOCKER_USER/$DOCKER_IMAGE
