package parser

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/iamolegga/rebus/internal/generator"
	"github.com/iamolegga/rebus/internal/parser/tags"
	"golang.org/x/tools/go/packages"
)

var relativePathPrefix = fmt.Sprintf(".%s", string(os.PathSeparator))

const (
	ResultStructNameSuffix = "Result"
)

type Parser struct {
	folder     string
	gen        *generator.Generator
	importPath string
}

func NewParser(folder string, gen *generator.Generator) *Parser {
	// required for proper parsing by `packages` pkg
	if !strings.HasPrefix(folder, relativePathPrefix) {
		folder = relativePathPrefix + folder
	}

	return &Parser{folder: folder, gen: gen}
}

func (p *Parser) Parse() error {
	if ok := p.parseImportPath(); !ok {
		return nil
	}

	fileSet := token.NewFileSet()
	dir, err := parser.ParseDir(fileSet, p.folder, skipOwnGeneratedFilesFilter, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("unable to parse source files %s: %w", p.folder, err)
	}

	for _, pkgAST := range dir {
		pkg := doc.New(pkgAST, p.folder, doc.AllDecls)

		tagsParsed, err := tags.New(pkg)
		if err != nil {
			return fmt.Errorf("unable to parse tags: %w", err)
		}

		outputByInput := make(map[*ast.TypeSpec]string)
		packageByInput := make(map[*ast.TypeSpec]string)
		folderByInput := make(map[*ast.TypeSpec]string)

		// O(2n) but KISS
		walkTypeStructs(pkgAST, func(x *ast.TypeSpec, f *ast.File) {
			if _, ok := tagsParsed.For(x.Name.String()); ok {
				outputByInput[x] = ""
				packageByInput[x] = pkgAST.Name
				folderByInput[x], _ = path.Split(fileSet.Position(f.Pos()).Filename)
			}
		})
		walkTypeStructs(pkgAST, func(x *ast.TypeSpec, _ *ast.File) {
			name := x.Name.String()
			if !strings.HasSuffix(name, ResultStructNameSuffix) {
				return
			}
			in := strings.TrimSuffix(name, ResultStructNameSuffix)
			for spec := range outputByInput {
				if spec.Name.String() == in {
					outputByInput[spec] = name
				}
			}
		})

		for input, output := range outputByInput {
			pkg := packageByInput[input]
			dir := folderByInput[input]

			inputName := input.Name.String()
			inputTags, _ := tagsParsed.For(inputName)

			genDest := inputTags.Out
			if genDest == "" {
				return fmt.Errorf("comment tag `%s%s` should be presented for %s", tags.TagPrefix, tags.TagOut, inputName)
			}
			genDest = path.Join(dir, genDest)

			genPkg := inputTags.Pkg
			if genPkg == "" {
				_, genPkg = path.Split(genDest)
			}

			handler := &generator.Handler{
				Package: pkg,
				Input:   inputName,
				Output:  output,
			}

			p.gen.Add(genDest, genPkg, p.importPath, handler)
		}
	}
	return nil
}

func (p *Parser) parseImportPath() (ok bool) {
	pkgs, _ := packages.Load(&packages.Config{Mode: packages.NeedName}, p.folder)
	if len(pkgs) == 1 {
		p.importPath = pkgs[0].PkgPath
		ok = true
	}
	return
}

func walkTypeStructs(pkgAST *ast.Package, fn func(*ast.TypeSpec, *ast.File)) {
	for _, file := range pkgAST.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.TypeSpec:
				switch x.Type.(type) {
				case *ast.StructType:
					fn(x, file)
				}
			}
			return true
		})
	}
}

func skipOwnGeneratedFilesFilter(info fs.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), generator.GeneratedCodeFileName)
}
