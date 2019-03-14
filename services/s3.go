package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"time"
)

const (
	expirationTime = 10 * time.Minute
)

func GetBucketLocation() {
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(s3BucketName),
	}

	result, err := s3Client.GetBucketLocation(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func GeneratePreSignedUrl(fileName string) (url string , err error) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key: aws.String(fileName),
		ContentType: aws.String("image/jpeg"),
	})

	url, err = req.Presign(expirationTime)

	return url, err
}

func UploadBlankProfilePicture() {
	fileName := "blank-profile-picture.jpg"

	f, err := os.Open("/home/huyen/Pictures/blank-profile-picture-973460_640.jpg")
	if err != nil {
		panic(err)
	}

	result, err := s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3BucketName),
		Key: aws.String(fileName),
		ContentType: aws.String("image/jpeg"),
		Body: f,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("blank profile location: %s\n", result.Location)
}

//func listObjects() {
//	input := &s3.ListObjectsInput{
//		Bucket: aws.String(s3BucketName),
//		MaxKeys: aws.Int64(10),
//	}
//
//	result, err := s3Client.ListObjects(input)
//
//	if err != nil {
//		if aerr, ok := err.(awserr.Error); ok {
//			switch aerr.Code() {
//			case s3.ErrCodeNoSuchBucket:
//				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
//			default:
//				fmt.Println(aerr.Error())
//			}
//		} else {
//			// Print the error, cast err to awserr.Error to get the Code and
//			// Message from an error.
//			fmt.Println(err.Error())
//		}
//		return
//	}
//
//	//deleteObjects(result.Contents)
//
//	fmt.Println(result)
//}

//func deleteObjects(Contents []*s3.Object) {
//	for _, content := range Contents {
//		if *content.Key != "blank-profile-picture.jpg" {
//			deleteObject(content.Key)
//		}
//	}
//}

//func deleteObject(key *string) {
//	input := &s3.DeleteObjectInput{
//		Bucket: aws.String(s3BucketName),
//		Key: aws.String(*key),
//	}
//
//	_, err := s3Client.DeleteObject(input)
//	if err != nil {
//		fmt.Println(err)
//	}
//}


