package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/sivistrukov/vk-assigment/internal/models"
)

type userRepoMock struct {
	Users map[string]models.User
}

func initUsers() map[string]models.User {
	users := make(map[string]models.User, 1)
	password, _ := HashPassword("password")

	users["user"] = models.User{
		Username: "user",
		Password: password,
	}

	return users
}

func initUsersRepo(users map[string]models.User) userRepoMock {
	var userRepo userRepoMock
	userRepo.Users = users

	return userRepo
}

func (r *userRepoMock) GetByUsername(_ context.Context, username string) (models.User, error) {
	user, ok := r.Users[username]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func TestService_Authenticate(t *testing.T) {
	users := initUsers()
	userRepo := initUsersRepo(users)

	type args struct {
		username string
		password string
		userRepo UserRepo
	}
	tests := []struct {
		name    string
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				username: "user",
				password: "password",
				userRepo: &userRepo,
			},
			want:    users["user"],
			wantErr: false,
		},
		{
			name: "wrong username",
			args: args{
				username: "user_not_exists",
				password: "password",
				userRepo: &userRepo,
			},
			want:    models.User{},
			wantErr: true,
		},
		{
			name: "wrong password",
			args: args{
				username: "user",
				password: "wrong_password",
				userRepo: &userRepo,
			},
			want:    models.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				userRepo: tt.args.userRepo,
			}

			got, err := s.Authenticate(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}
