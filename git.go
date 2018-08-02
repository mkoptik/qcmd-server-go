package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
)

var lastCommit string

// Checks if destination directory exists. If yes, it pulls repository
// from origin and if directory does not exist, it will be cloned.
func updateGitRepository(repositoryURL string, cloneDirectory string) bool {

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

		pullOptions := &git.PullOptions{
			Progress:   os.Stdout,
			RemoteName: "origin",
		}

		if err := workTree.Pull(pullOptions); err != nil && err != git.NoErrAlreadyUpToDate {
			log.Fatal(err)
		}
	}

	commit := getLastCommitHash(path)
	changeDetected := commit != lastCommit

	// lastCommit is empty string on first run
	if lastCommit != "" && changeDetected {
		log.Printf("Git change detected. Old commit %s, new commit %s", lastCommit, commit)
	}

	lastCommit = commit
	return changeDetected
}

func getLastCommitHash(repoPath string) string {

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		log.Fatal(err)
	}

	return commit.Hash.String()
}
