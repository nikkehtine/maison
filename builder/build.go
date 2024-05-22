/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package builder

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/nikkehtine/maison/lib"
	"github.com/nikkehtine/maison/options"
)

type Builder struct {
	Input       string
	Output      string
	Config      options.Config
	Directories []os.DirEntry
	Documents   []os.DirEntry
	Files       []os.DirEntry
}

// Initialize builder object
func (b *Builder) Init(conf options.Config) error {
	if b.Input == "" {
		b.Input = conf.Input
	}
	if b.Output == "" {
		b.Output = conf.Output
	}
	b.Config = conf

	input, err := os.Stat(b.Input)
	if err != nil {
		return err
	}
	if input.IsDir() {
		listDir, err := os.ReadDir(b.Input)
		if err != nil {
			log.Fatal(err)
		}

		b.Documents = lib.Filter(listDir, func(e os.DirEntry) bool {
			return (!lib.IsHidden(e) &&
				!e.IsDir() &&
				filepath.Ext(e.Name()) == ".md")
		})

		b.Directories = lib.Filter(listDir, func(e os.DirEntry) bool {
			return (!lib.IsHidden(e) && e.IsDir())
		})

		b.Files = lib.Filter(listDir, func(e os.DirEntry) bool {
			return (!lib.IsHidden(e) &&
				!e.IsDir() &&
				filepath.Ext(e.Name()) != ".md")
		})
		return nil
	} else {
		b.Input = filepath.Dir(b.Input)
		b.Files = append(b.Files, fs.FileInfoToDirEntry(input))
		return nil
	}
}

// Build the site, copying all files and building documents
func (b *Builder) Build() error {
	if err := os.MkdirAll(b.Output, 0755); err != nil {
		return err
	}

	if b.Input != "." && b.Input != "./" {
		blue := color.New(color.FgBlue).SprintFunc()
		log.Printf("building %s", blue(b.Input))
	}

	yellowBg := color.New(color.BgYellow).SprintFunc()
	cyanBg := color.New(color.BgCyan).SprintFunc()

	// Build documents
	for _, entry := range b.Documents {
		inFileName := filepath.Join(b.Input, entry.Name())
		log.Printf("%s %s", yellowBg(" BUILD "), entry.Name())

		input, err := os.ReadFile(inFileName)
		if err != nil {
			lib.LogError(err)
			continue
		}

		output, err := b.Parse(input)
		if err != nil {
			lib.LogError(err)
			continue
		}

		outFileName := filepath.Join(b.Output,
			strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".html")

		err = os.WriteFile(outFileName, output, 0644)
		if err != nil {
			lib.LogError(err)
			continue
		}
	}

	// Copy files
	for _, entry := range b.Files {
		inFileName := filepath.Join(b.Input, entry.Name())
		log.Printf("%s %s", cyanBg(" COPY  "), entry.Name())

		in, err := os.ReadFile(inFileName)
		if err != nil {
			lib.LogError(err)
			continue
		}

		if err = os.WriteFile(filepath.Join(b.Output, entry.Name()), in, 0644); err != nil {
			lib.LogError(err)
			continue
		}
	}

	// Build directories
	for _, entry := range b.Directories {
		dirBuilder := Builder{
			Input:  filepath.Join(b.Input, entry.Name()),
			Output: filepath.Join(b.Output, entry.Name()),
			Config: b.Config,
		}
		dirBuilder.Init(b.Config)
		err := dirBuilder.Build()
		if err != nil {
			return err
		}
	}
	return nil
}
