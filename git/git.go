package git

import (
	"fmt"
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

// git
func PullGitRepo(url string, directory string) {

	// Clone the given repository to the given directory
	Info("git clone %s %s --recursive", url, directory)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	checkIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	checkIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	checkIfError(err)

	fmt.Println(commit)
}

// CheckIfError should be used to naively panics if an error is not nil.
func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
