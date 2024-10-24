/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package builder

import (
	"bytes"
	"os"
	"text/template"
)

type PageRenderer struct {
	Title    string
	Body     string
	Template string
}

// Render outputs the rendered Markdown into a HTML template
func (r *PageRenderer) Output() ([]byte, error) {
	templatePath := "layout/main.html"

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(templatePath)
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
