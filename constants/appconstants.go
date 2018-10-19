package constants


const CONFIG_FILE_PERMISSIONS = 0777
const REPO = "https://github.com/chaitanya11/BodylessCMS"
const S3_INDEX_PAGE = "index.html"
const S3_STYLE_PAGE = "style.css"
const S3_IMG_FILE = "cow.png"
const S3_IMG_PATH = "content/img/"
const ASSESTS = "assets"



// aws
const PROFILE_ENV_KEY = "AWS_PROFILE"
const COGNITO_PROVIDER_PREFIX_TEMPLATE = "cognito-idp.{{.AwsRegion}}.amazonaws.com/"
const COGNITO_POOL_NAME = "bodyless_pool"
const AUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME = "Cognito_bodylesscms_identity_poolAuth_Role"
const AUTHENTICATED_USER_ROLE_POLICY_NAME = "Cognito_bodylesscms_identity_poolAuth_Policy"
const AUTHENTICATED_USER_ROLE_TRUST_POLICY_DESCRIPTION = "This role is applied for authenticated users from cognito."
const UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_NAME = "Cognito_bodylesscms_identity_poolUnauth_Role"
const UNAUTHENTICATED_USER_ROLE_POLICY_NAME = "Cognito_bodylesscms_identity_poolUnauth_Policy"
const UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_DESCRIPTION = "This role is applied for unAuthenticated users from cognito."

// templates
const PROJECT_CONF_TEMPLATE = `export class Config {
    public static readonly userPoolId= '{{.UserPoolId}}';
    public static readonly clientId = '{{.ClientId}}';
    public static readonly identityPoolId = '{{.IdentityPoolId}}';
    public static readonly awsRegion = '{{.AwsRegion}}';
    public static readonly bucketname = '{{.BucketName}}';
    public static readonly projectName = '{{.ProjectName}}';
}
`
const AUTHENTICATED_USER_ROLE_TRUST_POLICY_TEMPLATE = `{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Principal": {
			"Federated": "cognito-identity.amazonaws.com"
		},
		"Action": "sts:AssumeRoleWithWebIdentity",
		"Condition": {
			"StringEquals": {
				"cognito-identity.amazonaws.com:aud": "{{.IdentityPoolId}}"
			},
			"ForAnyValue:StringLike": {
				"cognito-identity.amazonaws.com:amr": "authenticated"
			}
		}
	}]
}`
const UNAUTHENTICATED_USER_ROLE_TRUST_POLICY_TEMPLATE = `{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Principal": {
			"Federated": "cognito-identity.amazonaws.com"
		},
		"Action": "sts:AssumeRoleWithWebIdentity",
		"Condition": {
			"StringEquals": {
				"cognito-identity.amazonaws.com:aud": "{{.IdentityPoolId}}"
			},
			"ForAnyValue:StringLike": {
				"cognito-identity.amazonaws.com:amr": "unauthenticated"
			}
		}
	}]
}`

const AUTHENTICATED_USER_ROLE_POLICY_TEMPLATE = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "mobileanalytics:PutEvents",
        "cognito-sync:*",
		"s3:*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}`

const PUBLIC_BUCKET_POLICY_TEMPLATE = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::{{.BucketName}}/*"
        }
    ]
}`


const UNAUTHENTICATED_USER_ROLE_POLICY_TEMPLATE = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "mobileanalytics:PutEvents",
        "cognito-sync:*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}`

// structs
type PROJECT_CONF_TEMPLATE_VARS struct {
	BucketName string
	ProjectName string
	UserPoolId string
	ClientId string
	IdentityPoolId string
	AwsRegion string
	ValidRoleArn string
	InValidRoleArn string
}
type BodylessProjectConfig struct {
	BucketName string
	Region string
	Profile string
	CognitoConfig PROJECT_CONF_TEMPLATE_VARS
}


// file paths
const CONFIG_DIR = ".bodyless-config"
const CONFIG_FILE_NAME = "config.json"
const PROJECT_CONFIG_PATH = "src/app/services/aws/config/index.ts";
const BUILD_FILES_PATH = "dist/BodylessCMS/"