package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/denormal/go-gitignore"
)

var CLI struct {
	Dir     string   `arg:"" type:"path" help:"Directory to walk through"`
	Ext     []string `short:"x" help:"File extensions to include"`
	Exclude []string `short:"e" help:"Exclude files or directories"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if err := run(); err != nil {
		ctx.FatalIfErrorf(err)
	}
}

func run() error {
	log.Printf("Walking through %s\n", CLI.Dir)

	var ignores []gitignore.GitIgnore
	ignored := func(path string) bool {
		if filepath.Base(path) == ".git" {
			return true
		}
		for _, ignore := range ignores {
			if ignore.Match(path) != nil {
				return true
			}
		}
		for _, pattern := range CLI.Exclude {
			match, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				panic(err)
			}
			if match {
				return true
			}
		}
		return false
	}

	// Parse .gitignore and .dockerignore files
	for _, fileName := range []string{".gitignore", ".dockerignore"} {
		ignore, err := gitignore.NewFromFile(filepath.Join(CLI.Dir, fileName))
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", fileName, err)
		}
		log.Printf("Parsed %s (base: %s)\n", fileName, ignore.Base())
		ignores = append(ignores, ignore)
	}

	// Walk through the directory and concatenate all files.
	var concatenation strings.Builder
	err := filepath.Walk(CLI.Dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Clean(path) == filepath.Clean(CLI.Dir) {
			return nil
		}
		if info.IsDir() && ignored(path) {
			return filepath.SkipDir
		}
		if info.IsDir() || ignored(path) {
			return nil
		}
		if len(CLI.Ext) > 0 {
			var match bool
			for _, ext := range CLI.Ext {
				if filepath.Ext(path) == ext {
					match = true
					break
				}
			}
			if !match {
				return nil
			}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(CLI.Dir, path)
		if err != nil {
			return err
		}
		concatenation.WriteString(fmt.Sprintf("# %s\n", relPath))
		concatenation.WriteString(fmt.Sprintf("```%s\n", filepath.Ext(path)))
		concatenation.WriteString(fmt.Sprintf("%s\n", data))
		concatenation.WriteString("```\n\n")
		return nil
	})
	fmt.Println(concatenation.String())

	return err
}
