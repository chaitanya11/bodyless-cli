package utils

import (
	"os"
)

func CheckNExitError(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func WriteToProjectConf(clientId *string, IdentityPool *string, userPoolId *string) {

}