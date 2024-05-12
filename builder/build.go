package builder

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

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

// Helper function
func filter[T any](slice []T, test func(T) bool) []T {
	var ret = make([]T, 0)
	for _, v := range slice {
		if test(v) {
			ret = append(ret, v)
		}
	}
	return ret
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

		isHidden := func(e os.DirEntry) bool {
			return strings.HasPrefix(e.Name(), ".") ||
				strings.HasPrefix(e.Name(), "_")
		}

		b.Documents = filter(listDir, func(e os.DirEntry) bool {
			return (!isHidden(e) &&
				!e.IsDir() &&
				filepath.Ext(e.Name()) == ".md")
		})

		b.Directories = filter(listDir, func(e os.DirEntry) bool {
			return (!isHidden(e) && e.IsDir())
		})

		b.Files = filter(listDir, func(e os.DirEntry) bool {
			return (!isHidden(e) &&
				!e.IsDir() &&
				filepath.Ext(e.Name()) == ".md")
		})
		return nil
	} else {
		b.Input = filepath.Dir(b.Input)
		b.Files = append(b.Files, fs.FileInfoToDirEntry(input))
		return nil
	}
}

func (b *Builder) Build() error {
	if err := os.MkdirAll(b.Output, 0755); err != nil {
		return err
	}

	// Build documents
	for _, entry := range b.Documents {
		inFileName := filepath.Join(b.Input, entry.Name())
		log.Printf("BUILD %s", entry.Name())

		input, err := os.ReadFile(inFileName)
		if err != nil {
			log.Fatal(err)
		}

		output, err := b.Parse(input)
		if err != nil {
			log.Fatal(err)
		}

		outFileName := filepath.Join(
			b.Output,
			strings.TrimSuffix(entry.Name(),
				filepath.Ext(entry.Name()),
			)+".html",
		)

		err = os.WriteFile(outFileName, output, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Build directories
	for _, entry := range b.Directories {
		log.Printf("MOVE DIR %s",
			filepath.Join(b.Input, entry.Name()))

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

		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}
	}
	return nil
}
