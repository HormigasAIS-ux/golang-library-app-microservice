package repository

import (
	"auth_service/domain/model"
	mocks "auth_service/mocks/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success_create_user",
			args: args{
				&model.User{
					UUID:     "test-uuid",
					Username: "test",
					Fullname: "test",
					Email:    "test",
					Password: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "failed_create_user",
			args: args{
				&model.User{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.IUserRepo{}
			if !tt.wantErr {
				repo.On("Create", tt.args.user).Return(nil)
			} else {
				repo.On("Create", tt.args.user).Return(errors.New("failed to create user"))
			}
			err := repo.Create(tt.args.user)
			if tt.wantErr != (err != nil) {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserRepo_GetByUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "success_getUserByUUID",
			args: args{
				uuid: "test-uuid",
			},
			want: &model.User{
				UUID:     "test-uuid",
				Username: "test",
				Fullname: "test",
				Email:    "test",
				Password: "test",
			},
			wantErr: false,
		},
		{
			name: "failed_getUserByUUID_notFound",
			args: args{
				uuid: "23123212",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.IUserRepo{}

			if !tt.wantErr {
				repo.On("GetByUUID", tt.args.uuid).Return(tt.want, nil)
			} else {
				repo.On("GetByUUID", tt.args.uuid).Return(tt.want, errors.New("user not found"))
			}

			got, err := repo.GetByUUID(tt.args.uuid)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got, "Expected nil but got a user")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got, "Returned user does not match the expected user")
			}

			repo.AssertExpectations(t)
		})
	}
}
