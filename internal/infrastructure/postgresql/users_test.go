package postgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sivistrukov/vk-assigment/internal/models"
)

func TestUserRepo_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	type args struct {
		context  context.Context
		username string
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		id           uint
	}{
		{
			name: "basic",
			args: args{
				context:  context.Background(),
				username: "user",
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT").
					WithArgs("user").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "is_admin"}).
						AddRow(1, "user", "password", true))
			},
			wantErr: false,
			id:      1,
		},
		{
			name: "user not exist",
			args: args{
				context:  context.Background(),
				username: "not_exist_user",
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery("SELECT").
					WithArgs("not_exist_user").
					WillReturnError(&ErrRecordNotFound{})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := repo.GetByUsername(tt.args.context, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID != tt.id {
				t.Errorf("UserRepo.GetByUsername() id = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestUserRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	type args struct {
		context context.Context
		user    *models.User
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "basic",
			args: args{
				context: context.Background(),
				user: &models.User{
					Username: "user",
					Password: "password",
					IsAdmin:  true,
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectPrepare("INSERT").ExpectQuery().
					WithArgs("user", "password", true).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			err := repo.Create(tt.args.context, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
