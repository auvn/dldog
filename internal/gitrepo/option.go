package gitrepo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Option func(*git.CloneOptions)

func WithAccessToken(t string) Option {
	return func(opt *git.CloneOptions) {
		if t == "" {
			return
		}

		opt.Auth = &http.BasicAuth{
			Username: "protomod",
			Password: t,
		}
	}
}
