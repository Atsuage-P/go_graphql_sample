package services

import (
	"context"
	"mygql/graph/model"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Services interface {
	UserService
	RepoService
	IssueService
	PullRequestService
}

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	ListUsersByID(ctx context.Context, IDs []string) ([]*model.User, error)
}

type RepoService interface {
	GetRepoByID(ctx context.Context, id string) (*model.Repository, error)
	GetRepoByFullName(ctx context.Context, owner, name string) (*model.Repository, error)
}

type IssueService interface {
	GetIssueByID(ctx context.Context, id string) (*model.Issue, error)
	GetIssueByRepoAndNumber(ctx context.Context, repoID string, number int) (*model.Issue, error)
	ListIssueInRepository(ctx context.Context, repoID string, after *string, before *string, first *int, last *int) (*model.IssueConnection, error)
}

type PullRequestService interface {
	GetPullRequestByID(ctx context.Context, id string) (*model.PullRequest, error)
	GetPullRequestByRepoAndNumber(ctx context.Context, repoID string, number int) (*model.PullRequest, error)
	ListPullRequestInRepository(ctx context.Context, repoID string, after *string, before *string, first *int, last *int) (*model.PullRequestConnection, error)
}

type services struct {
	*userService
	*repoService
	*issueService
	*pullRequestService
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		userService:        &userService{exec: exec},
		repoService:        &repoService{exec: exec},
		issueService:       &issueService{exec: exec},
		pullRequestService: &pullRequestService{exec: exec},
	}
}
