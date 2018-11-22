FROM golang:alpine

RUN apk update && apk upgrade && apk add --no-cache bash git

RUN go get github.com/satori/go.uuid
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jmoiron/sqlx

ENV SOURCES /work/src/github.com/UpToEleven1102/freezeapp-rest/

COPY . ${SOURCES}

RUN cd ${SOURCES} && CGO_ENABLED=0 go build

WORKDIR ${SOURCES}

CMD ${SOURCES}freezeapp-rest

EXPOSE 8080