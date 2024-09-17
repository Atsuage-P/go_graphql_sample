package services_test

import (
	"context"
	"mygql/graph/model"
	"mygql/graph/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func Test_userService_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	srv := services.New(db)
	ctx := context.Background()
	mockSetup := func(mock sqlmock.Sqlmock, id, name string) {
		columns := []string{"id", "name"}
		mock.ExpectQuery(".*").WithArgs(id).WillReturnRows(
			sqlmock.NewRows(columns).AddRow(id, name),
		)
	}

	tests := []struct {
		title   string
		id      string
		name    string
		want    *model.User
		wantErr bool
	}{
		{
			title: "case1",
			id:    "U_ABC",
			name:  "abc",
			want:  &model.User{ID: "U_ABC", Name: "abc"},
		},
		{
			title: "case2",
			id:    "U_DEF",
			name:  "def",
			want:  &model.User{ID: "U_DEF", Name: "def"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSetup(mock, tt.id, tt.name)

			got, err := srv.GetUserByID(ctx, tt.id)
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUserByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
