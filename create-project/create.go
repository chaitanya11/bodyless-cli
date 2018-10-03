package create_project

import (
	"bodyless-cli/git"
	"os"
	"bodyless-cli/configuration"
	"bodyless-cli/constants"
	"bodyless-cli/aws"
	"log"
)

func CreateProj(projectName string,
	path string,
	codeBucket string,
	profile string,
	region string) {
	// setting profile
	os.Setenv(constants.PROFILE_ENV_KEY, profile)

	// clone cms repo in to path
	log.Println("creating project...")
	path += projectName
	//cleaning given path
	os.RemoveAll(path)
	log.Println("cloneing to repo in to "+path)
	git.PullGitRepo(constants.REPO, path)
	log.Println("repo cloned.")

	// create-project Code Bucket
	aws.CreateBucket(codeBucket, &region)
	aws.SetWebSiteConfig(codeBucket, constants.S3_INDEX_PAGE, constants.S3_INDEX_PAGE, &region)

	// create cognito resources.
	cognitoConfig := aws.CreateCognitoResources(constants.COGNITO_POOL_NAME, &path, &region)

	// create configuration files.
	log.Println("writing repo configuraton...")
	path += "/"+constants.CONFIG_DIR
	os.Mkdir(path, constants.CONFIG_FILE_PERMISSIONS);
	filePath := path+"/"+constants.CONFIG_FILE_NAME
	configuration.WriteConfig(codeBucket, region, profile, filePath,cognitoConfig)
	log.Println("writing repo configuration is completed.")
}
