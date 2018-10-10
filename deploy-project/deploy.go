package deploy_project

import (
	"bodyless-cli/configuration"
	"bodyless-cli/aws"
	"bodyless-cli/utils"
	"errors"
	"bodyless-cli/constants"
)

func DeployProj(path string) {
	configPath := path + constants.CONFIG_DIR
	// check for config dir.
	var ngBucketName string
	if configuration.CheckConfigDir(&configPath) {
		// get configuration from config dir.
		bodylessProjectConfig := configuration.ReadConfig(&path)
		ngBucketName = "ng-"+bodylessProjectConfig.BucketName
		// empty ng bucket before deploying.
		aws.EmptyBucket(&ngBucketName, &bodylessProjectConfig.Region)
		// deploy project.
		aws.CreateNgCodeFiles(ngBucketName, &bodylessProjectConfig.Region, &path)
	} else {
		utils.CheckNExitError(errors.New("No configuration folder found in given path"))
	}

}