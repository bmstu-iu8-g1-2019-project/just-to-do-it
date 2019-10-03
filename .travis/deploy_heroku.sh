#!/bin/sh
wget -qO- https://toolbelt.heroku.com/install-ubuntu.sh | sh
heroku plugins:install @heroku-cli/plugin-container-registry

echo "$HEROKU_API_KEY" | docker login -u _  --password-stdin registry.heroku.com

docker build --tag=registry.heroku.com/${HEROKU_APP_NAME}/web .
docker push registry.heroku.com/${HEROKU_APP_NAME}/web
heroku container:release web --app ${HEROKU_APP_NAME}e web --app ${HEROKU_APP_NAME}