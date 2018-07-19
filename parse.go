package main

import (
	"io/ioutil"
	"log"
	"strings"
	"path"
	"fmt"
	"gopkg.in/russross/blackfriday.v2"
	"reflect"
)

func readMarkdownFilesInPath(absDirPath string, tagsStack []string, uniqueTags [][]string) ([]*Command, [][]string) {
	if !path.IsAbs(absDirPath) {
		panic(fmt.Sprintf("Path %s is not absolute", absDirPath))
	}
	files, err := ioutil.ReadDir(absDirPath)
	if err != nil {
		log.Fatal(err)
	}

	var allCommands []*Command

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "README") {
			continue
		}
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			commands, tags := readMarkdownFilesInPath(path.Join(absDirPath, file.Name()), append(tagsStack, file.Name()), uniqueTags)
			allCommands = append(allCommands, commands...)
			uniqueTags = tags
		}
		if strings.HasSuffix(file.Name(), ".md") {
			// COMMAND EXECUTABLE IS USED AS TAG INSTEAD OF FILE NAME
			commands := parseMarkdown(path.Join(absDirPath, file.Name()), tagsStack)
			for _, command := range commands {
				uniqueTags = addTagsIfNotExists(append(tagsStack, command.Executable), uniqueTags)
			}
			allCommands = append(allCommands, commands...)
		}
	}

	return allCommands, uniqueTags
}

func parseMarkdown(absPath string, tags []string) []*Command {
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal(err)
	}

	markdown := blackfriday.New(blackfriday.WithExtensions(blackfriday.FencedCode))
	node := markdown.Parse(bytes)

	return parseDocumentHeading(node.FirstChild, tags)
}

func parseDocumentHeading(node *blackfriday.Node, tags []string) []*Command {
	if node.Type != blackfriday.Heading || node.HeadingData.Level != 1 {
		log.Fatal("Node is not Heading with level 1")
	}

	var commands []*Command
	var groupName string
	var currentNode = node.Next

	for currentNode != nil {
		if currentNode.Type == blackfriday.Heading {
			if currentNode.HeadingData.Level == 2 {
				groupName = getTextFromNode(node.Next)
			} else if currentNode.HeadingData.Level == 3 {
				command, commandLastNode := parseCommandHeading(currentNode)
				currentNode = commandLastNode
				command.Tags = append(command.Tags, tags...)
				commands = append(commands, command)
				if groupName != "" {
					command.Tags = append(command.Tags, groupName)
				}
			}
		}
		currentNode = currentNode.Next
	}

	return commands
}

func parseCommandHeading(node *blackfriday.Node) (*Command, *blackfriday.Node) {
	if node.Type != blackfriday.Heading || node.HeadingData.Level != 3 {
		log.Fatal("Command heading node must be Heading with level 3")
	}

	command := Command{
		Label: getTextFromNode(node),
	}

	// Command is mandatory code block right after label
	node = node.Next
	if node.Type != blackfriday.CodeBlock {
		log.Fatalf("Node after heading of %s is not CodeBlock", command.Label)
	}
	command.CommandText = getTextFromNode(node)

	setExecutableFromCommandText(&command)

	// Command description is optional text after code
	if node.Next != nil {
		if node.Next.Type == blackfriday.Paragraph {
			node = node.Next
			command.Description = getTextFromNode(node)
		}
	}

	return &command, node

}

// Reads command text and get command executive. Skips sudo
func setExecutableFromCommandText(command *Command) {
	for _, field := range strings.Fields(command.CommandText) {
		if field != "sudo" {
			command.Executable = field
			return
		}
	}

}

func addTagsIfNotExists(commandTags []string, uniqueTags [][]string) [][]string {
	for _, tags := range uniqueTags {
		if reflect.DeepEqual(commandTags, tags) {
			return uniqueTags
		}
	}
	return append(uniqueTags, commandTags)
}

func getTextFromNode(node *blackfriday.Node) string {
	if node.Literal != nil {
		return string(node.Literal)
	} else if node.FirstChild != nil && node.FirstChild.Type == blackfriday.Text {
		return string(node.FirstChild.Literal)
	} else {
		log.Fatal("Cannot get text from node")
		return ""
	}
}