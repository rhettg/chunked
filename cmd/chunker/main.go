package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rhettg/sitter"
	treesitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func main() {
	parser := treesitter.NewParser()
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

	ch, err := sitter.NewChunks(golang.GetLanguage(), sourceCode, 256, 500)
	if err != nil {
		log.Fatal(err)
	}

	more := true
	var n []byte
	for more {
		n, more = ch.Next()
		fmt.Printf("Chunk (%d):\n", len(n))
		fmt.Println(string(n))
	}
}
