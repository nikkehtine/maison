package builder

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

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

// Check if a file or directory is hidden
func isHidden(e os.DirEntry) bool {
	return strings.HasPrefix(e.Name(), ".") ||
		strings.HasPrefix(e.Name(), "_")
}

// Log error and move on to the next entry. I don't know if you can continue a loop from within here so just PLEASE use 'continue' right after it in the error check!!!
func logError(err error) {
	redBg := color.New(color.BgRed).SprintFunc()
	if err != nil {
		log.Printf("%s %s", redBg(" ERROR "), err)
	}
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
	blue := color.New(color.FgBlue).SprintFunc()
	log.Printf("created %s, building", blue(b.Input))

	yellowBg := color.New(color.BgYellow).SprintFunc()
	cyanBg := color.New(color.BgCyan).SprintFunc()

	// Build documents
	for _, entry := range b.Documents {
		inFileName := filepath.Join(b.Input, entry.Name())
		log.Printf("%s %s", yellowBg(" BUILD "), entry.Name())

		input, err := os.ReadFile(inFileName)
		if err != nil {
			logError(err)
			continue
		}

		output, err := b.Parse(input)
		if err != nil {
			logError(err)
			continue
		}

		outFileName := filepath.Join(b.Output,
			strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".html")

		err = os.WriteFile(outFileName, output, 0644)
		if err != nil {
			logError(err)
			continue
		}
	}

	// Copy files
	for _, entry := range b.Files {
		inFileName := filepath.Join(b.Input, entry.Name())
		log.Printf("%s %s", cyanBg(" COPY  "), entry.Name())

		in, err := os.ReadFile(inFileName)
		if err != nil {
			logError(err)
			continue
		}

		if err = os.WriteFile(filepath.Join(b.Output, entry.Name()), in, 0644); err != nil {
			logError(err)
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
