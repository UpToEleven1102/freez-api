PROJECT_NAME := "freez-app-rest"
PKG := "git.nextgencode.io/huyen.vu/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all: build

build:
	##go get github.com/tbalthazar/onesignal-go
	##go get github.com/satori/go.uuid
	##go get github.com/jmoiron/sqlx
	##go get github.com/go-sql-driver/mysql
	##go get golang.org/x/crypto/bcrypt
	##go get github.com/dgrijalva/jwt-go
	##go get github.com/joho/godotenv
	##go get github.com/aws/aws-sdk-go
	##go get github.com/go-redis/redis
	##go get github.com/stripe/stripe-go
	##go get golang.org/x/net/websocket
	##go get -u github.com/huandu/facebook
	##go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -v -d
	CGO_ENABLED=0 GOOS=linux go build -o freez-app-rest

build-docker: build
	sudo docker build -t freez-app-rest .

dep: ##install dependencies
	@go get -v -d
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint: ##lint the files
	@golangci-lint run

test: ##test code
	echo "no test file was made"

push-code: ## push code to remote server
	scp -i /home/huyen/.ssh/Freeze.pem -r /home/huyen/gospace/src/git.nextgencode.io/huyen.vu/freez-app-rest/ ubuntu@35.162.158.187:/home/ubuntu/go/src/git.nextgencode.io/huyen.vu

generate-docs: ## generate swagger docs
	swagger -apiPackage=git.nextgencode.io/huyen.vu/freez-app-rest -format=swagger -output=./docs

dev-up: build-docker
	sudo docker-compose kill
	sudo docker-compose up

dev:
	go run main.go
