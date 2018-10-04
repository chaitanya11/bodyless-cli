package aws

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"bodyless-cli/utils"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"bodyless-cli/constants"
	"github.com/aws/aws-sdk-go/service/iam"
	"log"
)

func createCognitoUserPool(poolName string, region *string) *string {
	log.Println("Createing cognito userpool ....")
	// create userpool.
	createUserPoolInput := cognitoidentityprovider.CreateUserPoolInput{PoolName: &poolName}
	Cognitoidentityprovider := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: region,
	}))
	createUserPoolOutput, err := Cognitoidentityprovider.CreateUserPool(&createUserPoolInput)
	utils.CheckNExitError(err)
	userPoolId := createUserPoolOutput.UserPool.Id;
	log.Printf("Created cognito userpool. userpool id %s", *userPoolId)
	return userPoolId
}

func createUserPoolClient(
	clientName string,
	userPoolId *string,
	region *string) *cognitoidentityprovider.CreateUserPoolClientOutput {
		log.Println("Creating cognito userpool client ...")
	generateSecrete := false
	Cognitoidentityprovider := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: region,
	}))
	createUserPoolClientInput := cognitoidentityprovider.CreateUserPoolClientInput{ClientName: &clientName,
		UserPoolId: userPoolId,
		GenerateSecret: &generateSecrete}
	createUserPoolClientOutput, err := Cognitoidentityprovider.CreateUserPoolClient(&createUserPoolClientInput)
	utils.CheckNExitError(err)
	log.Printf("Created cognito userpool client, client id %s",
		*createUserPoolClientOutput.UserPoolClient.ClientId)
	return createUserPoolClientOutput
}

func createIdentityPool(poolName string, clientId *string, userPoolId string, region *string) (*string, *string, *string) {
	log.Println("Creating cognito identity pool ...")
	allowUnauthenticatedIdentities := false
	providerName := utils.GetStringFromTemplate(constants.COGNITO_PROVIDER_PREFIX_TEMPLATE,
		constants.PROJECT_CONF_TEMPLATE_VARS{
			AwsRegion: *region,
		}) + userPoolId

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
	IdentityPool, createPoolErr := CognitoIdentity.CreateIdentityPool(&createIdentityPoolInput)
	utils.CheckNExitError(createPoolErr)
	log.Printf("Created cognito identity pool, identity pool id %s", *IdentityPool.IdentityPoolId)
	// TODO set roles with permissions for identity pool.
	// get roles arns
	validRole, invalidRole := createIamRoles(region, IdentityPool.IdentityPoolId);
	log.Println("Attaching roles to identity pool ....")
	_, setRolesErr := CognitoIdentity.SetIdentityPoolRoles(&cognitoidentity.SetIdentityPoolRolesInput{
		IdentityPoolId: IdentityPool.IdentityPoolId,
		Roles:          map[string]*string{
			"authenticated" : validRole,
			"unauthenticated" : invalidRole,
		},
	})
	utils.CheckNExitError(setRolesErr)
	log.Println("Roles attached to identity pool")
	return IdentityPool.IdentityPoolId, validRole, invalidRole
}

func createIamRoles(region *string, identityPoolId *string) (*string, *string) {
	Iam := iam.New(session.New(&aws.Config{
		Region: region,
	}))
	log.Println("Creating roles for identity pool ...")
	projectVars := constants.PROJECT_CONF_TEMPLATE_VARS{
		IdentityPoolId: *identityPoolId,
	}
	// valid role.
	validUserRoleDoc := utils.GetStringFromTemplate(constants.AUTHENTICATED_USER_ROLE_POLICY_TEMPLATE,
		projectVars)
	validUserROleName := constants.AUTHENTICATED_USER_ROLE_POLICY_NAME
	validUserROleDes := constants.AUTHENTICATED_USER_ROLE_POLICY_DESCRIPTION
	createValidUserRoleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: &validUserRoleDoc,
		RoleName: &validUserROleName,
		Description: &validUserROleDes}
	createRoleOut, createValidUserRoleErr := Iam.CreateRole(&createValidUserRoleInput)
	utils.CheckNExitError(createValidUserRoleErr)
	validUserRoleArn := createRoleOut.Role.Arn
	log.Printf("Created role for valid users, role arn %s", *validUserRoleArn)

	// invalid role.
	inValidUserRoleDoc := utils.GetStringFromTemplate(constants.UNAUTHENTICATED_USER_ROLE_POLICY_TEMPLATE,
		projectVars)
	inValidUserROleName := constants.UNAUTHENTICATED_USER_ROLE_POLICY_NAME
	inValidUserROleDes := constants.UNAUTHENTICATED_USER_ROLE_POLICY_DESCRIPTION

	createinValidUserRoleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: &inValidUserRoleDoc,
		RoleName: &inValidUserROleName,
		Description: &inValidUserROleDes}
	createRoleOut, createinValidUserRoleErr := Iam.CreateRole(&createinValidUserRoleInput)
	utils.CheckNExitError(createinValidUserRoleErr)
	inValidUserRoleArn := createRoleOut.Role.Arn
	log.Printf("Created role for invalid users, role arn %s", *inValidUserRoleArn)
	return validUserRoleArn, inValidUserRoleArn
}

func CreateCognitoResources(poolName string, path *string, region *string) constants.PROJECT_CONF_TEMPLATE_VARS {
	log.Println("Creating cognito resources ...")
	userPoolId := createCognitoUserPool(poolName, region)
	// create user pool client.
	createUserPoolClientOutput := createUserPoolClient(poolName, userPoolId, region)
	clientId := createUserPoolClientOutput.UserPoolClient.ClientId
	IdentityPoolId, validRoleArn, inValidRoleArn := createIdentityPool(poolName, clientId, *userPoolId, region)
	// write to config file of project.
	log.Println("Writing project configuration ...")
	cognitoConfig := utils.WriteProjectConfig(userPoolId, clientId, IdentityPoolId, region, path,
		validRoleArn, inValidRoleArn)
	return cognitoConfig
}
