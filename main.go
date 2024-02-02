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

	ch, err := NewChunks(golang.GetLanguage(), sourceCode, 256)
	if err != nil {
		log.Fatal(err)
	}

	more := true
	var n []byte
	for more {
		n, more = ch.Next()
		fmt.Println("Chunk:")
		fmt.Println(string(n))
	}
}

type Chunks struct {
	sourceCode []byte
	size       uint32
	c          *sitter.TreeCursor
}

func NewChunks(l *sitter.Language, sourceCode []byte, chunkSize uint32) (*Chunks, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(l)

	tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
	if err != nil {
		return nil, err
	}

	c := sitter.NewTreeCursor(tree.RootNode())
	if !c.GoToFirstChild() {
		return nil, fmt.Errorf("no first child")
	}

	ch := Chunks{
		sourceCode: sourceCode,
		size:       chunkSize,
		c:          c,
	}

	return &ch, nil
}

func (c *Chunks) Next() ([]byte, bool) {
	n := c.c.CurrentNode()
	//fmt.Println(n)
	start := n.StartByte()
	end := n.EndByte()

	more := false
	for c.c.GoToNextSibling() {
		if c.c.CurrentNode().EndByte()-start > c.size {
			more = true
			break
		}

		end = c.c.CurrentNode().EndByte()
	}
	return c.sourceCode[start:end], more
}
