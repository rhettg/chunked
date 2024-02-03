package plaintext

import (
	"bufio"
	"io"
	"unicode"
)

func isByteSpace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

type Plaintext struct {
	s       *bufio.Scanner
	minSize int
	maxSize int
}

func New(r io.Reader, minSize, maxSize int) *Plaintext {
	p := Plaintext{
		minSize: minSize,
		maxSize: maxSize,
	}

	p.s = bufio.NewScanner(r)
	p.s.Split(p.splitFunc)
	return &p
}

func (p *Plaintext) splitFunc(data []byte, atEOF bool) (int, []byte, error) {
	lastWhitespace := -1
	lastNewline := -1
	lastDoubleNewline := -1
	for ndx := 0; ndx < len(data); ndx++ {
		if ndx > p.maxSize {
			if lastDoubleNewline > p.minSize {
				return lastDoubleNewline, data[:lastDoubleNewline], nil
			}

			if lastNewline > p.minSize {
				return lastNewline, data[:lastNewline], nil
			}

			if lastWhitespace > p.minSize {
				return lastWhitespace, data[:lastWhitespace], nil
			}

			// At this point we're past the max split but nothing structurally is giving us
			// a better split.  We'll just return what we have.
			return ndx, data[:ndx], nil
		}

		if isByteSpace(data[ndx]) {
			lastWhitespace = ndx
		}

		if data[ndx] == '\n' {
			lastNewline = ndx
		}

		if ndx > 0 && data[ndx] == '\n' && data[ndx-1] == '\n' {
			lastDoubleNewline = ndx
		}
	}

	// At this point we've been through all the data. We must need more data to
	// find an appropriate sized split.

	if atEOF {
		// If we're at EOF, we have a final split, no matter how big it is
		return len(data), data, bufio.ErrFinalToken
	}

	// Need more data
	return 0, nil, nil
}

func (p *Plaintext) Next() ([]byte, bool) {
	if !p.s.Scan() {
		return nil, false
	}

	return p.s.Bytes(), true
}
