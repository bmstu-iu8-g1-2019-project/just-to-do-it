#!/usr/bin/env bash

export DB_host=todoapp_postgres
docker-compose build
docker-compose up -d
