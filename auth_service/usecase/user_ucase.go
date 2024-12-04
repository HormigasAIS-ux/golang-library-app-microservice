package ucase

import (
	"auth_service/domain/dto"
	"auth_service/domain/model"
	"auth_service/repository"
	bcrypt_util "auth_service/utils/bcrypt"
	error_utils "auth_service/utils/error"
	validator_util "auth_service/utils/validator/user"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type UserUcase struct {
	userRepo repository.IUserRepo
}

type IUserUcase interface {
	GetByUUID(ctx context.Context, ginCtx *gin.Context, userUUID string) (*dto.GetUserByUUIDResp, error)
	CreateUser(
		ctx context.Context,
		ginCtx *gin.Context,
		payload dto.CreateUserReq,
	) (*dto.CreateUserRespData, error)
	UpdateUser(
		ctx context.Context,
		ginCtx *gin.Context,
		userUUID string,
		payload dto.UpdateUserReq,
	) (*dto.UpdateUserRespData, error)
	DeleteUser(
		ctx context.Context,
		ginCtx *gin.Context,
		userUUID string,
	) (*dto.DeleteUserRespData, error)
}

func NewUserUcase(userRepo repository.IUserRepo) IUserUcase {
	return &UserUcase{userRepo: userRepo}
}

func (ucase *UserUcase) GetByUUID(ctx context.Context, ginCtx *gin.Context, userUUID string) (*dto.GetUserByUUIDResp, error) {
	user, err := ucase.userRepo.GetByUUID(userUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	return &dto.GetUserByUUIDResp{
		UUID:      user.UUID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (ucase *UserUcase) CreateUser(
	ctx context.Context,
	ginCtx *gin.Context,
	payload dto.CreateUserReq,
) (*dto.CreateUserRespData, error) {
	// validate input
	err := validator_util.ValidateUsername(payload.Username)
	if err != nil {
		logger.Errorf("error validating username: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.InvalidArgument,
			Message:  err.Error(),
		}
	}

	err = validator_util.ValidateEmail(payload.Email)
	if err != nil {
		logger.Errorf("error validating email: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.InvalidArgument,
			Message:  err.Error(),
		}
	}

	err = validator_util.ValidatePassword(payload.Password)
	if err != nil {
		logger.Errorf("error validating password: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.InvalidArgument,
			Message:  err.Error(),
		}
	}

	err = validator_util.ValidateRole(payload.Role)
	if err != nil {
		logger.Errorf("error validating role: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.InvalidArgument,
			Message:  err.Error(),
		}
	}

	// check if user exists
	user, _ := ucase.userRepo.GetByEmail(payload.Email)
	logger.Debugf("user by email: %v", user)
	if user != nil {
		logger.Errorf("user with email %s already exists", payload.Email)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.AlreadyExists,
			Message:  fmt.Sprintf("user with email %s already exists", payload.Email),
		}
	}

	user, _ = ucase.userRepo.GetByUsername(payload.Username)
	if user != nil {
		logger.Errorf("user with username %s already exists", payload.Username)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.AlreadyExists,
			Message:  fmt.Sprintf("user with username %s already exists", payload.Username),
		}
	}

	// create password
	password, err := bcrypt_util.Hash(payload.Password)
	if err != nil {
		logger.Errorf("error hashing password: %v", err)
		return nil, err
	}

	// create user
	user = &model.User{
		UUID:     uuid.New(),
		Username: payload.Username,
		Password: password,
		Email:    payload.Email,
		Role:     "user",
	}
	err = user.Validate()
	if err != nil {
		return nil, err
	}

	err = ucase.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserRespData{
		UUID:      user.UUID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (ucase *UserUcase) UpdateUser(
	ctx context.Context,
	ginCtx *gin.Context,
	userUUID string,
	payload dto.UpdateUserReq,
) (*dto.UpdateUserRespData, error) {
	// validate input
	if payload.Username != nil {
		err := validator_util.ValidateUsername(*payload.Username)
		if err != nil {
			logger.Errorf("error validating username: %s", err.Error())
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				GrpcCode: codes.InvalidArgument,
				Message:  err.Error(),
			}
		}
	}

	if payload.Email != nil {
		err := validator_util.ValidateEmail(*payload.Email)
		if err != nil {
			logger.Errorf("error validating email: %s", err.Error())
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				GrpcCode: codes.InvalidArgument,
				Message:  err.Error(),
			}
		}
	}

	if payload.Password != nil {
		err := validator_util.ValidatePassword(*payload.Password)
		if err != nil {
			logger.Errorf("error validating password: %s", err.Error())
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				GrpcCode: codes.InvalidArgument,
				Message:  err.Error(),
			}
		}
	}

	if payload.Role != nil {
		err := validator_util.ValidateRole(*payload.Role)
		if err != nil {
			logger.Errorf("error validating role: %s", err.Error())
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				GrpcCode: codes.InvalidArgument,
				Message:  err.Error(),
			}
		}
	}

	// get existing user
	user, err := ucase.userRepo.GetByUUID(userUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// update user obj
	if payload.Username != nil {
		user.Username = *payload.Username
	}
	if payload.Email != nil {
		user.Email = *payload.Email
	}
	if payload.Password != nil {
		password, err := bcrypt_util.Hash(*payload.Password)
		if err != nil {
			logger.Errorf("error hashing password: %v", err)
			return nil, err
		}
		user.Password = password
	}
	if payload.Role != nil {
		user.Role = *payload.Role
	}

	// update user
	err = ucase.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateUserRespData{
		UUID:      user.UUID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (ucase *UserUcase) DeleteUser(
	ctx context.Context,
	ginCtx *gin.Context,
	userUUID string,
) (*dto.DeleteUserRespData, error) {
	// find user
	user, err := ucase.userRepo.GetByUUID(userUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// delete user
	err = ucase.userRepo.Delete(user.UUID.String())
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	return &dto.DeleteUserRespData{
		UUID:      user.UUID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
