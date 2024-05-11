package builder

import (
	"fmt"
	"log"
	"os"
)

type Builder struct {
	Dir      string
	OutDir   string
	Contents []os.DirEntry
}

func (b *Builder) Init(path string) error {
	contents, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	b.Contents = contents
	return nil
}

func (b *Builder) Build() error {
	for _, entry := range b.Contents {
		if entry.IsDir() {
			fmt.Print("* ")
		} else {
			fmt.Print("  ")
		}
		fmt.Println(entry.Name())
	}
	return nil
}
