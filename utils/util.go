package utils

import (
	"os"
	"bodyless-cli/constants"
	"text/template"
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
	path *string) constants.PROJECT_CONF_TEMPLATE_VARS {
	templateVars := constants.PROJECT_CONF_TEMPLATE_VARS{UserPoolId: *userPoolId,
		ClientId: *clientId,
		IdentityPoolId: *identityPoolId,
		AwsRegion: *awsRegion}
	template := template.New("projectConfig")
	template, parseErr := template.Parse(constants.PROJECT_CONF_TEMPLATE)
	CheckNExitError(parseErr)
	configFile, configFileReadErr := os.OpenFile(*path + "/" +constants.PROJECT_CONFIG_PATH,
		os.O_RDWR, os.ModePerm)
	CheckNExitError(configFileReadErr)
	exeErr := template.Execute(configFile, templateVars)
	CheckNExitError(exeErr)
	configFile.Close()
	return templateVars
}