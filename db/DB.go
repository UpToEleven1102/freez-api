package db

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"os"
)

var DB *sqlx.DB
var err error

func init() {
	config.SetEnv()
}

func seed(DB *sqlx.DB) {
	DB.MustExec("DROP TABLE IF EXISTS merchant")
	DB.MustExec("DROP TABLE IF EXISTS user")
	DB.MustExec("DROP TABLE IF EXISTS request")
	DB.MustExec("DROP TABLE IF EXISTS location")
	DB.MustExec("DROP TABLE IF EXISTS favorite")
	DB.MustExec(schemaMerchant)
	DB.MustExec(schemaUser)
	DB.MustExec(schemaRequest)
	DB.MustExec(schemaLocation)
	DB.MustExec(schemaFavorites)

	tx := DB.MustBegin()

	tx.MustExec("INSERT INTO request (user_id, merchant_id, location) VALUES (123, '3412',ST_GeomFromText('POINT(1 1)'))")
	uid, _ := uuid.NewV4()
	tx.MustExec("INSERT INTO merchant (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","hotdog@truck.com", "Hot Dog Truck", "hot dog password")
	uid, _ = uuid.NewV4()
	tx.MustExec("INSERT INTO user (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","a@truck.com", "H", "hot dog password")

	tx.Commit()
}

func Config() (*sqlx.DB, error) {
	dbUri := getMysqlUri()

	if DB == nil {
		DB, err = sqlx.Connect("mysql", dbUri)
		if err != nil {
			panic(err)
		}
	}

	seed(DB)

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
	//DB.Close()
}