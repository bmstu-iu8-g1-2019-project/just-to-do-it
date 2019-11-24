FROM golang:1.9.2

RUN apt-get -y update
RUN apt-get install -y tree wget curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

USER root

WORKDIR $GOPATH/src/github.com/bmstu-iu8-g1-2019-project/just-to-do-it
ADD ./ $GOPATH/src/github.com/bmstu-iu8-g1-2019-project/just-to-do-it

RUN chmod +x ./scripts/*
RUN ./scripts/build.sh

CMD ["./server.app"]
