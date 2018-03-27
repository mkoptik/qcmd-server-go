package main

import (
	"github.com/blevesearch/bleve"
	"path/filepath"
	"log"
	"os"
	"os/user"
	"strconv"
)

type Config struct {

}

type Command struct {
	Id string
	Label string
	Description string
	CommandText string
	Tags []string
}

var bleveIndex bleve.Index = nil

func main() {

	currentUser, _ := user.Current()

	mdFilesPath := filepath.Join(currentUser.HomeDir, ".qcmd/commands")
	indexPath := filepath.Join(currentUser.HomeDir, ".qcmd/index.bleve")

	updateGitRepository("https://github.com/mkoptik/qcmd-commands", mdFilesPath)
	absPath, err := filepath.Abs(mdFilesPath)
	if err != nil {
		log.Fatal(err)
	}
	commands := readMarkdownFilesInPath(absPath, []string{})

	if _, err := os.Stat(indexPath); err == nil {
		err = os.RemoveAll(indexPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		log.Fatal("Error creating bleve index", err)
	}

	log.Printf("Indexing %d commands into %s", len(commands), indexPath)
	for i, command := range commands {
		index.Index(strconv.Itoa(i), command)
	}

	bleveIndex = index
	StartHttpServer()

	log.Printf("Closing bleve index")
	index.Close()

}