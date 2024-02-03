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

func FindSplitBounds(data []byte, minSize, maxSize int) int {
	lastWhitespace := -1
	lastNewline := -1
	lastDoubleNewline := -1
	for ndx := 0; ndx < len(data); ndx++ {
		if ndx > maxSize {
			if lastDoubleNewline > minSize {
				return lastDoubleNewline
			}

			if lastNewline > minSize {
				return lastNewline
			}

			if lastWhitespace > minSize {
				return lastWhitespace
			}

			// At this point we're past the max split but nothing structurally is giving us
			// a better split.  We'll just return what we have.
			return ndx
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

	return -1
}

func (p *Plaintext) splitFunc(data []byte, atEOF bool) (int, []byte, error) {
	n := FindSplitBounds(data, p.minSize, p.maxSize)
	if n < p.minSize {
		if atEOF {
			return len(data), data, bufio.ErrFinalToken
		}

		// Need more data
		return 0, nil, nil
	}

	return n, data[:n], nil
}

func (p *Plaintext) Next() ([]byte, bool) {
	if !p.s.Scan() {
		return nil, false
	}

	return p.s.Bytes(), true
}
