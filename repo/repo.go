package repo

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os/exec"

	"github.com/assistcontrol/whathappened/date"
)

const (
	revisionFormat = "%ct %H"
	dateFormat     = "%a %d %b %H:%M"
)

type Repo struct {
	Path      string
	Date      string
	dateRange []string
	commits   []string
	hashes    hashMap
	revisions []string
}

type hashMap map[[16]byte]bool

func New(path, day string) (*Repo, error) {
	d, err := date.Range(day)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Path:      path,
		Date:      day,
		dateRange: d,
		commits:   []string{},
		hashes:    make(hashMap),
		revisions: []string{},
	}, nil
}

func (r *Repo) Add(revision []byte) {
	if len(revision) == 0 {
		return
	}

	sum := md5.Sum(revision)

	if _, exists := r.hashes[sum]; !exists {
		r.hashes[sum] = true
		r.revisions = append(r.revisions, string(revision))
	}
}

// Query returns a list of revisions that match a set of
// queries.
func (r *Repo) Query(limiters []string) error {
	args := []string{}
	args = append(args, "-C", r.Path)
	args = append(args, "log")
	args = append(args, r.dateRange...)
	args = append(args, fmt.Sprintf("--format='%s'", revisionFormat))
	args = append(args, limiters...)

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return err
	}

	for _, rev := range bytes.Split(out, []byte("\n")) {
		r.Add(rev)
	}

	return nil
}
