package builder

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikkehtine/maison/options"
	"github.com/nikkehtine/maison/parser"
)

type Builder struct {
	Input    string
	Output   string
	Config   options.Config
	Contents []os.DirEntry
}

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
		b.Contents = listDir
		return nil
	} else {
		b.Input = filepath.Dir(b.Input)
		b.Contents = append(b.Contents, fs.FileInfoToDirEntry(input))
		return nil
	}
}

func (b *Builder) Build() error {
	for _, entry := range b.Contents {
		if entry.IsDir() {
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
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		log.Printf("building %s", entry.Name())

		input, err := os.ReadFile(filepath.Join(b.Input, entry.Name()))
		if err != nil {
			log.Fatal(err)
		}

		output, err := parser.Parse(input)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = os.Stat(b.Output); os.IsNotExist(err) {
			err := os.MkdirAll(b.Output, 0755)
			if err != nil {
				return err
			}
		}

		outFile := filepath.Join(b.Output, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".html")

		err = os.WriteFile(outFile, output, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
