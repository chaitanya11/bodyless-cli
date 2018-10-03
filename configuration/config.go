package configuration

import (
	"encoding/json"
	"os"
	"fmt"
	"io/ioutil"
	"bodyless-cli/constants"
)

type BodylessProjectConfig struct {
	BucketName string
	Region string
	Profile string
}


func WriteConfig(
	bucketName string,
	awsRegion string,
	awsProfile string,
	fileName string) {
	bodylessProjectConfig := BodylessProjectConfig{BucketName:bucketName, Region:awsRegion, Profile:awsProfile}

	b, jsonErr := json.Marshal(bodylessProjectConfig)

	checkNExitOnError(jsonErr)

	// write data to file
	fileWriteErr := ioutil.WriteFile(fileName, b, constants.CONFIG_FILE_PERMISSIONS);

	checkNExitOnError(fileWriteErr)
}

func ReadConfig() BodylessProjectConfig {
	file, _ := os.Open(constants.CONFIG_DIR + "/" + constants.CONFIG_FILE_NAME)
	defer file.Close()
	decoder := json.NewDecoder(file)
	bodylessProjectConfig := BodylessProjectConfig{}
	err := decoder.Decode(&bodylessProjectConfig)
	checkNExitOnError(err)
	return bodylessProjectConfig
}


func checkNExitOnError(err error) {
	if err != nil {
		fmt.Println(err);
		os.Exit(1);
	}
}