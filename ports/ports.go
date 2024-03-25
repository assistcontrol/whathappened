package ports

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

const indexGlob = "/usr/ports/INDEX-1[0-9]*"

// Local gets a list of ports from the current live system, whether installed or
// not.
func Local() ([]string, error) {
	pkgquery, err := exec.Command("pkg", "query", "%o").Output()
	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, pkg := range strings.Split(string(pkgquery), "\n") {
		if len(pkg) == 0 {
			continue
		}

		if stat, err := os.Stat("/data/freebsd/ports/" + pkg); err != nil || !stat.IsDir() {
			continue
		}

		list = append(list, pkg)
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
