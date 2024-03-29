package db

import (
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DB *sqlx.DB
var err error

func init() {
	//config.SetEnv()
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func seed(DB *sqlx.DB) {
	DB.MustExec("DROP TABLE IF EXISTS m_order_product")
	DB.MustExec("DROP TABLE IF EXISTS product")
	DB.MustExec("DROP TABLE IF EXISTS location")
	DB.MustExec("DROP TABLE IF EXISTS favorite")
	DB.MustExec("DROP TABLE IF EXISTS merchant_m_option")
	DB.MustExec("DROP TABLE IF EXISTS m_option")
	DB.MustExec("DROP TABLE IF EXISTS user")
	DB.MustExec("DROP TABLE IF EXISTS request")
	DB.MustExec("DROP TABLE IF EXISTS  merchant_notification")
	DB.MustExec("DROP TABLE IF EXISTS user_notification")
	DB.MustExec("DROP TABLE IF EXISTS activity_type")
	DB.MustExec("DROP TABLE IF EXISTS m_order")
	DB.MustExec("DROP TABLE IF EXISTS merchant")
	DB.MustExec("DROP TABLE IF EXISTS merchant_category")

	DB.MustExec(schemaMerchantCategory)
	DB.MustExec(schemaActivityType)

	tx := DB.MustBegin()
	tx.MustExec(`INSERT INTO merchant_category (category) VALUE (?)`, config.MerchantCategoryBreakfast)
	tx.MustExec(`INSERT INTO merchant_category (category) VALUE (?)`, config.MerchantCategoryLunch)
	tx.MustExec(`INSERT INTO merchant_category (category) VALUE (?)`, config.MerchantCategoryDinner)
	tx.MustExec(`INSERT INTO merchant_category (category) VALUE (?)`, config.MerchantCategorySweets)

	tx.MustExec(`INSERT INTO activity_type (type) VALUE (?)`, config.NotifTypeFlagRequest)
	tx.MustExec(`INSERT INTO activity_type (type) VALUE (?)`, config.NotifTypePaymentMade)
	tx.MustExec(`INSERT INTO activity_type (type) VALUE (?)`, config.NotifTypeRefundMade)
	tx.MustExec(`INSERT INTO activity_type (type) VALUE (?)`, config.NotifTypeRefundMade)
	//
	//tx.MustExec("INSERT INTO request (user_id, merchant_id, location) VALUES (123, '3412',ST_GeomFromText('POINT(1 1)'))")
	//uid, _ := uuid.NewV4()
	//tx.MustExec("INSERT INTO merchant (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","hotdog@truck.com", "Hot Dog Truck", "hot dog password")
	//uid, _ = uuid.NewV4()
	//tx.MustExec("INSERT INTO user (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)", uid.String(), "8013215431","a@truck.com", "H", "hot dog password")
	//
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	DB.MustExec(schemaMerchant)
	DB.MustExec(schemaProduct)
	DB.MustExec(schemaMerchantMOption)
	DB.MustExec(triggerInsertMerchantMOption)
	DB.MustExec(schemaUser)
	DB.MustExec(schemaMOption)
	DB.MustExec(triggerInsertUserMOption)
	DB.MustExec(schemaRequest)
	DB.MustExec(schemaLocation)
	DB.MustExec(schemaFavorites)
	DB.MustExec(schemaMerchantNotification)
	DB.MustExec(schemaUserNotification)
	DB.MustExec(schemaOrder)
	DB.MustExec(schemaOrderProduct)
}

func Config() (*sqlx.DB, error) {
	dbUri := getMysqlUri()

	if DB == nil {
		DB, err = sqlx.Connect("mysql", dbUri)

		if err != nil {
			log.Println("Failed to connect to DB. Sleep for awhile")
			return nil, err
		}
	}

	log.Printf("Reset db: %s", os.Getenv("RESET_DB"))

	if DB != nil {
		if os.Getenv("RESET_DB") == "true" {
			seed(DB)
			log.Println("Successfully reset DB")
		}
	}

	return DB, err
}

func getMysqlUri() (uri string) {
	uri = os.Getenv("MYSQL_URI")
	if len(uri) == 0 {
		uri = `h@tcp(127.0.0.1:3306)/freeze_app`
	}
	return uri
}

func Close() {
	//DB.Close()
}