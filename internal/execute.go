package internal

import (
	"io/fs"
	"path/filepath"

	"github.com/iamolegga/rebus/internal/generator"
	"github.com/iamolegga/rebus/internal/parser"
)

// Opts is a struct to pass options to Execute function.
type Opts struct {
	//Root is a folder to search recursively for rebus comment tags.
	Root string
	//WithContext is a flag to add context.Context to function signatures.
	WithContext bool
}

// Execute is an entrypoint for generating code by this package.
func Execute(opts Opts) error {
	gen := generator.New(opts.WithContext)

	if err := filepath.Walk(opts.Root, func(p string, f fs.FileInfo, err error) error {
		if f.IsDir() {
			return parser.NewParser(p, gen).Parse()
		}
		return nil
	}); err != nil {
		return err
	}

	return gen.Generate()
}
