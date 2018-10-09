package build_project

import (
	"bodyless-cli/utils"
	"log"
)

func BuildProj(path string) {
	log.Printf("Building project in directory %s ...", path)
	if path == "" {
		path = "."
	}

	// install project dependencies.
	log.Print("Installing project dependencies ...")
	utils.ExecuteCmd(path, "npm" ,"i")
	//log.Println(instDepOutput)

	// build project.
	utils.ExecuteCmd(path, "npm", "run", "build:prod")
	//log.Println(string(buildOutput))
	log.Printf("Building project in directory %s is completed", path)
}
