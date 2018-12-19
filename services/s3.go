package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
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
	})

	url, err = req.Presign(expirationTime)

	return url, err
}