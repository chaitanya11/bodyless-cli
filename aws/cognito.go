package aws

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"bodyless-cli/utils"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"bodyless-cli/constants"
)

func createCognitoUserPool(poolName string) *string {
	// create userpool.
	createUserPoolInput := cognitoidentityprovider.CreateUserPoolInput{PoolName: &poolName}
	Cognitoidentityprovider := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))
	createUserPoolOutput, err := Cognitoidentityprovider.CreateUserPool(&createUserPoolInput)
	utils.CheckNExitError(err)
	userPoolId := createUserPoolOutput.UserPool.Id;
	return userPoolId
}

func createUserPoolClient(
	clientName string,
	userPoolId *string) *cognitoidentityprovider.CreateUserPoolClientOutput {
	generateSecrete := false
	Cognitoidentityprovider := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))
	createUserPoolClientInput := cognitoidentityprovider.CreateUserPoolClientInput{ClientName: &clientName,
		UserPoolId: userPoolId,
		GenerateSecret: &generateSecrete}
	createUserPoolClientOutput, err := Cognitoidentityprovider.CreateUserPoolClient(&createUserPoolClientInput)
	utils.CheckNExitError(err)
	return createUserPoolClientOutput
}

func createIdentityPool(poolName string, clientId *string, userPoolId string) *string {
	allowUnauthenticatedIdentities := false
	providerName := constants.COGNITO_PROVIDER_PREFIX + userPoolId

	provider := cognitoidentity.Provider{
		ClientId:             clientId,
		ProviderName:         &providerName,
		ServerSideTokenCheck: &allowUnauthenticatedIdentities,
	}
	providers := []*cognitoidentity.Provider{&provider}
	CognitoIdentity := cognitoidentity.New(session.New(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))
	createIdentityPoolInput := cognitoidentity.CreateIdentityPoolInput{IdentityPoolName: &poolName,
		AllowUnauthenticatedIdentities: &allowUnauthenticatedIdentities,
		CognitoIdentityProviders: providers}
	IdentityPool, err := CognitoIdentity.CreateIdentityPool(&createIdentityPoolInput)
	utils.CheckNExitError(err)
	return IdentityPool.IdentityPoolId
}

func CreateCognitoResources(poolName string) {
	userpoolId := createCognitoUserPool(poolName)
	// create user pool client.
	createUserPoolClientOutput := createUserPoolClient(poolName, userpoolId)
	clientId := createUserPoolClientOutput.UserPoolClient.ClientId
	IdentityPoolId := createIdentityPool(poolName, clientId, *userpoolId)
	// write to config file of project.
	utils.WriteToProjectConf(clientId, IdentityPoolId, userpoolId)
}
