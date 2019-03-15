package services

import (
	"git.nextgencode.io/huyen.vu/freez-app-rest/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go"
	"github.com/tbalthazar/onesignal-go"
	"log"
	"os"
	"time"
)

var (
	DB              *sqlx.DB
	oneSignalClient *onesignal.Client
	oneSignalAppID  string

	s3Client     *s3.S3
	s3Uploader   *s3manager.Uploader
	s3BucketName string

	RedisClient *redis.Client
)

const (
	minDistance = 3000
	minNotifyDistance = 1000
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

func redisConfig() {
	redisServer := os.Getenv("REDIS_ADDRESS")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisServer,
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Printf("%s\n", err)
	}
}

//func loadNotifTypeToRedis() {
//	type ActivityType struct {
//		ID   int    `json:"id"`
//		Type string `json:"type"`
//	}
//
//	r, err := DB.Query(`SELECT id, type FROM activity_type`)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	//var activityTypes []ActivityType
//	var activityType ActivityType
//	defer r.Close()
//
//	for r.Next() {
//		_ = r.Scan(&activityType.ID, &activityType.ID)
//		//activityTypes = append(activityTypes, activityType)
//		RedisClient.Set(fmt.Sprintf("%d", activityType.ID), activityType, 3600*time.Minute)
//	}
//}

func connectDB() {
	var err error
	DB, err = db.Config()

	if err != nil {
		time.AfterFunc(5 * time.Second, connectDB)
	}
}

func init() {
	connectDB()
	oneSignalConfig()

	awsConfig()

	redisConfig()

	stripe.Key = os.Getenv("STRIPE_KEY")

	//listObjects()

	//CreateEmailNotification("quyhuyen.vu@gmail.com", "", "hello")

	//UploadBlankProfilePicture()
	//GetBucketLocation()
}
