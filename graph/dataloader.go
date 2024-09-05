package graph

import (
	"context"
	"errors"
	"mygql/graph/services"

	"github.com/graph-gophers/dataloader"
)

type Loaders struct {
	UserLoader *dataloader.Loader
}

func NewLoaders(Srv services.Services) *Loaders {
	userBatcher := &userBatcher{Srv: Srv}
	return &Loaders{
		UserLoader: dataloader.NewBatchedLoader(userBatcher.BatchGetUsers),
	}
}

type userBatcher struct {
	Srv services.Services
}

func (u *userBatcher) BatchGetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	IDs := make([]string, len(keys))
	for i, key := range keys {
		IDs[i] = key.String()
	}

	results := make([]*dataloader.Result, len(IDs))
	for i := range results {
		results[i] = &dataloader.Result{
			Error: errors.New("not found"),
		}
	}

	indexs := make(map[string]int, len(IDs))
	for i, ID := range IDs {
		indexs[ID] = i
	}

	users, err := u.Srv.ListUsersByID(ctx, IDs)

	for _, user := range users {
		var rsl *dataloader.Result
		if err != nil {
			rsl = &dataloader.Result{
				Error: err,
			}
		} else {
			rsl = &dataloader.Result{
				Data: user,
			}
		}
		results[indexs[user.ID]] = rsl
	}
	return results
}
