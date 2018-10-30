package main

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"log"
	"net/http"
	"os"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	commandsQuery := buildCommandQuery(r)
	tagsQuery := buildTagsQuery(r)

	request := &bleve.SearchRequest{}
	rootQuery := bleve.NewConjunctionQuery(commandsQuery, tagsQuery)
	request = bleve.NewSearchRequest(rootQuery)

	request.Fields = []string{"label", "commandText", "description", "executable", "tags"}
	request.Size = 500

	searchResult, err := commandsIndex.Search(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d commands in %s", searchResult.Hits.Len(), searchResult.Took)

	foundCommands := make([]Command, searchResult.Hits.Len())
	for i, hit := range searchResult.Hits {
		foundCommands[i] = Command{
			Label:       hit.Fields["label"].(string),
			CommandText: hit.Fields["commandText"].(string),
			Executable:  hit.Fields["executable"].(string),
			Tags:        extractStringsArray(hit, "tags"),
		}
		if hit.Fields["description"] != nil {
			foundCommands[i].Description = hit.Fields["description"].(string)
		}
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	encoder := json.NewEncoder(w)
	encoder.Encode(foundCommands)
}

func buildCommandQuery(r *http.Request) query.Query {
	searchString := r.URL.Query().Get("search")

	termsQuery := bleve.NewDisjunctionQuery()
	tokens := commandsIndex.Mapping().AnalyzerNamed("standard").Analyze([]byte(searchString))
	if len(tokens) == 0 {
		return bleve.NewMatchAllQuery()
	}

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

		prefixQuery = bleve.NewPrefixQuery(tokenTerm)
		prefixQuery.SetField("executable")
		prefixQuery.SetBoost(8)
		fieldsQuery.AddQuery(prefixQuery)

		termsQuery.AddQuery(fieldsQuery)
	}

	return termsQuery
}

func buildTagsQuery(r *http.Request) query.Query {
	searchTagsString := r.URL.Query().Get("tag")
	if len(searchTagsString) > 0 {
		query := bleve.NewMatchQuery(searchTagsString)
		query.SetField("tags")
		return query
	}
	return bleve.NewMatchAllQuery()
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	matchAllQuery := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(matchAllQuery)

	searchRequest.Fields = []string{"path"}
	searchRequest.Size = 1000

	searchResults, err := tagsIndex.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	foundTags := make([][]string, searchResults.Hits.Len())
	for i, hit := range searchResults.Hits {
		foundTags[i] = extractStringsArray(hit, "path")
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	encoder := json.NewEncoder(w)
	encoder.Encode(foundTags)
}

// SINGLE VALUE IS NOT STORED AS ARRAY IN BLEVE SEARCH
func extractStringsArray(hit *search.DocumentMatch, fieldName string) []string {
	if hit.Fields[fieldName] == nil {
		return nil
	}
	singleString, ok := hit.Fields[fieldName].(string)
	if ok {
		return []string{singleString}
	}
	tags := make([]string, len(hit.Fields[fieldName].([]interface{})))
	for i2, tagObj := range hit.Fields[fieldName].([]interface{}) {
		tags[i2] = tagObj.(string)
	}
	return tags
}

func startHTTPServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	fsJs := http.FileServer(http.Dir("./web/client/dist"))
	fsAssets := http.FileServer(http.Dir("./web/assets"))

	http.Handle("/", http.StripPrefix("/", fsJs))
	http.Handle("/assets", http.StripPrefix("/", fsAssets))

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/tags", tagsHandler)
	log.Printf("Starting http server on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
