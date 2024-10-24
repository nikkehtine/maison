/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package builder

import (
	"bytes"
	"html/template"
)

type PageRenderer struct {
	Title string
	Body  string
}

// Render outputs the rendered Markdown into a HTML template
func (r *PageRenderer) Output(tmpl *template.Template) ([]byte, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, r)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
