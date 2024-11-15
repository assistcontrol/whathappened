package ports

import (
	"bufio"
	"database/sql"
	"os"
	"path/filepath"
	"slices"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	indexGlob = "/usr/ports/INDEX-1[0-9]*"
	repoBase  = "/data/freebsd/ports/"
	repoDB    = "/var/db/pkg/local.sqlite"
)

// Local gets a list of ports from the current live system, whether installed or
// not.
func Local() ([]string, error) {
	db, err := sql.Open("sqlite3", repoDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT origin FROM packages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []string{}
	for rows.Next() {
		var origin string
		if err := rows.Scan(&origin); err != nil {
			return nil, err
		}

		// Check if the port exists in the ports tree.
		if stat, err := os.Stat(repoBase + origin); err != nil || !stat.IsDir() {
			continue
		}

		list = append(list, origin)
	}

	return list, nil
}

// Mine returns a list of ports that I maintain.
func Mine() ([]string, error) {
	indexMatches, err := filepath.Glob(indexGlob)
	if err != nil {
		return nil, err
	}
	slices.Sort(indexMatches)
	slices.Reverse(indexMatches)
	index := indexMatches[0]

	file, err := os.Open(index)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	matches := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if parts[5] == "adamw@FreeBSD.org" {
			dir := strings.TrimPrefix(parts[1], "/usr/ports/")
			matches = append(matches, dir)
		}
	}

	return matches, nil
}
