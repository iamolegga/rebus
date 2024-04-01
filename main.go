package main

import (
	"log"

	"github.com/iamolegga/rebus/internal"
	"github.com/jessevdk/go-flags"
)

var cfg struct {
	Context bool `short:"c" long:"context" description:"Add context.Context to function signatures"`
}

var path string

func main() {
	ff, err := flags.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	switch len(ff) {
	case 0:
		path = "."
	case 1:
		path = ff[0]
	default:
		log.Fatal("rebus should have zero or one argument")
	}

	if err := internal.Execute(internal.Opts{
		Root:        path,
		WithContext: cfg.Context,
	}); err != nil {
		log.Fatal(err)
	}
}
