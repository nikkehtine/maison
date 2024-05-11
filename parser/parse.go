package parser

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func Parse(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := goldmark.Convert(input, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
