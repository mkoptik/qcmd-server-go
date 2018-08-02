package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/jasonlvhit/gocron"
)

type config struct {
}

// Command TODO: finish comment
type Command struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	CommandText string   `json:"commandText"`
	Executable  string   `json:"executable"`
	Tags        []string `json:"tags"`
}

// Tag TODO: finish comment
type Tag struct {
	Path []string `json:"path"`
}

var commandsIndex bleve.IndexAlias
var commandsIndexLastAdded bleve.Index

var tagsIndex bleve.IndexAlias
var tagsIndexLastAdded bleve.Index

func main() {

	fetchFromGitAndReindexCommandsIfNeeded()
	gocron.Every(5).Minutes().Do(fetchFromGitAndReindexCommandsIfNeeded)
	gocron.Start()

	defer commandsIndex.Close()
	defer tagsIndex.Close()

	startHTTPServer()

	log.Printf("Quiting, bye...")
}

func fetchFromGitAndReindexCommandsIfNeeded() {

	currentUser, _ := user.Current()

	mdFilesPath := filepath.Join(currentUser.HomeDir, ".qcmd/commands")
	commandsIndexPath := filepath.Join(currentUser.HomeDir, ".qcmd/commands.index.bleve")
	tagsIndexPath := filepath.Join(currentUser.HomeDir, ".qcmd/tags.index.bleve")

	changed := updateGitRepository("https://github.com/mkoptik/qcmd-commands", mdFilesPath)
	if !changed {
		return
	}

	absPath, err := filepath.Abs(mdFilesPath)
	if err != nil {
		log.Fatal(err)
	}

	commands, tags := readMarkdownFilesInPath(absPath, []string{}, [][]string{})

	indexCommands(commandsIndexPath, commands)
	indexTags(tagsIndexPath, tags)
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

	if commandsIndex == nil {
		commandsIndex = bleve.NewIndexAlias(index)
	} else {
		commandsIndex.Swap([]bleve.Index{index}, []bleve.Index{commandsIndexLastAdded})
	}
	commandsIndexLastAdded = index
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

	if tagsIndex == nil {
		tagsIndex = bleve.NewIndexAlias(index)
	} else {
		tagsIndex.Swap([]bleve.Index{index}, []bleve.Index{tagsIndexLastAdded})
	}
	tagsIndexLastAdded = index
}
