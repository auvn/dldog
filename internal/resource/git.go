package resource

import (
	"os"

	"github.com/auvn/dldog/internal/gitrepo"
)

var DefaultGit = &git{
	opts: []gitrepo.Option{
		gitrepo.WithAccessToken(os.Getenv("GITHUB_TOKEN")),
	},
}

type git struct {
	opts []gitrepo.Option
}

func (g *git) Clone(dir string, repo gitrepo.Desc, opts ...gitrepo.Option) error {
	return gitrepo.Clone(dir, repo, append(g.opts, opts...)...)
}
