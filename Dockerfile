FROM golang:latest
RUN mkdir -p /go/src/github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/
ADD . /go/src/github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/main
WORKDIR /go/src/github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/main/
RUN go get -v