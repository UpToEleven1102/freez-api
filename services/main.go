package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/db"
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
}
