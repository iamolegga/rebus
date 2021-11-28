package internal

import (
	"io/fs"
	"path/filepath"

	"github.com/iamolegga/rebus/internal/generator"
	"github.com/iamolegga/rebus/internal/parser"
)

//Execute is an entrypoint for generating code by this package.
//root is the folder to search recursively for rebus comment tags
func Execute(root string) error {
	gen := generator.New()

	if err := filepath.Walk(root, func(p string, f fs.FileInfo, err error) error {
		if f.IsDir() {
			return parser.NewParser(p, gen).Parse()
		}
		return nil
	}); err != nil {
		return err
	}

	return gen.Generate()
}
