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
	DB.MustExec(schema)

	tx := DB.MustBegin()
	uid, _ := uuid.NewV4()
	tx.MustExec("INSERT INTO merchant (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "3023324324","icecream@truck.com","Ice Cream Truck", "Password")
	uid, _ = uuid.NewV4()
	tx.MustExec("INSERT INTO merchant (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","hotdog@truck.com", "Hot Dog Truck", "hot dog password")
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