package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/assistcontrol/whathappened/date"
	"github.com/assistcontrol/whathappened/ports"
	"github.com/assistcontrol/whathappened/repo"
)

const (
	repoRoot = "/data/freebsd/"
)

var (
	Date    string
	OSVer   string
	Queries = map[string][][]string{
		"relevant": {
			{"--committer", "adamw@FreeBSD.org"},
			{"--grep", "adamw"},
		},
		"other": {
			{"Mk", "Tools", "Templates"}, // Dirs to always watch
		},
		"src": {
			{"stable/"}, // OSVer gets appended later...
		},
	}
)

func init() {
	log.SetFlags(0) // Just the facts, ma'am

	flag.StringVar(&Date, "date", date.Yesterday(), "Date (YYYY-MM-DD), default: yesterday")
	flag.Parse()

	repo.Date = Date
	repo.Base = repoRoot

	version, err := exec.Command("uname", "-U").Output()
	if err != nil {
		log.Fatal(err)
	}
	OSVer = string(version[0:2])
}

func main() {
	local, err := ports.Local()
	if err != nil {
		log.Fatal(err)
	}

	mine, err := ports.Mine()
	if err != nil {
		log.Fatal(err)
	}

	// RELEVANT
	Queries["relevant"] = append(Queries["relevant"], local, mine)
	relevant, err := repo.Commits(repo.Config{
		Repo:    "ports",
		Queries: Queries["relevant"],
		Format:  commitFmt("ports"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// OTHER
	other, err := repo.Commits(repo.Config{
		Repo:    "ports",
		Queries: Queries["other"],
		Format:  commitFmt("ports"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// SRC
	Queries["src"][0][0] += OSVer
	src, err := repo.Commits(repo.Config{
		Repo:    "src",
		Queries: Queries["src"],
		Format:  commitFmt("src"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output
	if len(relevant) != 0 {
		fmt.Print(title("relevant ports"))
		fmt.Println(relevant)
	}

	if len(other) != 0 {
		fmt.Print(title("other ports"))
		fmt.Println(other)
	}

	if len(src) != 0 {
		fmt.Print(title("src"))
		fmt.Println(src)
	}
}

// title returns a formatted title.
func title(s string) string {
	return fmt.Sprintf("███ %s\n\n", s)
}

// commitFmt returns a format string for a commit log.
func commitFmt(repo string) string {
	h := ""
	for range 20 {
		h += "━"
	}

	h += fmt.Sprintf("%%nCommitter: %%cl (%%cn)%%nDate: %%cd%%nCommit: https://cgit.freebsd.org/%s/commit/?id=%%h%%n%%n%%B", repo)

	return h
}
