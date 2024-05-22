/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package builder

import (
	"bytes"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
	),
	goldmark.WithExtensions(
		extension.GFM,
		emoji.Emoji,
		extension.DefinitionList,
		extension.Footnote,
		extension.Typographer,
	),
)

func (b *Builder) Parse(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := md.Convert(input, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
