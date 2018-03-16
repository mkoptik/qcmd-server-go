package main

import (
	//"github.com/blevesearch/bleve"
	//"log"
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

func main() {

	/*command1 := Command{
		Id: "id1",
	 	Label: "label1",
	 	CommandText: "command text 1",
	 	Description: "description1",
	 	Tags: []string{"tag1", "tag2"},
	}

	command2 := Command{
		Id: "id2",
		Label: "label2",
		CommandText: "command text 2",
		Description: "description1",
		Tags: []string{"tag1", "tag2"},
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("/home/mkoptik/.qcmd/index.bleve", mapping)
	if err != nil {
		log.Fatal("Error creating bleve index", err)
	}

	index.Index("doc1", command1)
	index.Index("doc2", command2)

	index.Close()

	log.Print("Finished")*/

	updateGitRepository("https://github.com/mkoptik/qcmd-commands", "./data/qcmd-commands")

}