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

	ch, err := NewChunks(golang.GetLanguage(), sourceCode, 256, false)
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

	offset      uint32
	splitChunks bool
}

func NewChunks(l *sitter.Language, sourceCode []byte, chunkSize uint32, splitChunks bool) (*Chunks, error) {
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

		splitChunks: splitChunks,
	}

	return &ch, nil
}

func (c *Chunks) Next() ([]byte, bool) {
	n := c.c.CurrentNode()
	start := n.StartByte()
	end := n.EndByte()

	start += c.offset

	if c.splitChunks && end-start > c.size {
		// TODO: It would be better to search for line endings maybe? Or other
		// context sensitive break point.
		end = start + c.size
		if end >= n.EndByte() {
			c.offset = 0
			return c.sourceCode[start:n.EndByte()], c.c.GoToNextSibling()
		}

		c.offset += c.size
		return c.sourceCode[start:end], true
	}

	c.offset = 0

	more := false
	for c.c.GoToNextSibling() {
		c.offset = 0
		if c.c.CurrentNode().EndByte()-start > c.size {
			more = true
			break
		}

		end = c.c.CurrentNode().EndByte()
	}
	return c.sourceCode[start:end], more
}
