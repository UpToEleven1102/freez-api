all:
	go get github.com/tbalthazar/onesignal-go
	go get github.com/satori/go.uuid
	go get github.com/jmoiron/sqlx
	go get github.com/go-sql-driver/mysql
	go get golang.org/x/crypto/bcrypt
	go get github.com/dgrijalva/jwt-go
	go get github.com/joho/godotenv
	go get github.com/aws/aws-sdk-go
	go get github.com/go-redis/redis
	go get github.com/stripe/stripe-go
	go get golang.org/x/net/websocket
	go get -u github.com/huandu/facebook
	go build -o freeze-app
	./freeze-app
