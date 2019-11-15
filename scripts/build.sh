#!/bin/bash

main=src/main/main.go
app=server.app

if [[ ! -d 'vendor' ]]; then
    dep ensure -update
    dep ensure
fi

go build -o ${app} ${main}

md5sum ${app}
