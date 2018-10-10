package aws

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"bodyless-cli/constants"
	"bodyless-cli/utils"
)

func CreateIamRoles(region *string, identityPoolId *string) (*string, *string) {
	Iam := iam.New(session.New(&aws.Config{
		Region: region,
	}))
	log.Println("Creating roles for identity pool ...")
	projectVars := constants.PROJECT_CONF_TEMPLATE_VARS{
		IdentityPoolId: *identityPoolId,
	}
	// valid role.
	validUserRoleDoc := utils.GetStringFromTemplate(constants.AUTHENTICATED_USER_ROLE_TRUST_POLICY_TEMPLATE,
		projectVars)
	validUserROleName := constants.AUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME
	validUserROleDes := constants.AUTHENTICATED_USER_ROLE_TRUST_POLICY_DESCRIPTION
	createValidUserRoleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: &validUserRoleDoc,
		RoleName: &validUserROleName,
		Description: &validUserROleDes}
	createRoleOut, createValidUserRoleErr := Iam.CreateRole(&createValidUserRoleInput)
	utils.CheckNExitError(createValidUserRoleErr)
	validUserRoleArn := createRoleOut.Role.Arn
	log.Printf("Created role for valid users, role arn %s", *validUserRoleArn)

	// invalid role.
	inValidUserRoleDoc := utils.GetStringFromTemplate(constants.UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_TEMPLATE,
		projectVars)
	inValidUserROleName := constants.UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME
	inValidUserROleDes := constants.UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_DESCRIPTION

	createinValidUserRoleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: &inValidUserRoleDoc,
		RoleName: &inValidUserROleName,
		Description: &inValidUserROleDes}
	createRoleOut, createinValidUserRoleErr := Iam.CreateRole(&createinValidUserRoleInput)
	utils.CheckNExitError(createinValidUserRoleErr)
	inValidUserRoleArn := createRoleOut.Role.Arn
	log.Printf("Created role for invalid users, role arn %s", *inValidUserRoleArn)

	// Adding policies to rolesAttaching roles to identity pool
	AttachPolicy(constants.AUTHENTICATED_USER_ROLE_POLICY_TEMPLATE,
		region,
		constants.AUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME,
		constants.AUTHENTICATED_USER_ROLE_POLICY_NAME)
	AttachPolicy(constants.UNAUTHENTICATED_USER_ROLE_POLICY_TEMPLATE,
		region,
		constants.UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME,
		constants.UNAUTHENTICATED_USER_ROLE_POLICY_NAME)
	return validUserRoleArn, inValidUserRoleArn
}


func AttachPolicy(
	policyTemplate string,
	region *string,
	roleName string,
	policyName string) *iam.PutRolePolicyOutput {
	log.Printf("Creating policy document for role %s ...", roleName)
	svc := iam.New(session.New(&aws.Config{
		Region: region,
	}))

	policyInput := iam.PutRolePolicyInput{RoleName: &roleName,
		PolicyDocument: &policyTemplate,
		PolicyName: &policyName,
	}
	PutRolePolicyOutput, error := svc.PutRolePolicy(&policyInput)
	utils.CheckNExitError(error)
	log.Printf("Created policy document for role %s ...", roleName)
	return PutRolePolicyOutput
}


func DeleteIamRole(roleName string, policyName string, region *string) {
	log.Printf("Deleting role %s ...", roleName)
	svc := iam.New(session.New(&aws.Config{
		Region: region,
	}))
	DetachRolePolicy(roleName, policyName, svc)
	_,err := svc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: &roleName,
	})
	utils.CheckNExitError(err)
	log.Printf("Deleted role %s.", roleName)
}

func DetachRolePolicy(roleName string, policyName string, svc *iam.IAM) {
	log.Printf("Detaching policy from role %s ...", roleName)
	_, err := svc.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
		RoleName: &roleName,
		PolicyName: &policyName,
	})
	utils.CheckNExitError(err)
	log.Printf("Detached policy from role %s.", roleName)
}