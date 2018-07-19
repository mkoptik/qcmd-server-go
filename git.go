package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
)

// Checks if destination directory exists. If yes, it pulls repository
// from origin and if directory does not exist, it will be cloned.
func updateGitRepository(repositoryURL string, cloneDirectory string) {

	path, err := filepath.Abs(cloneDirectory)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {

		log.Printf("Directory %s does not exist, creating new", path)
		err = os.MkdirAll(path, os.ModePerm)

		log.Printf("Cloning %s into %s", repositoryURL, path)
		git.PlainClone(path, false, &git.CloneOptions{
			URL:      repositoryURL,
			Progress: os.Stdout,
		})

	} else {

		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Fatal(err)
		}

		workTree, err := repo.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Git pull from origin")

		pullOptions := &git.PullOptions{
			Progress:   os.Stdout,
			RemoteName: "origin",
		}

		if err := workTree.Pull(pullOptions); err != nil && err != git.NoErrAlreadyUpToDate {
			log.Fatal(err)
		}

	}

}
