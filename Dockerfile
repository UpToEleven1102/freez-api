FROM golang:alpine

RUN apk update && apk upgrade && apk add --no-cache bash git

RUN go get github.com/tbalthazar/onesignal-go
RUN go get github.com/satori/go.uuid
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jmoiron/sqlx
RUN go get github.com/joho/godotenv
RUN go get github.com/aws/aws-sdk-go
RUN go get github.com/go-redis/redis

ENV SOURCES /go/src/git.nextgencode.io/huyen.vu/freeze-app-rest/

COPY . ${SOURCES}

RUN cd ${SOURCES} && CGO_ENABLED=0 go build

WORKDIR ${SOURCES}

CMD ${SOURCES}freeze-app-rest

EXPOSE 8080