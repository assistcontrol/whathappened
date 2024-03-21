package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/assistcontrol/whathappened/date"
)

const (
	revisionFormat = "%ct %H"
	dateFormat     = "%a %d %b %H:%M"
)

// Query returns a list of revisions that match a set of
// queries.
func Query(repo, day string, limiters []string) ([]string, error) {
	dateRange, err := date.Range(day)
	if err != nil {
		return nil, err
	}

	args := []string{}
	args = append(args, "-C", repo)
	args = append(args, "log")
	args = append(args, dateRange...)
	args = append(args, fmt.Sprintf("--format='%s'", revisionFormat))
	args = append(args, limiters...)

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), "\n"), nil
}
