package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"fmt"
	"mygql/graph/model"
	"mygql/internal"
	"strings"

	"github.com/graph-gophers/dataloader"
)

// Author is the resolver for the author field.
func (r *issueResolver) Author(ctx context.Context, obj *model.Issue) (*model.User, error) {
	thunk := r.Loaders.UserLoader.Load(ctx, dataloader.StringKey(obj.Author.ID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	user, ok := result.(*model.User)
	if !ok {
		return nil, errors.New("failed to cast to *model.User")
	}
	return user, nil
}

// Repository is the resolver for the repository field.
func (r *issueResolver) Repository(ctx context.Context, obj *model.Issue) (*model.Repository, error) {
	panic(fmt.Errorf("not implemented: Repository - repository"))
}

// ProjectItems is the resolver for the projectItems field.
func (r *issueResolver) ProjectItems(ctx context.Context, obj *model.Issue, after *string, before *string, first *int, last *int) (*model.ProjectV2ItemConnection, error) {
	panic(fmt.Errorf("not implemented: ProjectItems - projectItems"))
}

// AddProjectV2ItemByID is the resolver for the addProjectV2ItemById field.
func (r *mutationResolver) AddProjectV2ItemByID(ctx context.Context, input model.AddProjectV2ItemByIDInput) (*model.AddProjectV2ItemByIDPayload, error) {
	panic(fmt.Errorf("not implemented: AddProjectV2ItemByID - addProjectV2ItemById"))
}

// Items is the resolver for the items field.
func (r *projectV2Resolver) Items(ctx context.Context, obj *model.ProjectV2, after *string, before *string, first *int, last *int) (*model.ProjectV2ItemConnection, error) {
	panic(fmt.Errorf("not implemented: Items - items"))
}

// Owner is the resolver for the owner field.
func (r *projectV2Resolver) Owner(ctx context.Context, obj *model.ProjectV2) (*model.User, error) {
	return r.Srv.GetUserByID(ctx, obj.Owner.ID)
}

// Content is the resolver for the content field.
func (r *projectV2ItemResolver) Content(ctx context.Context, obj *model.ProjectV2Item) (model.ProjectV2ItemContent, error) {
	panic(fmt.Errorf("not implemented: Content - content"))
}

// Repository is the resolver for the repository field.
func (r *pullRequestResolver) Repository(ctx context.Context, obj *model.PullRequest) (*model.Repository, error) {
	return r.Srv.GetRepoByID(ctx, obj.Repository.ID)
}

// ProjectItems is the resolver for the projectItems field.
func (r *pullRequestResolver) ProjectItems(ctx context.Context, obj *model.PullRequest, after *string, before *string, first *int, last *int) (*model.ProjectV2ItemConnection, error) {
	panic(fmt.Errorf("not implemented: ProjectItems - projectItems"))
}

// Repository is the resolver for the repository field.
func (r *queryResolver) Repository(ctx context.Context, name string, owner string) (*model.Repository, error) {
	user, err := r.Srv.GetUserByName(ctx, owner)
	if err != nil {
		return nil, err
	}
	return r.Srv.GetRepoByFullName(ctx, user.ID, name)
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, name string) (*model.User, error) {
	return r.Srv.GetUserByName(ctx, name)
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	nElems := strings.SplitN(id, "_", 2)
	nType, _ := nElems[0], nElems[1]

	switch nType {
	case "U":
		return r.Srv.GetUserByID(ctx, id)
	case "REPO":
		return r.Srv.GetRepoByID(ctx, id)
	case "ISSUE":
		return r.Srv.GetIssueByID(ctx, id)
	// case "PJ":
	// 	return r.Srv.GetProjectByID(ctx, id)
	// case "PR":
	// 	return r.Srv.GetPullRequestByID(ctx, id)
	default:
		return nil, errors.New("invalid ID")
	}
}

// Owner is the resolver for the owner field.
func (r *repositoryResolver) Owner(ctx context.Context, obj *model.Repository) (*model.User, error) {
	return r.Srv.GetUserByID(ctx, obj.Owner.ID)
}

// Issue is the resolver for the issue field.
func (r *repositoryResolver) Issue(ctx context.Context, obj *model.Repository, number int) (*model.Issue, error) {
	return r.Srv.GetIssueByRepoAndNumber(ctx, obj.ID, number)
}

// Issues is the resolver for the issues field.
func (r *repositoryResolver) Issues(ctx context.Context, obj *model.Repository, after *string, before *string, first *int, last *int) (*model.IssueConnection, error) {
	return r.Srv.ListIssueInRepository(ctx, obj.ID, after, before, first, last)
}

// PullRequest is the resolver for the pullRequest field.
func (r *repositoryResolver) PullRequest(ctx context.Context, obj *model.Repository, number int) (*model.PullRequest, error) {
	return r.Srv.GetPullRequestByRepoAndNumber(ctx, obj.ID, number)
}

// PullRequests is the resolver for the pullRequests field.
func (r *repositoryResolver) PullRequests(ctx context.Context, obj *model.Repository, after *string, before *string, first *int, last *int) (*model.PullRequestConnection, error) {
	return r.Srv.ListPullRequestInRepository(ctx, obj.ID, after, before, first, last)
}

// ProjectV2 is the resolver for the projectV2 field.
func (r *userResolver) ProjectV2(ctx context.Context, obj *model.User, number int) (*model.ProjectV2, error) {
	panic(fmt.Errorf("not implemented: ProjectV2 - projectV2"))
}

// ProjectV2s is the resolver for the projectV2s field.
func (r *userResolver) ProjectV2s(ctx context.Context, obj *model.User, after *string, before *string, first *int, last *int) (*model.ProjectV2Connection, error) {
	panic(fmt.Errorf("not implemented: ProjectV2s - projectV2s"))
}

// Issue returns internal.IssueResolver implementation.
func (r *Resolver) Issue() internal.IssueResolver { return &issueResolver{r} }

// Mutation returns internal.MutationResolver implementation.
func (r *Resolver) Mutation() internal.MutationResolver { return &mutationResolver{r} }

// ProjectV2 returns internal.ProjectV2Resolver implementation.
func (r *Resolver) ProjectV2() internal.ProjectV2Resolver { return &projectV2Resolver{r} }

// ProjectV2Item returns internal.ProjectV2ItemResolver implementation.
func (r *Resolver) ProjectV2Item() internal.ProjectV2ItemResolver { return &projectV2ItemResolver{r} }

// PullRequest returns internal.PullRequestResolver implementation.
func (r *Resolver) PullRequest() internal.PullRequestResolver { return &pullRequestResolver{r} }

// Query returns internal.QueryResolver implementation.
func (r *Resolver) Query() internal.QueryResolver { return &queryResolver{r} }

// Repository returns internal.RepositoryResolver implementation.
func (r *Resolver) Repository() internal.RepositoryResolver { return &repositoryResolver{r} }

// User returns internal.UserResolver implementation.
func (r *Resolver) User() internal.UserResolver { return &userResolver{r} }

type issueResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type projectV2Resolver struct{ *Resolver }
type projectV2ItemResolver struct{ *Resolver }
type pullRequestResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type repositoryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
