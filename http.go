package main

import (
	"net/http"
	"log"
	"github.com/blevesearch/bleve"
	"encoding/json"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchString := r.URL.Query().Get("search")

	termsQuery := bleve.NewDisjunctionQuery()

	tokens := bleveIndex.Mapping().AnalyzerNamed("standard").Analyze([]byte(searchString))
	for _, token := range tokens {
		tokenTerm := string(token.Term)
		fieldsQuery := bleve.NewDisjunctionQuery()

		prefixQuery := bleve.NewPrefixQuery(tokenTerm)
		prefixQuery.SetField("label")
		prefixQuery.SetBoost(4)
		fieldsQuery.AddQuery(prefixQuery)

		prefixQuery = bleve.NewPrefixQuery(tokenTerm)
		prefixQuery.SetField("description")
		prefixQuery.SetBoost(2)
		fieldsQuery.AddQuery(prefixQuery)

		termsQuery.AddQuery(fieldsQuery)
	}

	request := bleve.NewSearchRequest(termsQuery)
	request.Fields = []string { "label", "commandText", "description", "tags" }
	request.Size = 500

	searchResult, err := bleveIndex.Search(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d commands in %s", searchResult.Hits.Len(), searchResult.Took)

	foundCommands := make([]Command, searchResult.Hits.Len())
	for i, hit := range searchResult.Hits {
		foundCommands[i] = Command{
			Label: hit.Fields["label"].(string),
			CommandText: hit.Fields["commandText"].(string),
		}
		if hit.Fields["description"] != nil {
			foundCommands[i].Description = hit.Fields["description"].(string)
		}

		// SINGLE VALUE IS NOT STORED AS ARRAY IN BLEVE SEARCH
		tagsString, ok := hit.Fields["tags"].(string)
		if ok {
			foundCommands[i].Tags = []string {tagsString}
		} else {
			tags := make([]string, len(hit.Fields["tags"].([]interface{})))
			for i2, tagObj := range hit.Fields["tags"].([]interface{}) {
				tags[i2] = tagObj.(string)
			}
			foundCommands[i].Tags = tags
		}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(foundCommands)
}

func StartHttpServer() {
	http.HandleFunc("/search", searchHandler)
	log.Printf("Starting http server on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}