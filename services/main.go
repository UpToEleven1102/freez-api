package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jmoiron/sqlx"
	"github.com/tbalthazar/onesignal-go"
	"os"
)

var (
	DB              *sqlx.DB
	oneSignalClient *onesignal.Client
	oneSignalAppID  string

	s3Client     *s3.S3
	s3Uploader *s3manager.Uploader
	s3BucketName string

)

func oneSignalConfig() {
	oneSignalAppID = os.Getenv("ONE_SIGNAL_APP_ID")
	appKey := os.Getenv("ONE_SIGNAL_APP_KEY")
	userKey := os.Getenv("ONE_SIGNAL_USER_KEY")

	oneSignalClient = onesignal.NewClient(nil)
	oneSignalClient.AppKey = appKey
	oneSignalClient.UserKey = userKey
}

func awsConfig() {
	s3BucketName = os.Getenv("AWS_BUCKET_NAME")
	awsRegion := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		panic(err)
	}

	s3Client = s3.New(sess)
	s3Uploader = s3manager.NewUploader(sess)
}


func init() {
	DB, _ = db.Config()
	oneSignalConfig()

	awsConfig()
	listObjects()
	//UploadBlankProfilePicture()
	//GetBucketLocation()
}
