package aws

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"bodyless-cli/utils"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"bodyless-cli/constants"
)

func createCognitoUserPool(poolName string, region *string) *string {
	// create userpool.
	createUserPoolInput := cognitoidentityprovider.CreateUserPoolInput{PoolName: &poolName}
	Cognitoidentityprovider := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: region,
	}))
	createUserPoolOutput, err := Cognitoidentityprovider.CreateUserPool(&createUserPoolInput)
	utils.CheckNExitError(err)
	userPoolId := createUserPoolOutput.UserPool.Id;
	return userPoolId
}

func createUserPoolClient(
	clientName string,
	userPoolId *string,
	region *string) *cognitoidentityprovider.CreateUserPoolClientOutput {
	generateSecrete := false
	Cognitoidentityprovider := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: region,
	}))
	createUserPoolClientInput := cognitoidentityprovider.CreateUserPoolClientInput{ClientName: &clientName,
		UserPoolId: userPoolId,
		GenerateSecret: &generateSecrete}
	createUserPoolClientOutput, err := Cognitoidentityprovider.CreateUserPoolClient(&createUserPoolClientInput)
	utils.CheckNExitError(err)
	return createUserPoolClientOutput
}

func createIdentityPool(poolName string, clientId *string, userPoolId string, region *string) *string {
	allowUnauthenticatedIdentities := false
	providerName := constants.COGNITO_PROVIDER_PREFIX + userPoolId

	provider := cognitoidentity.Provider{
		ClientId:             clientId,
		ProviderName:         &providerName,
		ServerSideTokenCheck: &allowUnauthenticatedIdentities,
	}
	providers := []*cognitoidentity.Provider{&provider}
	CognitoIdentity := cognitoidentity.New(session.New(&aws.Config{
		Region: region,
	}))
	createIdentityPoolInput := cognitoidentity.CreateIdentityPoolInput{IdentityPoolName: &poolName,
		AllowUnauthenticatedIdentities: &allowUnauthenticatedIdentities,
		CognitoIdentityProviders: providers}
	IdentityPool, err := CognitoIdentity.CreateIdentityPool(&createIdentityPoolInput)
	// TODO set roles for identity pool.
	utils.CheckNExitError(err)
	return IdentityPool.IdentityPoolId
}

func CreateCognitoResources(poolName string, path *string, region *string) constants.PROJECT_CONF_TEMPLATE_VARS {
	userPoolId := createCognitoUserPool(poolName, region)
	// create user pool client.
	createUserPoolClientOutput := createUserPoolClient(poolName, userPoolId, region)
	clientId := createUserPoolClientOutput.UserPoolClient.ClientId
	IdentityPoolId := createIdentityPool(poolName, clientId, *userPoolId, region)
	// write to config file of project.
	cognitoConfig := utils.WriteProjectConfig(userPoolId, clientId, IdentityPoolId, region, path)
	return cognitoConfig
}
