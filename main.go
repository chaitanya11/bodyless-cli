package main

import (
	"fmt"
	"os"

	"github.com/chaitanya11/bodyless/git"
	"github.com/chaitanya11/bodyless/s3"
	flag "github.com/ogier/pflag"
)

var (
	projectName   string
	path          string
	codeBucket    string
	contentBucket string
	profile       string
)

func main() {
	flag.Parse()
	// if user does not supply flags, print usage
	if flag.NFlag() == 0 {
		printUsage()
	}

	// clone cms repo in to path
	git.PullGitRepo("https://github.com/chaitanya11/BodylessCMS", path)

	// create Code Bucket
	s3.CreateBucket(codeBucket)
	s3.SetWebSiteConfig(codeBucket, "index", "index.html")

	// create content Bucket
	s3.CreateBucket(contentBucket)
}

func init() {
	flag.StringVarP(&projectName, "ProjectName", "N", "", "Name of the project.")
	flag.StringVarP(&path, "Path", "P", ".", "Project Location.")
	flag.StringVarP(&codeBucket, "CodeBucketName", "w", "", "Name of the bucket where website code is deployed.")
	flag.StringVarP(&contentBucket, "ContentBucketName", "c", "", "Name of the bucket where content is stored.")
	flag.StringVarP(&profile, "profile", "p", "default", "Name of the aws profile configured.")
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
