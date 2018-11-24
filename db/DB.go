package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"github.com/satori/go.uuid"
)

var DB *sqlx.DB
var err error

func Config() (*sqlx.DB, error) {
	dbUri := getMysqlUri()

	if DB == nil {
		DB, err = sqlx.Connect("mysql", dbUri)
		if err != nil {
			panic(err)
		}
	}


	DB.MustExec("DROP TABLE IF EXISTS merchant")
	DB.MustExec("DROP TABLE IF EXISTS user")
	DB.MustExec("DROP TABLE IF EXISTS request")
	DB.MustExec(schemaMerchant)
	DB.MustExec(schemaUser)
	DB.MustExec(schemaRequest)

	tx := DB.MustBegin()

	tx.MustExec("INSERT INTO request (user_id, location) VALUES (123, ST_GeomFromText('POINT(1 1)'))")

	uid, _ := uuid.NewV4()
	tx.MustExec("INSERT INTO merchant (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "3023324324","icecream@truck.com","Ice Cream Truck", "Password")
	uid, _ = uuid.NewV4()
	tx.MustExec("INSERT INTO merchant (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","hotdog@truck.com", "Hot Dog Truck", "hot dog password")
	uid, _ = uuid.NewV4()
	tx.MustExec("INSERT INTO user (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","h@truck.com", "AJ", "hot dog password")
	uid, _ = uuid.NewV4()
	tx.MustExec("INSERT INTO user (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","a@truck.com", "H", "hot dog password")

	tx.Commit()

	return DB, err
}

func getMysqlUri() (uri string) {
	uri = os.Getenv("MYSQL_URI")
	if len(uri) == 0 {
		uri = `root@tcp(127.0.0.1:3306)/freeze_app`
	}
	return uri
}

func Close() {
	DB.Close()
}