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
	Id string `json:"id"`
	Label string `json:"label"`
	Description string `json:"description"`
	CommandText string `json:"commandText"`
	Executive string `json:"executive"`
	Tags []string `json:"tags"`
}

type Tag struct {
	Path []string `json:"path"`
}

var commandsIndex bleve.Index = nil
var tagsIndex bleve.Index = nil

func main() {

	currentUser, _ := user.Current()

	mdFilesPath := filepath.Join(currentUser.HomeDir, ".qcmd/commands")
	commandsIndexPath := filepath.Join(currentUser.HomeDir, ".qcmd/commands.index.bleve")
	tagsIndexPath := filepath.Join(currentUser.HomeDir, ".qcmd/tags.index.bleve")

	updateGitRepository("https://github.com/mkoptik/qcmd-commands", mdFilesPath)
	absPath, err := filepath.Abs(mdFilesPath)
	if err != nil {
		log.Fatal(err)
	}

	commands, tags := readMarkdownFilesInPath(absPath, []string{}, [][]string{})

	indexCommands(commandsIndexPath, commands)
	defer commandsIndex.Close()

	indexTags(tagsIndexPath, tags)
	defer tagsIndex.Close()

	StartHttpServer()

	log.Printf("Closing bleve index")

}

func indexCommands(indexPath string, commands []*Command) {

	if _, err := os.Stat(indexPath); err == nil {
		err = os.RemoveAll(indexPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	enTextFieldMapping := bleve.NewTextFieldMapping()
	enTextFieldMapping.Analyzer = "en"

	documentMapping := bleve.NewDocumentMapping()
	documentMapping.AddFieldMappingsAt("label", enTextFieldMapping)
	documentMapping.AddFieldMappingsAt("description", enTextFieldMapping)
	documentMapping.AddFieldMappingsAt("tags", enTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("command", documentMapping)

	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		log.Fatal("Error creating bleve index for commands", err)
	}

	log.Printf("Indexing %d commands into %s", len(commands), indexPath)
	for i, command := range commands {
		index.Index(strconv.Itoa(i), command)
	}

	commandsIndex = index

}

func indexTags(indexPath string, uniqueTags [][]string) {

	if _, err := os.Stat(indexPath); err == nil {
		err = os.RemoveAll(indexPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	indexMapping := bleve.NewIndexMapping()

	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		log.Fatal("Error creating bleve index for tags", err)
	}

	log.Printf("Indexing %d tas into %s", len(uniqueTags), indexPath)
	for i, tag := range uniqueTags {
		tagObj := Tag{
			Path: tag,
		}
		index.Index(strconv.Itoa(i), tagObj)
	}

	tagsIndex = index
}