package utils

import (
	"os"
	"bodyless-cli/constants"
	"text/template"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func CheckNExitError(err error) {
	if err != nil {
		panic(err)
		os.Exit(1)
	}
}

func WriteProjectConfig(bucketName *string,
	userPoolId *string,
	clientId *string,
	identityPoolId *string,
	awsRegion *string,
	path *string,
	validRoleArn *string,
	invalidRoleArn *string,
	projectName *string) constants.PROJECT_CONF_TEMPLATE_VARS {
	templateVars := constants.PROJECT_CONF_TEMPLATE_VARS{
		BucketName: *bucketName,
		ProjectName: *projectName,
		UserPoolId: *userPoolId,
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
	configFile.Truncate(0)
	configFile.Seek(0,0)
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

func ExecuteCmd(executionDir string,name string, args ...string) {
	log.Printf("Executing %s in directory %s ...", name + " " + strings.Join(args, " "), executionDir)
	cmd := exec.Command("npm", args...)
	cmd.Dir = executionDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	log.Printf("Executing %s in directory %s is completed", name + " " + strings.Join(args, " "), executionDir)
}


func IsDirectory(path string) bool {
	fd, err := os.Stat(path)
	CheckNExitError(err)
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}

func RemoveDirectory(path string) {
	log.Printf("Removing directory %s ...", path)
	err := os.RemoveAll(path)
	CheckNExitError(err)
	log.Printf("Removed directory %s.", path)
}
