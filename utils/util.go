package utils

import (
	"os"
	"bodyless-cli/constants"
	"text/template"
	"bytes"
	"log"
)

func CheckNExitError(err error) {
	if err != nil {
		panic(err)
		os.Exit(1)
	}
}

func WriteProjectConfig(userPoolId *string,
	clientId *string,
	identityPoolId *string,
	awsRegion *string,
	path *string,
	validRoleArn *string,
	invalidRoleArn *string) constants.PROJECT_CONF_TEMPLATE_VARS {
	templateVars := constants.PROJECT_CONF_TEMPLATE_VARS{UserPoolId: *userPoolId,
		ClientId: *clientId,
		IdentityPoolId: *identityPoolId,
		AwsRegion: *awsRegion,
		ValidRoleArn: *validRoleArn,
		InValidRoleArn: *invalidRoleArn,
	}
	template := template.New("projectConfig")
	template, parseErr := template.Parse(constants.PROJECT_CONF_TEMPLATE)
	CheckNExitError(parseErr)
	log.Println("Opening config file")
	configFile, configFileReadErr := os.OpenFile(*path + "/" +constants.PROJECT_CONFIG_PATH,
		os.O_RDWR, os.ModePerm)
	CheckNExitError(configFileReadErr)
	log.Println("Writing data to config file")
	exeErr := template.Execute(configFile, templateVars)
	CheckNExitError(exeErr)
	configFile.Close()
	log.Println("Config file is loaded with configuration.")
	return templateVars
}


func GetStringFromTemplate(templateValue string, strcuture interface{}) string{
	log.Println("Resolving template")
	var data bytes.Buffer
	temp := template.New("genericTemplate")
	temp, err := temp.Parse(templateValue)
	CheckNExitError(err)
	exeErr := temp.Execute(&data, strcuture)
	CheckNExitError(exeErr)
	log.Println("Rosolved template")
	return data.String()
}