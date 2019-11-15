#!/bin/bash

main=src/main/main.go
app=server.app

if [[ ! -d 'vendor' ]]; then
    dep ensure -update
    dep ensure
fi

[[ ! -f ${app} ]] && go build -o ${app} ${main}

#service postgresql start

./${app} postgres://docker:docker@localhost:5555/todoapp