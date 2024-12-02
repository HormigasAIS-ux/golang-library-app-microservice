package repository

import (
	"auth_service/domain/model"
	mocks "auth_service/mocks/repository"
	"errors"
	"testing"
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

// func TestUserRepo_GetByUUID(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *model.User
// 		wantErr bool
// 	}{
// 		{
// 			name: "success_get_user_by_uuid",
// 			args: args{
// 				uuid: "test-uuid",
// 			},
// 			want: &model.User{
// 				UUID:     "test-uuid",
// 				Username: "test",
// 				Fullname: "test",
// 				Email:    "test",
// 				Password: "test",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "failed_get_user_by_uuid",
// 			args: args{
// 				uuid: "",
// 			},
// 			want:    &model.User{},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			repo := &mocks.IUserRepo{}
// 			if !tt.wantErr {
// 				repo.On("GetByUUID", tt.args.uuid).Return(nil)
// 			} else {
// 				repo.On("GetByUUID", tt.args.uuid).Return(errors.New("user not found"))
// 			}
// 			got, err := repo.GetByUUID(tt.args.uuid)
// 			if tt.wantErr != (err != nil) {
// 				t.Errorf("Repository.GetByUUID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if tt.want != got {
// 				t.Errorf("UserRepo.GetByUUID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
