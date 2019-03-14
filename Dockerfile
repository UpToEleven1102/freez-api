FROM golang:alpine

RUN apk update && apk upgrade && apk add --no-cache bash git make

ENV SOURCES /go/src/git.nextgencode.io/huyen.vu/freez-app-rest/

COPY . ${SOURCES}

RUN cd ${SOURCES} && make dep && CGO_ENABLED=0 go build

WORKDIR ${SOURCES}

CMD ${SOURCES}freez-app-rest

EXPOSE 8080