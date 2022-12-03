package repo

import (
	"context"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/sirupsen/logrus"
)

// Service of repository
type Service struct {
	config     *config.Config
	scmService core.SCMService
	userStore  core.UserStore
	repoStore  core.RepoStore
}

// NewService of repository
func NewService(
	config *config.Config,
	scmService core.SCMService,
	userStore core.UserStore,
	repoStore core.RepoStore,
) *Service {
	return &Service{
		config:     config,
		scmService: scmService,
		userStore:  userStore,
		repoStore:  repoStore,
	}
}

// Synchronize repository with remote and store to database
func (s *Service) Synchronize(ctx context.Context, user *core.User) error {
	userRepos := make([]*core.Repo, 0)
	provider := s.config.Provider()
	client, err := s.scmService.Client(provider)
	if err != nil {
		return err
	}
	repos, err := client.Repositories().List(ctx, user)
	if err != nil {
		logrus.Warnln(err)
		return err
	}
	if err := s.repoStore.BatchUpdateOrCreate(repos); err != nil {
		return err
	}
	userRepos = append(userRepos, repos...)

	return s.userStore.UpdateRepositories(user, userRepos)
}
