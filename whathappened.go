package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/assistcontrol/whathappened/date"
	"github.com/assistcontrol/whathappened/repo"
)

const (
	repoRoot = "/Users/adamw/build/dotfiles"
)

var Date string

func init() {
	log.SetFlags(0) // Just the facts, ma'am

	flag.StringVar(&Date, "date", date.Yesterday(), "Date (YYYY-MM-DD), default: yesterday")
	flag.Parse()
}

func main() {
	r, err := repo.New("/Users/adamw/build/dotfiles", Date)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.Query(nil); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", r.Rev())
}
