/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package builder

import (
	"bytes"
	"text/template"
)

type PageRenderer struct {
	Title    string
	Body     string
	Template string
}

// Render outputs the rendered Markdown into a HTML template
func (r *PageRenderer) Output() ([]byte, error) {
	tmpl, err := template.ParseFiles("layout/main.html")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, r)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
