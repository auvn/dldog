package gitrepo

import (
	"log"
	"reflect"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"gopkg.in/validator.v2"
)

type Desc struct {
	Url    string `validate:"nonzero"`
	Tag    string `validate:"semver"`
	Branch string
	Sha    string `validate:"sha"`
}

func (c Desc) Version() (string, bool) {
	version := firstNonEmpty(c.Branch, c.Tag, c.Sha)
	return version, version != ""
}

func Clone(dir string, desc Desc, opts ...Option) error {
	cloneOpts := git.CloneOptions{
		URL:  desc.Url,
		Tags: git.NoTags,
		//	Progress: os.Stdout,
	}

	for _, opt := range opts {
		opt(&cloneOpts)
	}

	checkout, err := prepareClone(desc, &cloneOpts)
	if err != nil {
		return errors.WithMessage(err, "prepareClone")
	}

	log.Printf("Cloning %q into %q...\n", desc.Url, dir)

	repo, err := git.PlainClone(dir, false, &cloneOpts)
	if err != nil {
		return errors.WithMessage(err, "git.Clone")
	}

	if checkout != nil {
		w, err := repo.Worktree()
		if err != nil {
			return errors.WithMessage(err, "repo.Worktree")
		}

		if err := w.Checkout(checkout); err != nil {
			return errors.WithMessage(err, "w.Checkout")
		}
	}
	return nil
}

func prepareClone(c Desc, cloneOpts *git.CloneOptions) (*git.CheckoutOptions, error) {
	if c.Sha != "" {
		return &git.CheckoutOptions{
			Hash: plumbing.NewHash(c.Sha),
		}, nil
	}

	if c.Branch != "" {
		cloneOpts.Depth = 1
		cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(c.Branch)
		cloneOpts.SingleBranch = true
		return nil, nil
	}

	if c.Tag != "" {
		cloneOpts.Depth = 1
		cloneOpts.ReferenceName = plumbing.NewTagReferenceName(c.Tag)
		cloneOpts.SingleBranch = true

		return nil, nil
	}

	return nil, errors.New("either Branch, Tag or Sha should be specified")
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func init() {
	validSha := regexp.MustCompile("[a-f0-9]{40}")

	//nolint:errcheck
	_ = validator.SetValidationFunc("semver",
		func(v interface{}, param string) error {
			st := reflect.ValueOf(v)
			if st.Kind() != reflect.String {
				return validator.ErrUnsupported
			}
			if s := st.String(); s != "" && !semver.IsValid(s) {
				return errors.New("invalid semver")
			}
			return nil
		})

	//nolint:errcheck
	_ = validator.SetValidationFunc("sha",
		func(v interface{}, param string) error {
			st := reflect.ValueOf(v)
			if st.Kind() != reflect.String {
				return validator.ErrUnsupported
			}
			if s := st.String(); s != "" && !validSha.MatchString(s) {
				return errors.New("invalid sha")
			}
			return nil
		})
}
