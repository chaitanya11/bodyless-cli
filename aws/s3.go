package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"bodyless-cli/utils"
	"bodyless-cli/constants"
	"bodyless-cli/build-project"
	"path/filepath"
	"mime"
	"path"
	"strings"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// aws s3
func CreateBucket(bucketName string,
	region *string) {
	log.Printf("Creating s3 bucket with name %s ...", bucketName)
	svc := s3.New(session.New(&aws.Config{
		Region: region,
	}))
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

	log.Printf("Created s3 bucket with name %s", bucketName)

	// make bucket public-read
	setPublicBucketPolicy(&bucketName, svc)

	// enable versioning on s3 bucket
	// TODO fix multi version delete objects.
	//enableBucketVersioning(&bucketName, svc)

	// create cors config
	createCorsConfig(&bucketName, svc)

	log.Println(result)
}

func CreateDeploymentFiles(s3BucketName string, region *string) {
	bucketName := s3BucketName
	// add one bodylesscms logo, index.html to s3 bucket
	imgFileName := constants.S3_IMG_FILE
	imgObjectName := constants.S3_IMG_PATH + imgFileName
	uploadFile(&bucketName, region, &imgObjectName, imgFileName)

	indexFileName := constants.S3_INDEX_PAGE
	uploadFile(&bucketName, region, &indexFileName, indexFileName)
	styleFileName := constants.S3_STYLE_PAGE
	uploadFile(&bucketName, region, &styleFileName, styleFileName)

}

func CreateNgCodeFiles(s3BucketName string, region *string, path *string) {
	log.Printf("Uploading angular files to %s bucket ...", s3BucketName)
	bucketName := s3BucketName
	// build ng project.
	build_project.BuildProj(*path)
	// recursively upload files ng build files to s3.
	uploadDirToS3(bucketName, region, *path+"/"+constants.BUILD_FILES_PATH)
	log.Printf("Uploaded angular files to %s bucket ...", s3BucketName)
}

func enableBucketVersioning(bucketName *string, svc *s3.S3) {
	log.Println("Enabling versioning s3 bucket ...")
	_, err := svc.PutBucketVersioning(&s3.PutBucketVersioningInput{
		Bucket: bucketName,
		VersioningConfiguration: &s3.VersioningConfiguration{
			MFADelete: aws.String("Disabled"),
			Status:    aws.String("Enabled"),
		},
	})
	utils.CheckNExitError(err)
	log.Println("Enabled versioning on s3 bucket")
}

func setPublicBucketPolicy(bucketName *string, svc *s3.S3) {
	log.Printf("Making %s bucket as public ...", *bucketName)
	policyTemplate := constants.PUBLIC_BUCKET_POLICY_TEMPLATE
	policyDoc := utils.GetStringFromTemplate(policyTemplate, constants.BodylessProjectConfig{
		BucketName: *bucketName,
	})
	input := &s3.PutBucketPolicyInput{
		Bucket: bucketName,
		Policy: aws.String(policyDoc),
	}
	_, err := svc.PutBucketPolicy(input)
	utils.CheckNExitError(err)
	log.Printf("Made %s bucket as public", *bucketName)
}

func uploadFile(bucketName *string, region *string, objectName *string, fileName string) {
	log.Printf("Uploading file %s ...", fileName)

	file, FileOpenErr := os.Open(fileName)
	utils.CheckNExitError(FileOpenErr)
	defer file.Close()
	contentType := mime.TypeByExtension(path.Ext(*objectName))
	svc := s3.New(session.New(&aws.Config{
		Region: region,
	}))
	input := &s3.PutObjectInput{
		Body:   file,
		Bucket: bucketName,
		Key:    objectName,
		ContentType: &contentType,
	}

	_, err := svc.PutObject(input)
	utils.CheckNExitError(err)
	log.Printf("Uploaded file %s", fileName)
}

func uploadDirToS3(bucketName string, region *string, dirPath string) {
	log.Printf("Uploading files from %s directory ...", dirPath)
	fileList := []string{}
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		log.Println("PATH ==> " + path)
		if utils.IsDirectory(path) {
			// Do nothing
			return nil
		} else {
			fileList = append(fileList, path)
			return nil
		}
	})

	for _, file := range fileList {
		fileName := filepath.Base(file)
		if strings.Contains(file, "assets") {
			fileName = "assets/" + fileName
		}
		uploadFile(&bucketName, region, &fileName, file)
	}
	log.Printf("Uploaded files from %s directory ...", dirPath)
}

func createCorsConfig(bucketName *string, svc *s3.S3) {
	log.Printf("Setting cors for %s bucket ...", *bucketName)
	_, err := svc.PutBucketCors(&s3.PutBucketCorsInput{
		Bucket: bucketName,
		CORSConfiguration:  &s3.CORSConfiguration{
			CORSRules: []*s3.CORSRule{
				{
					AllowedHeaders: []*string{
						aws.String("*"),
					},
					AllowedMethods: []*string{
						aws.String("PUT"),
						aws.String("POST"),
						aws.String("GET"),
					},
					AllowedOrigins: []*string{
						aws.String("*"),
					},
				},
			},
		},
	})
	utils.CheckNExitError(err)
	log.Printf("Setting cors for %s bucket is completed.", *bucketName)
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

func DeleteBucket(bucketName *string, region *string) {
	EmptyBucket(bucketName, region)
	log.Printf("Deleting %s bucket ...", *bucketName)
	svc := s3.New(session.New(&aws.Config{
		Region: region,
	}))
	input := &s3.DeleteBucketInput{
		Bucket: bucketName,
	}
	_, err := svc.DeleteBucket(input)
	utils.CheckNExitError(err)
	log.Printf("Deleted %s bucket.", *bucketName)
}

func EmptyBucket(bucketName *string, region *string) {
	log.Printf("Deleting all files in %s bucket ...", *bucketName)
	sess, _ := session.NewSession(&aws.Config{
		Region: region},
	)

	// Create S3 service client
	svc := s3.New(sess)

	// Setup BatchDeleteIterator to iterate through a list of objects.
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: bucketName,
	})

	// Traverse iterator deleting each object
	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		exitErrorf("Unable to delete objects from bucket %q, %v", bucketName, err)
	}

	log.Printf("Deleted object(s) from bucket: %s", *bucketName)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
