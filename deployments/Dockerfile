FROM golang:1.16.4-alpine3.13 as builder

WORKDIR /usr/local/go/src/just-to-do-it/cmd

COPY . /usr/local/go/src/just-to-do-it

RUN go build -o main .

FROM alpine:3.13

COPY --from=builder /usr/local/go/src/just-to-do-it/cmd/main /usr/local/go/src/just-to-do-it/cmd/main

ENTRYPOINT ["/usr/local/go/src/just-to-do-it/cmd/main"]
