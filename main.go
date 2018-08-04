package main

import (
	"fmt"
	"os"

	"github.com/chaitanya11/bodyless/git"
	"github.com/chaitanya11/bodyless/s3"
	flag "github.com/ogier/pflag"
)

var (
	projectName   string
	path          string
	codeBucket    string
	contentBucket string
	profile       string
	createCommand *flag.FlagSet
	buildCommand  *flag.FlagSet
	deployCommand *flag.FlagSet
)

func main() {

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
	case "build":
		buildCommand.Parse(os.Args[2:])
	case "deploy":
		deployCommand.Parse(os.Args[2:])
	default:
		printHelp(commands)
		os.Exit(2)
	}
	// flag.Parse()
	// // if user does not supply flags, print usage
	// if flag.NFlag() == 0 {
	// 	printUsage()
	// }

	for _, command := range commands {
		if command.Parsed() {
			if command.NFlag() == 0 {
				printDefaults(command)
			}
		}
	}

	// clone cms repo in to path
	git.PullGitRepo("https://github.com/chaitanya11/BodylessCMS", path)

	// create Code Bucket
	s3.CreateBucket(codeBucket)
	s3.SetWebSiteConfig(codeBucket, "index", "index.html")

	// create content Bucket
	s3.CreateBucket(contentBucket)
}

func init() {
	// create command
	createCommand = flag.NewFlagSet("create", flag.ExitOnError)
	createCommand.StringVarP(&projectName, "ProjectName", "N", "", "Name of the project.")
	createCommand.StringVarP(&path, "Path", "P", ".", "Project Location.")
	createCommand.StringVarP(&codeBucket, "CodeBucketName", "w", "", "Name of the bucket where website code is deployed.")
	createCommand.StringVarP(&contentBucket, "ContentBucketName", "c", "", "Name of the bucket where content is stored.")
	createCommand.StringVarP(&profile, "profile", "p", "default", "Name of the aws profile configured.")

	buildCommand = flag.NewFlagSet("build", flag.ExitOnError)
	deployCommand = flag.NewFlagSet("deploy", flag.ExitOnError)
}

func printUsage() {
	// fmt.Printf("Usage: %s [options]\n", os.Args[0])
	// fmt.Println("Options:")
	// flag.PrintDefaults()

	if len(os.Args) == 1 {
		fmt.Println("usage: bodyless <command> [<args>]")
		return
	}
	os.Exit(1)
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
