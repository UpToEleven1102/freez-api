all:
	go get github.com/tbalthazar/onesignal-go
	go get github.com/paulsmith/gogeos/geos
	go get github.com/satori/go.uuid
	go get github.com/jmoiron/sqlx
	go get github.com/go-sql-driver/mysql
	go get golang.org/x/crypto/bcrypt
	go get github.com/dgrijalva/jwt-go
	go get github.com/joho/godotenv
	go get github.com/aws/aws-sdk-go
	go build -o freeze-app
	./freeze-app
