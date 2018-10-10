package remove_project

import (
	"bodyless-cli/constants"
	"bodyless-cli/configuration"
	"bodyless-cli/aws"
	"log"
	"bodyless-cli/utils"
)

func RemoveResources(path string) {
	configPath := path + constants.CONFIG_DIR
	var ngBucketName string
	// check for config.
	if configuration.CheckConfigDir(&configPath) {
		bodylessProjectConfig := configuration.ReadConfig(&path)
		ngBucketName = "ng-"+bodylessProjectConfig.BucketName
		log.Println("Removing all aws resources.")
		// remove s3 resources.
		aws.DeleteBucket(&bodylessProjectConfig.BucketName, &bodylessProjectConfig.Region)
		aws.DeleteBucket(&ngBucketName, &bodylessProjectConfig.Region)
		// remove cognito resources.
		aws.DeleteCognitoPool(&bodylessProjectConfig.CognitoConfig.UserPoolId, &bodylessProjectConfig.Region)
		aws.DeleteIdentityPool(&bodylessProjectConfig.CognitoConfig.IdentityPoolId, &bodylessProjectConfig.Region)
		// remove iam resources.
		aws.DeleteIamRole(constants.AUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME,
			constants.AUTHENTICATED_USER_ROLE_POLICY_NAME,
			&bodylessProjectConfig.Region)
		aws.DeleteIamRole(constants.UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME,
			constants.UNAUTHENTICATED_USER_ROLE_POLICY_NAME,
			&bodylessProjectConfig.Region)
		// delete project.
		utils.RemoveDirectory(path)
	}
}
