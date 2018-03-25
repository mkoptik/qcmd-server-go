package main

import (
	"net/http"
	"log"
	"github.com/blevesearch/bleve"
	"encoding/json"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {

	searchString := r.URL.Query().Get("search")

	query := bleve.NewMatchQuery(searchString)
	request := bleve.NewSearchRequest(query)
	request.Fields = []string { "Label", "CommandText", "Description" }

	searchResult, err := bleveIndex.Search(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d commands in %s", searchResult.Hits.Len(), searchResult.Took)

	foundCommands := make([]Command, searchResult.Hits.Len())
	for i, hit := range searchResult.Hits {
		foundCommands[i] = Command{
			Label: hit.Fields["Label"].(string),
			CommandText: hit.Fields["CommandText"].(string),
		}
		if hit.Fields["Description"] != nil {
			foundCommands[i].Description = hit.Fields["Description"].(string)
		}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(foundCommands)
}

func StartHttpServer() {
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}