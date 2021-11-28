package generator

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"text/template"
)

//go:embed bus.go.tmpl
var tmpl string
var t = template.Must(template.New("bus").Parse(tmpl))

const GeneratedCodeFileName = "generated.go"

func New() *Generator {
	return &Generator{make(map[string]*File)}
}

type Generator struct {
	data map[string]*File
}

func (g *Generator) Generate() error {
	for pkgPath, payload := range g.data {
		file := path.Join(pkgPath, GeneratedCodeFileName)

		if err := os.MkdirAll(pkgPath, os.ModePerm); err != nil {
			return fmt.Errorf("unable to create directory %s: %w", pkgPath, err)
		}
		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("unable to create file %s: %w", file, err)
		}
		if err = t.Execute(f, payload); err != nil {
			return fmt.Errorf("unable to write to file %s: %w", file, err)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("unable to close file %s: %w", file, err)
		}
	}
	return nil
}

func (g *Generator) Add(filePath, pkg, imp string, handler *Handler) {
	existing, ok := g.data[filePath]
	if ok {
		existing.Handlers = append(existing.Handlers, handler)
		for _, v := range existing.Imports {
			if v == imp {
				return
			}
		}
		existing.Imports = append(existing.Imports, imp)
		g.data[filePath] = existing
	} else {
		g.data[filePath] = &File{
			Package:  pkg,
			Imports:  []string{imp},
			Handlers: []*Handler{handler},
		}
	}
}

type File struct {
	Package  string
	Imports  []string
	Handlers []*Handler
}

type Handler struct {
	Package string
	Input   string
	Output  string
}

type GenData = map[string]File
