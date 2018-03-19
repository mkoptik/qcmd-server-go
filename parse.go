package main

import (
	"io/ioutil"
	"log"
	"strings"
	"path"
	"fmt"
	"github.com/russross/blackfriday"
)

func readMarkdownFilesInPath(absDirPath string, tags []string) {
	if !path.IsAbs(absDirPath) {
		panic(fmt.Sprintf("Path %s is not absolute", absDirPath))
	}
	files, err := ioutil.ReadDir(absDirPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "README") {
			continue
		}
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			readMarkdownFilesInPath(path.Join(absDirPath, file.Name()), append(tags, file.Name()))
		}
		if strings.HasSuffix(file.Name(), ".md") {
			parseMarkdown(path.Join(absDirPath, file.Name()))
		}
	}
}

func parseMarkdown(absPath string) {
	log.Printf("Parsing %s", absPath)
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	markdown := blackfriday.New()
	node := markdown.Parse(bytes)

	log.Printf("Parsed %s", node)
}