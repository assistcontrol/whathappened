package repo

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os/exec"

	"github.com/assistcontrol/whathappened/date"
)

const (
	revisionFormat = "%H"
	dateFormat     = "%a %d %b %H:%M"
)

var Date, Base string

type Config struct {
	Repo    string
	Queries [][]string
	Format  string
}

type Repo struct {
	Path      string
	dateRange []string
	commits   []string
	hashes    hashMap
	revisions []string
}

type hashMap map[[16]byte]bool

func New(path string) (*Repo, error) {
	d, err := date.Range(Date)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Path:      path,
		dateRange: d,
		commits:   []string{},
		hashes:    make(hashMap),
		revisions: []string{},
	}, nil
}

func Commits(c Config) (string, error) {
	r, err := New(Base + c.Repo)
	if err != nil {
		return "", err
	}

	err = r.Update()
	if err != nil {
		return "", err
	}

	for _, q := range c.Queries {
		if err := r.Query(q); err != nil {
			return "", err
		}
	}

	logs, err := r.Logs(c.Format)
	if err != nil {
		return "", err
	}

	return logs, nil
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
	args = append(args, fmt.Sprintf("--format=%s", revisionFormat))
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

func (r *Repo) Logs(format string) (string, error) {
	args := []string{}
	args = append(args, "-C", r.Path)
	args = append(args, "show", "--no-patch")
	args = append(args, fmt.Sprintf("--date=format-local:%s", dateFormat))
	args = append(args, fmt.Sprintf("--format=%s", format))
	args = append(args, r.revisions...)

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Update runs git pull on the repository.
func (r *Repo) Update() error {
	args := []string{}
	args = append(args, "-C", r.Path)
	args = append(args, "pull", "-q")

	_, err := exec.Command("git", args...).Output()
	if err != nil {
		return err
	}

	return nil
}
