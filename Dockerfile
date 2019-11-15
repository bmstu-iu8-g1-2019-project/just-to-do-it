FROM ubuntu:18.04

RUN apt-get -y update
RUN apt-get install -y wget tree git curl

RUN wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz && mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

USER root

EXPOSE 5000
EXPOSE 8080

ENV GO_DEEP "psql -U docker -h 127.0.0.1 -d todoapp"

WORKDIR $GOPATH/src/github.com/bmstu-iu8-g1-2019-project/just-to-do-it
ADD ./ $GOPATH/src/github.com/bmstu-iu8-g1-2019-project/just-to-do-it

RUN tree -L 4 ./

RUN chmod +x ./scripts/*
RUN ./scripts/build.sh


RUN md5sum server.app

ENTRYPOINT ["./scripts/start.sh"]
CMD [""]
