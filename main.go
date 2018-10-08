package main

import (
	"fmt"
	"os"
	"bodyless-cli/create-project"
	flag "github.com/ogier/pflag"
	"bodyless-cli/deploy-project"
	"bodyless-cli/build-project"
)

var (
	projectName   string
	path          string
	codeBucket    string
	profile       string
	region 		  string
	createCommand *flag.FlagSet
	buildCommand  *flag.FlagSet
	deployCommand *flag.FlagSet
)

func main() {
	// TODO add remove resources command to remove all created aws resources.
	commands := []*flag.FlagSet{
		createCommand,
		buildCommand,
		deployCommand,
	}

	if len(os.Args) == 1 {
		printHelp(commands)
		return
	}

	switch os.Args[1] {

	case "create":
		createCommand.Parse(os.Args[2:])
		create_project.CreateProj(
			projectName,
			path,
			codeBucket,
			profile,
			region)

	case "build":
		buildCommand.Parse(os.Args[2:])
		build_project.BuildProj()

	case "deploy":
		deployCommand.Parse(os.Args[2:])
		// build project.
		build_project.BuildProj()
		// deploy project.
		deploy_project.DeployProj()

	default:
		printHelp(commands)
		os.Exit(2)
	}

	for _, command := range commands {
		if command.Parsed() {
			if command.NFlag() == 0 {
				printDefaults(command)
			}
		}
	}
}

func init() {
	// create-project command
	createCommand = flag.NewFlagSet("create-project", flag.ExitOnError)
	createCommand.StringVarP(&projectName, "ProjectName", "N", "", "Name of the project.")
	createCommand.StringVarP(&path, "Path", "P", ".", "Project Location.")
	createCommand.StringVarP(&codeBucket, "CodeBucketName", "w", "", "Name of the bucket where website code is deployed.")
	createCommand.StringVarP(&profile, "profile", "p", "default", "Name of the aws profile configured.")
	createCommand.StringVarP(&region, "region", "r", "us-east-1", "Name of the aws region.")

	// build command
	buildCommand = flag.NewFlagSet("build", flag.ExitOnError)

	// deploy command
	deployCommand = flag.NewFlagSet("deploy", flag.ExitOnError)
}


func printDefaults(command *flag.FlagSet) {
	command.PrintDefaults()
}

func printHelp(commands []*flag.FlagSet) {
	fmt.Println("usage: bodyless <command> [<args>]")
	fmt.Println("commands:")
	fmt.Println("create")
	createCommand.PrintDefaults()
	fmt.Println("build")
	buildCommand.PrintDefaults()
	fmt.Println("deploy")
	deployCommand.PrintDefaults()
}