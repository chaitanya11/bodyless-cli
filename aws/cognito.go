package aws

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"bodyless-cli/utils"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"bodyless-cli/constants"
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

func DeleteCognitoPool(userPoolId *string, region *string) {
	log.Printf("Deleting userpool %s ...", *userPoolId)
	Cognitoidentityprovider := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: region,
	}))
	_,err := Cognitoidentityprovider.DeleteUserPool(&cognitoidentityprovider.DeleteUserPoolInput{
		UserPoolId: userPoolId,
	})
	utils.CheckNExitError(err)
	log.Printf("Deleted userpool %s.", *userPoolId)
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
	// get roles arns
	validRole, invalidRole := CreateIamRoles(region, IdentityPool.IdentityPoolId);
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

func DeleteIdentityPool(identityPoolId *string, region *string) {
	log.Printf("Deleting Identity pool %s ...", *identityPoolId)
	CognitoIdentity := cognitoidentity.New(session.New(&aws.Config{
		Region: region,
	}))
	_,err :=CognitoIdentity.DeleteIdentityPool(&cognitoidentity.DeleteIdentityPoolInput{
		IdentityPoolId: identityPoolId,
	})
	utils.CheckNExitError(err)
	log.Printf("Deleted Identity pool %s ...", *identityPoolId)
}

func createUser(region *string, userName string, password string, email string, userPoolId *string, clientId *string) {
	log.Println("Creating user in cognito pool ...")
	CognitoSvc := cognitoidentityprovider.New(session.New(&aws.Config{
		Region: region,
	}))
	attributeList := []*cognitoidentityprovider.AttributeType{
		&cognitoidentityprovider.AttributeType{
			Name: aws.String("email"),
			Value: &email,
		},
	}
	_, userCreateErr := CognitoSvc.SignUp(&cognitoidentityprovider.SignUpInput{
		Username: &userName,
		Password: &password,
		UserAttributes: attributeList,
		ClientId: clientId,
	})
	utils.CheckNExitError(userCreateErr)
	log.Printf("User created with details, username : %s, password: %s, email: %s.",
		userName, password, email)
	log.Println("Use username and password to login to application.")

	// confirm created user.
	_, userCnfmErr := CognitoSvc.AdminConfirmSignUp(&cognitoidentityprovider.AdminConfirmSignUpInput{
		Username: &userName,
		UserPoolId: userPoolId,
	})
	utils.CheckNExitError(userCnfmErr)
}

func CreateCognitoResources(poolName string,
	path *string,
	region *string,
	bucketName *string,
	projectName *string) constants.PROJECT_CONF_TEMPLATE_VARS {
	log.Println("Creating cognito resources ...")
	userPoolId := createCognitoUserPool(poolName, region)
	// create user pool client.
	createUserPoolClientOutput := createUserPoolClient(poolName, userPoolId, region)
	clientId := createUserPoolClientOutput.UserPoolClient.ClientId
	IdentityPoolId, validRoleArn, inValidRoleArn := createIdentityPool(poolName, clientId, *userPoolId, region)
	// creating user in cognito pool.
	createUser(region, "bodyless", "Hello@1234", "bodylesscms@mailinator.com", userPoolId, clientId)

	// write to config file of project.
	log.Println("Writing project configuration ...")
	cognitoConfig := utils.WriteProjectConfig(bucketName, userPoolId, clientId, IdentityPoolId, region, path,
		validRoleArn, inValidRoleArn, projectName)
	return cognitoConfig
}
