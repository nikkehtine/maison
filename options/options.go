/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package options

import (
	"os"

	"github.com/nikkehtine/maison/lib"
)

type Config struct {
	Input        string
	Output       string
	TemplatePath string
	IgnoreList   []string
}

func (c *Config) IsIgnored(e os.DirEntry) bool {
	return lib.Includes(e.Name(), c.IgnoreList)
}

var DefaultConfig = Config{
	Input:        ".",
	Output:       "./public",
	TemplatePath: "layout/main.html",
	IgnoreList:   []string{"public", "layout", "maison.config.toml", "maison.config.json"},
}
