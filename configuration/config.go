package configuration

import (
	"encoding/json"
	"os"
	"io/ioutil"
	"bodyless-cli/constants"
	"bodyless-cli/utils"
	"log"
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
	fileWriteErr := ioutil.WriteFile(fileName, b, constants.CONFIG_FILE_PERMISSIONS)

	utils.CheckNExitError(fileWriteErr)
}

func ReadConfig(path *string) constants.BodylessProjectConfig {
	file, _ := os.Open(*path + constants.CONFIG_DIR + "/" + constants.CONFIG_FILE_NAME)
	defer file.Close()
	decoder := json.NewDecoder(file)
	bodylessProjectConfig := constants.BodylessProjectConfig{}
	err := decoder.Decode(&bodylessProjectConfig)
	utils.CheckNExitError(err)
	log.Printf("Project configuration %+v", bodylessProjectConfig)
	return bodylessProjectConfig
}

func CheckConfigDir(path *string) bool {
	log.Printf("checking %s directory for configuration ....", *path)
	var result bool
	_, err := os.Stat(*path)
	if err == nil {
			log.Printf("configuration found in %s directory.", *path)
			result = true
		}
	if os.IsNotExist(err) {
		log.Printf("configuration not found in %s directory.", *path)
		result = false
		}
	return result
}