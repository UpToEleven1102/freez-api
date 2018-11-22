all:
    go get golang.org/x/crypto/bcrypt
    go get github.com/dgrijalva/jwt-go
	go build -o freeze-app
	./freeze-app