package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// The main function
func main() {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <file>")
	}
	fName := os.Args[1]

	f, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}
	sourceCode, _ := io.ReadAll(f)
	tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
	if err != nil {
		log.Fatal(err)
	}

	c := sitter.NewTreeCursor(tree.RootNode())
	if !c.GoToFirstChild() {
		log.Fatal("no first child")
	}

	for {
		n := c.CurrentNode()
		length := n.EndPoint().Row - n.StartPoint().Row + 1

		name := "anonymous"
		i := n.ChildByFieldName("name")
		if i != nil {
			name = i.Content(sourceCode)
		}

		if n.Type() != "\n" && n.Type() != "comment" {
			fmt.Printf("%s %s (%d lines)\n", n.Type(), name, length)
		}

		if !c.GoToNextSibling() {
			break
		}
	}

	os.Exit(0)
}
