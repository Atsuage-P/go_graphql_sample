package services

import (
	"context"
	"mygql/graph/db"
	"mygql/graph/model"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type userService struct {
	exec boil.ContextExecutor
}

func convertUser(user *db.User) *model.User {
	return &model.User{
		ID:   user.ID,
		Name: user.Name,
	}
}

func convertUserSlice(users []*db.User) []*model.User {
	result := make([]*model.User, 0, len(users))
	for _, user := range users {
		result = append(result, convertUser(user))
	}
	return result
}

func (u *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := db.FindUser(ctx, u.exec, id,
		db.UserTableColumns.ID,
		db.UserTableColumns.Name,
	)
	if err != nil {
		return nil, err
	}
	return convertUser(user), nil
}

func (u *userService) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	user, err := db.Users( // FROM
		qm.Select(db.UserTableColumns.ID, db.UserTableColumns.Name), // SELECT
		db.UserWhere.Name.EQ(name),                                  // WHERE
	).One(ctx, u.exec) // LIMIT

	if err != nil {
		return nil, err
	}
	return convertUser(user), nil
}

func (u *userService) ListUsersByID(ctx context.Context, IDs []string) ([]*model.User, error) {
	users, err := db.Users(
		qm.Select(db.UserTableColumns.ID, db.UserTableColumns.Name),
		db.UserWhere.ID.IN(IDs),
	).All(ctx, u.exec)
	if err != nil {
		return nil, err
	}

	return convertUserSlice(users), nil
}
