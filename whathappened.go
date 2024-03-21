package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/assistcontrol/whathappened/date"
	"github.com/assistcontrol/whathappened/git"
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
	list, err := git.Query(repoRoot, Date, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", list)
}
