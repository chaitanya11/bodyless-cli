package constants


const CONFIG_FILE_PERMISSIONS = 0777
const REPO = "https://github.com/chaitanya11/BodylessCMS"
const S3_INDEX_PAGE = "index.html"



// aws
const PROFILE_ENV_KEY = "AWS_PROFILE"
const COGNITO_PROVIDER_PREFIX = "cognito-idp.us-east-1.amazonaws.com/"
const COGNITO_POOL_NAME = "bodyless_pool"


// templates
const PROJECT_CONF_TEMPLATE = `
export class Config {
    public static readonly userPoolId= '{{.UserPoolId}}';
    public static readonly clientId = '{{.ClientId}}';
    public static readonly identityPoolId = '{{.IdentityPoolId}}';
    public static readonly awsRegion = '{{.AwsRegion}}';
}
`
type PROJECT_CONF_TEMPLATE_VARS struct {
	UserPoolId string
	ClientId string
	IdentityPoolId string
	AwsRegion string
}


// file paths
const CONFIG_DIR = ".bodyless-config"
const CONFIG_FILE_NAME = "config.json"
const PROJECT_CONFIG_PATH = "src/app/aws-services/config/index.ts";