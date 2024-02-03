package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rhettg/chunker/plaintext"
	"github.com/rhettg/chunker/sitter"
	treesitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

var (
	chunkerType = "sitter"
	minSize     = 256
	maxSize     = 0
)

func init() {
	flag.StringVar(&chunkerType, "type", chunkerType, "chunker type")
	flag.IntVar(&minSize, "min", minSize, "minimum chunk size")
	flag.IntVar(&maxSize, "max", maxSize, "maximum chunk size")
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

type nexter interface {
	Next() ([]byte, bool)
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Usage: go run main.go <file>")
	}
	fName := flag.Args()[0]

	f, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}

	var nr nexter
	switch chunkerType {
	case "sitter":
		parser := treesitter.NewParser()
		parser.SetLanguage(golang.GetLanguage())
		sourceCode, _ := io.ReadAll(f)

		nr, err = sitter.New(golang.GetLanguage(), sourceCode, minSize, maxSize)
		if err != nil {
			log.Fatal(err)
		}
	case "plaintext":
		nr = plaintext.New(f, minSize, maxSize)
	default:
		errorf("Unknown chunker type: %s", chunkerType)
	}

	more := true
	var n []byte
	for more {
		n, more = nr.Next()
		fmt.Printf("Chunk (%d):\n", len(n))
		fmt.Println(string(n))
	}
}
