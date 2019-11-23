FROM golang:1.9.2

EXPOSE 8080

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

# RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh
# RUN chmod +x wait-for-it.sh


RUN chmod +x ./scripts/*
RUN ./scripts/start.sh
RUN go env
RUN tree -L 4 ./vendor

CMD [""]
