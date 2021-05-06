FROM golang:1.16.3-buster

RUN go get -u github.com/civo/civogo
RUN go get -u github.com/sirupsen/logrus

RUN mkdir -p /go/src/github.com/valerauko/panther
WORKDIR /go/src/github.com/valerauko/panther
