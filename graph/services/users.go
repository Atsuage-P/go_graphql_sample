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

func convertUser(user *db.User) *model.User {
	return &model.User{
		ID:   user.ID,
		Name: user.Name,
	}
}
