package main

import (
	"log"
	"os"

	"github.com/iamolegga/rebus/internal"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("rebus should have single argument")
	}
	if err := internal.Execute(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
