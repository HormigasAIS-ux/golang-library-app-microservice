package ucase

import (
	"auth_service/config"
	"auth_service/domain/dto"
	"auth_service/domain/model"
	author_grpc "auth_service/interface/grpc/genproto/author"
	"auth_service/repository"
	bcrypt_util "auth_service/utils/bcrypt"
	error_utils "auth_service/utils/error"
	"auth_service/utils/helper"
	jwt_util "auth_service/utils/jwt"
	validator_util "auth_service/utils/validator/user"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUcase struct {
	userRepo                repository.IUserRepo
	refreshTokenRepo        repository.IRefreshTokenRepo
	authorGrpcServiceClient author_grpc.AuthorServiceClient
}

type IAuthUcase interface {
	Register(ctx *gin.Context, payload dto.RegisterUserReq) (*dto.RegisterUserRespData, error)
	Login(payload dto.LoginReq) (*dto.LoginRespData, error)
	RefreshToken(payload dto.RefreshTokenReq) (*dto.RefreshTokenRespData, error)
	CheckToken(payload dto.CheckTokenReq) (*dto.CheckTokenRespData, error)
}

func NewAuthUcase(
	userRepo repository.IUserRepo,
	refreshTokenRepo repository.IRefreshTokenRepo,
	authorGrpcServiceClient author_grpc.AuthorServiceClient,
) IAuthUcase {
	return &AuthUcase{
		userRepo:                userRepo,
		refreshTokenRepo:        refreshTokenRepo,
		authorGrpcServiceClient: authorGrpcServiceClient,
	}
}

func (s *AuthUcase) Register(ctx *gin.Context, payload dto.RegisterUserReq) (*dto.RegisterUserRespData, error) {
	// validate input
	err := validator_util.ValidateUsername(payload.Username)
	if err != nil {
		logger.Errorf("error validating username: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	err = validator_util.ValidateEmail(payload.Email)
	if err != nil {
		logger.Errorf("error validating email: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	err = validator_util.ValidatePassword(payload.Password)
	if err != nil {
		logger.Errorf("error validating password: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	// check if user exists
	user, _ := s.userRepo.GetByEmail(payload.Email)
	logger.Debugf("user by email: %v", user)
	if user != nil {
		logger.Errorf("user with email %s already exists", payload.Email)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  fmt.Sprintf("user with email %s already exists", payload.Email),
		}
	}

	user, _ = s.userRepo.GetByUsername(payload.Username)
	if user != nil {
		logger.Errorf("user with username %s already exists", payload.Username)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  fmt.Sprintf("user with username %s already exists", payload.Username),
		}
	}

	// create password
	password, err := bcrypt_util.Hash(payload.Password)
	if err != nil {
		logger.Errorf("error hashing password: %v", err)
		return nil, err
	}

	// create author through author service
	_, err = s.authorGrpcServiceClient.CreateAuthor(
		ctx, &author_grpc.CreateAuthorReq{
			FirstName: payload.Username,
			LastName:  "",
			BirthDate: "",
			Bio:       "",
		},
	)
	grpcCode := status.Code(err)

	if grpcCode != codes.OK || err != nil {
		logger.Errorf("error creating author: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "error creating author",
			Detail:   err,
		}
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

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// generate token
	token, err := jwt_util.GenerateJwtToken(user, config.Envs.JWT_SECRET_KEY, config.Envs.JWT_EXP_HOURS, nil)
	if err != nil {
		logger.Errorf("error generating token: %v", err)
		return nil, err
	}

	// invalidate old refresh token
	s.refreshTokenRepo.InvalidateManyByUserUUID(user.UUID.String())

	// create refresh token
	refreshTokenExpiredAt := helper.TimeNowUTC().Add(time.Hour * time.Duration(config.Envs.JWT_REFRESH_EXP_HOURS))
	newRefreshTokenObj := model.RefreshToken{
		Token:     uuid.New().String(),
		UserUUID:  user.UUID,
		UsedAt:    nil,
		ExpiredAt: &refreshTokenExpiredAt,
	}
	logger.Debugf("new refresh token: %+v", newRefreshTokenObj)
	err = s.refreshTokenRepo.Create(&newRefreshTokenObj)
	if err != nil {
		logger.Errorf("error creating refresh token: %v", err)
		return nil, err
	}

	resp := &dto.RegisterUserRespData{
		AccessToken:  token,
		RefreshToken: newRefreshTokenObj.Token,
	}
	return resp, nil
}

func (s *AuthUcase) Login(payload dto.LoginReq) (*dto.LoginRespData, error) {
	// validate username
	if strings.Contains(payload.UsernameOrEmail, "@") {
		err := validator_util.ValidateEmail(payload.UsernameOrEmail)
		if err != nil {
			logger.Errorf("invalid username: %s\n%v", payload.UsernameOrEmail, err)
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  err.Error(),
			}
		}
	} else {
		err := validator_util.ValidateUsername(payload.UsernameOrEmail)
		if err != nil {
			logger.Errorf("invalid email: %s\n%v", payload.UsernameOrEmail, err)
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  err.Error(),
			}
		}
	}

	// validate password
	err := validator_util.ValidatePassword(payload.Password)
	if err != nil {
		logger.Errorf("invalid password: %s\n%v", payload.Password, err)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	// check if user exists
	var existing_user *model.User
	if strings.Contains(payload.UsernameOrEmail, "@") {
		existing_user, _ = s.userRepo.GetByEmail(payload.UsernameOrEmail)
	} else {
		existing_user, _ = s.userRepo.GetByUsername(payload.UsernameOrEmail)
	}
	if existing_user == nil {
		logger.Errorf("user not found")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Credentials",
		}
	}
	logger.Debugf("user by username or email: %v", helper.PrettyJson(existing_user))

	// check password
	if !bcrypt_util.Compare(payload.Password, existing_user.Password) {
		logger.Errorf("invalid password")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Credentials",
		}
	}

	// generate token
	token, err := jwt_util.GenerateJwtToken(existing_user, config.Envs.JWT_SECRET_KEY, config.Envs.JWT_EXP_HOURS, nil)
	if err != nil {
		logger.Errorf("error generating token: %v", err)
		return nil, err
	}

	// invalidate old refresh token
	err = s.refreshTokenRepo.InvalidateManyByUserUUID(existing_user.UUID.String())
	if err != nil {
		logger.Errorf("error invalidating old refresh token: %v", err)
		return nil, err
	}

	// create refresh token
	refreshTokenExpiredAt := helper.TimeNowUTC().Add(time.Hour * time.Duration(config.Envs.JWT_REFRESH_EXP_HOURS))
	newRefreshTokenObj := model.RefreshToken{
		Token:     uuid.New().String(),
		UserUUID:  existing_user.UUID,
		UsedAt:    nil,
		ExpiredAt: &refreshTokenExpiredAt,
	}
	logger.Debugf("new refresh token: %+v", helper.PrettyJson(newRefreshTokenObj))
	err = s.refreshTokenRepo.Create(&newRefreshTokenObj)
	if err != nil {
		logger.Errorf("error creating refresh token: %v", err)
		return nil, err
	}

	return &dto.LoginRespData{
		AccessToken:  token,
		RefreshToken: newRefreshTokenObj.Token,
	}, nil
}

func (s *AuthUcase) RefreshToken(payload dto.RefreshTokenReq) (*dto.RefreshTokenRespData, error) {
	// get refresh token
	refreshToken, err := s.refreshTokenRepo.GetByToken(payload.RefreshToken)
	if err != nil {
		logger.Errorf("refresh token not found: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Refresh Token",
		}
	}

	// check if refresh token is expired
	if refreshToken.ExpiredAt != nil {
		if refreshToken.ExpiredAt.Before(helper.TimeNowUTC()) {
			logger.Errorf("refresh token is expired")
			return nil, &error_utils.CustomErr{
				HttpCode: 401,
				Message:  "Invalid Refresh Token",
			}
		}
	}

	// check if refresh token is used
	if refreshToken.UsedAt != nil {
		logger.Errorf("refresh token is used")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Refresh Token",
		}
	}

	// check if refresh token is valid
	if refreshToken.Invalid {
		logger.Errorf("refresh token is invalid")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Refresh Token",
		}
	}

	// mark refresh token as used
	timeNow := helper.TimeNowUTC()
	refreshToken.UsedAt = &timeNow
	err = s.refreshTokenRepo.Update(refreshToken)
	if err != nil {
		logger.Errorf("error updating refresh token: %v", err)
		return nil, err
	}

	// get user
	user, err := s.userRepo.GetByUUID(refreshToken.UserUUID.String())
	if err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "Internal server error",
			Detail:   err.Error(),
		}
	}

	// generate token
	token, err := jwt_util.GenerateJwtToken(user, config.Envs.JWT_SECRET_KEY, config.Envs.JWT_EXP_HOURS, nil)
	if err != nil {
		logger.Errorf("error generating token: %v", err)
		return nil, err
	}

	// invalidate old refresh token
	err = s.refreshTokenRepo.InvalidateManyByUserUUID(user.UUID.String())
	if err != nil {
		logger.Errorf("error invalidating old refresh token: %v", err)
		return nil, err
	}

	// create refresh token
	refreshTokenExpiredAt := helper.TimeNowUTC().Add(time.Hour * time.Duration(config.Envs.JWT_REFRESH_EXP_HOURS))
	newRefreshTokenObj := model.RefreshToken{
		Token:     uuid.New().String(),
		UserUUID:  user.UUID,
		UsedAt:    nil,
		ExpiredAt: &refreshTokenExpiredAt,
	}
	err = s.refreshTokenRepo.Create(&newRefreshTokenObj)
	if err != nil {
		logger.Errorf("error creating refresh token: %v", err)
		return nil, err
	}

	return &dto.RefreshTokenRespData{
		AccessToken:  token,
		RefreshToken: newRefreshTokenObj.Token,
	}, nil
}

func (s *AuthUcase) CheckToken(payload dto.CheckTokenReq) (*dto.CheckTokenRespData, error) {
	claims, err := jwt_util.ValidateJWT(payload.AccessToken, config.Envs.JWT_SECRET_KEY)
	if err != nil || claims == nil {
		logger.Errorf("error validating token: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			GrpcCode: codes.Unauthenticated,
			Message:  "Invalid Access Token",
			Detail:   err.Error(),
		}
	}

	resp := &dto.CheckTokenRespData{
		UUID:     claims.UUID,
		Username: claims.Username,
		Role:     claims.Role,
		Email:    claims.Email,
	}

	return resp, nil
}
