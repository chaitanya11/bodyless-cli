package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"log"
)

// aws s3
func CreateBucket(bucketName string,
	region *string) {
	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))
	// TODO add versioning.
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: region,
		},
	}

	result, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				log.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				log.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		os.Exit(1)
		return
	}

	log.Println(result)
}

func SetWebSiteConfig(bucketName string, indexSuffix string, errorPage string, region *string) {
	sess, err := session.NewSession(&aws.Config{
		Region: region},
	)

	// Create S3 service client
	svc := s3.New(sess)
	params := s3.PutBucketWebsiteInput{
		Bucket: aws.String(bucketName),
		WebsiteConfiguration: &s3.WebsiteConfiguration{
			IndexDocument: &s3.IndexDocument{
				Suffix: aws.String(indexSuffix),
			},
		},
	}

	params.WebsiteConfiguration.ErrorDocument = &s3.ErrorDocument{
		Key: aws.String(errorPage),
	}
	_, err = svc.PutBucketWebsite(&params)
	if err != nil {
		exitErrorf("Unable to set bucket %q website configuration, %v",
			bucketName, err)
	}

	log.Printf("Successfully set bucket %q website configuration\n", bucketName)

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
