package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/db"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/jmoiron/sqlx"
	"github.com/tbalthazar/onesignal-go"
	"os"
)

var (
	DB *sqlx.DB
	oneSignalClient *onesignal.Client
	appID, appKey, userKey string
)

func init() {
	DB, _ = db.Config()

	appID = os.Getenv("ONE_SIGNAL_APP_ID")
	appKey = os.Getenv("ONE_SIGNAL_APP_KEY")
	userKey = os.Getenv("ONE_SIGNAL_USER_KEY")

	oneSignalClient = onesignal.NewClient(nil)
	oneSignalClient.AppKey = appKey
	oneSignalClient.UserKey = userKey

	var data models.RequestData

	data.UserId = "User Id"
	data.Data = "s3cr3t"

	CreateNotificationByUserId("0033dfe3-29e8-4be4-b2ad-e2814fdefa2c", "title", "Hello", data)
}
