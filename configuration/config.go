package configuration

import (
	"encoding/json"
	"os"
	"io/ioutil"
	"bodyless-cli/constants"
	"bodyless-cli/utils"
)




func WriteConfig(
	bucketName string,
	awsRegion string,
	awsProfile string,
	fileName string,
	cognitoConfig constants.PROJECT_CONF_TEMPLATE_VARS) {
	bodylessProjectConfig := constants.BodylessProjectConfig{BucketName:bucketName,
		Region:awsRegion,
		Profile:awsProfile,
		CognitoConfig:cognitoConfig}

	b, jsonErr := json.Marshal(bodylessProjectConfig)

	utils.CheckNExitError(jsonErr)

	// write data to file
	fileWriteErr := ioutil.WriteFile(fileName, b, constants.CONFIG_FILE_PERMISSIONS);

	utils.CheckNExitError(fileWriteErr)
}

func ReadConfig() constants.BodylessProjectConfig {
	file, _ := os.Open(constants.CONFIG_DIR + "/" + constants.CONFIG_FILE_NAME)
	defer file.Close()
	decoder := json.NewDecoder(file)
	bodylessProjectConfig := constants.BodylessProjectConfig{}
	err := decoder.Decode(&bodylessProjectConfig)
	utils.CheckNExitError(err)
	return bodylessProjectConfig
}