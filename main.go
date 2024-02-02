package main

import (
	"fmt"
	"io"
	"log"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func main() {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	f, _ := os.Open("main.go")
	sourceCode, _ := io.ReadAll(f)
	tree := parser.Parse(nil, sourceCode)

	c := sitter.NewTreeCursor(tree.RootNode())
	if !c.GoToFirstChild() {
		log.Fatal("no first child")
	}
	for {
		n := c.CurrentNode()
		fmt.Println(n.Type())
		if n.Type() == "function_declaration" {
			fmt.Println(n.Symbol())
			fmt.Println(n.EndPoint().Row - n.StartPoint().Row)
			//fmt.Println(n.Content(sourceCode))
			// find the identifier
			i := n.ChildByFieldName("name")
			if i == nil {
				fmt.Println("no identifier")
			} else {
				fmt.Println(i.Content(sourceCode))
			}
		}
		if !c.GoToNextSibling() {
			break
		}
	}

	os.Exit(0)
}
