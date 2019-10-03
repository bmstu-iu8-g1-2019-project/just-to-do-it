#!/bin/sh

wget -qO- https://toolbelt.heroku.com/install-ubuntu.sh | sh
heroku plugins:install heroku-container-registry

docker login -e _ -u _ --password=$HEROKU_API_KEY registry.heroku.com

docker build --tag=registry.heroku.com/${HEROKU_APP_NAME}/web .
docker push registry.heroku.com/${HEROKU_APP_NAME}/web
heroku container:release web --app ${HEROKU_APP_NAME}