package tags

import (
	"fmt"
	"go/doc"
	"strings"
)

const (
	TagPrefix = "+rebus:"
	TagOut    = "out"
	TagPkg    = "pkg"
)

type Package struct {
	data map[string]Struct
}

type Struct struct {
	Out string
	Pkg string
}

func New(pkg *doc.Package) (*Package, error) {
	data := make(map[string]Struct)
	for _, t := range pkg.Types {
		if t.Doc == "" {
			continue
		}
		comments := strings.Split(t.Doc, "\n")
		for _, commentLine := range comments {
			commentLine = strings.TrimSpace(commentLine)
			if strings.HasPrefix(commentLine, TagPrefix) {
				commentLine = strings.TrimPrefix(commentLine, TagPrefix)
			} else {
				continue
			}
			pair := strings.Split(commentLine, "=")
			if len(pair) != 2 {
				return nil, fmt.Errorf("wrong comment annotation: %s", t.Doc)
			}
			s, ok := data[t.Name]
			if !ok {
				s = Struct{}
			}
			switch v := pair[1]; pair[0] {
			case TagOut:
				s.Out = v
			case TagPkg:
				s.Pkg = v
			default:
				return nil, fmt.Errorf("unknown tag: `%s%s`", TagPrefix, pair[0])
			}
			data[t.Name] = s
		}
	}
	return &Package{data}, nil
}

func (t *Package) For(structName string) (structTags Struct, ok bool) {
	structTags, ok = t.data[structName]
	return
}
