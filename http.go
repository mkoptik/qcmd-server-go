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
	tokens := commandsIndex.Mapping().AnalyzerNamed("standard").Analyze([]byte(searchString))
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

	request := &bleve.SearchRequest{}
	if len(tokens) > 0 {
		request = bleve.NewSearchRequest(termsQuery)
	} else {
		request = bleve.NewSearchRequest(bleve.NewMatchAllQuery())
	}

	request.Fields = []string { "label", "commandText", "description", "tags" }
	request.Size = 500

	searchResult, err := commandsIndex.Search(request)
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


	}

	encoder := json.NewEncoder(w)
	encoder.Encode(foundCommands)
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	matchAllQuery := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(matchAllQuery)

	searchRequest.Fields = []string { "path" }
	searchRequest.Size = 1000

	searchResults, err := tagsIndex.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	foundTags := make([][]string, searchResults.Hits.Len())
	for i, hit := range searchResults.Hits {
		// SINGLE VALUE IS NOT STORED AS ARRAY IN BLEVE SEARCH
		pathString, ok := hit.Fields["path"].(string)
		if ok {
			foundTags[i] = []string {pathString}
		} else {
			tags := make([]string, len(hit.Fields["path"].([]interface{})))
			for i2, tagObj := range hit.Fields["path"].([]interface{}) {
				tags[i2] = tagObj.(string)
			}
			foundTags[i] = tags
		}

	}

	encoder := json.NewEncoder(w)
	encoder.Encode(foundTags)
}

func StartHttpServer() {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/tags", tagsHandler)
	log.Printf("Starting http server on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}