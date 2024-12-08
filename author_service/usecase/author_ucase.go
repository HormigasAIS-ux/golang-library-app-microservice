package ucase

import (
	"author_service/domain/dto"
	"author_service/domain/model"
	auth_pb "author_service/interface/grpc/genproto/auth"
	book_pb "author_service/interface/grpc/genproto/book"
	"author_service/repository"
	error_utils "author_service/utils/error"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/op/go-logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = logging.MustGetLogger("main")

type AuthorUcase struct {
	authorRepo            repository.IAuthorRepo
	authGrpcServiceClient auth_pb.AuthServiceClient
	bookGrpcServiceClient book_pb.BookServiceClient
}

type IAuthorUcase interface {
	CreateNewAuthor(
		ctx context.Context,
		payload dto.CreateNewAuthorReq,
	) (*dto.CreateNewAuthorRespData, error) // admin only
	EditAuthor(
		ctx *gin.Context,
		authorUUID string,
		payload dto.EditAuthorReq,
	) (*dto.EditAuthorRespData, error) // admin only or owner
	DeleteAuthor(ctx *gin.Context, authorUUID string) (*dto.DeleteAuthorRespData, error) // admin only
	GetAuthorDetail(ctx *gin.Context, authorUUID string) (*dto.GetAuthorDetailRespData, error)
	GetList(
		ctx *gin.Context, query dto.GetAuthorListReq,
	) ([]dto.GetAuthorListRespDataItem, int64, error)
	GetAuthorByUserUUID(
		ctx context.Context, userUUID string,
	) (*dto.GetAuthorByUserUUIDRespData, error)
}

func NewAuthorUcase(
	authorRepo repository.IAuthorRepo,
	authGrpcServiceClient auth_pb.AuthServiceClient,
	bookGrpcServiceClient book_pb.BookServiceClient,
) IAuthorUcase {
	return &AuthorUcase{
		authorRepo:            authorRepo,
		authGrpcServiceClient: authGrpcServiceClient,
		bookGrpcServiceClient: bookGrpcServiceClient,
	}
}

func (u *AuthorUcase) CreateNewAuthor(
	ctx context.Context,
	payload dto.CreateNewAuthorReq,
) (*dto.CreateNewAuthorRespData, error) {
	logger.Debugf("CreateNewAuthor in")
	var parsedUserUUID uuid.UUID
	var userEmail string
	var userUsername string
	var userRole string
	var err error

	if payload.UserUUID != nil { // user uuid provided for auth service grpc call
		logger.Debugf("payload.UserUUID: %s", *payload.UserUUID)
		getUserResp, err := u.authGrpcServiceClient.GetUserByUUID(
			ctx,
			&auth_pb.GetUserByUUIDRequest{
				Uuid: *payload.UserUUID,
			},
		)
		grpcCode := status.Code(err)
		logger.Debugf("getUserResp: %v, grpcCode: %v, err: %v", getUserResp, grpcCode, err)
		if grpcCode != codes.OK {
			switch grpcCode {
			case codes.NotFound:
				logger.Errorf("user not found: %s", err.Error())
				return nil, &error_utils.CustomErr{
					HttpCode: 400,
					GrpcCode: grpcCode,
					Message:  err.Error(),
					Detail:   err.Error(),
				}
			default:
				logger.Errorf("error getting user: %s", err.Error())
				return nil, &error_utils.CustomErr{
					HttpCode: 500,
					GrpcCode: grpcCode,
					Message:  err.Error(),
					Detail:   err.Error(),
				}
			}
		}

		if getUserResp == nil {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				GrpcCode: grpcCode,
				Message:  "internal server error",
				Detail:   "get user response is nil",
			}
		}

		parsedUserUUID, err = uuid.Parse(getUserResp.Uuid)
		if err != nil {
			logger.Errorf("error parsing user uuid: %s", err.Error())
			return nil, err
		}
		userEmail = getUserResp.Email
		userUsername = getUserResp.Username
		userRole = getUserResp.Role

	} else { // user uuid not provided for client call
		createUserResp, err := u.authGrpcServiceClient.CreateUser(
			ctx,
			&auth_pb.CreateUserReq{
				Email:    payload.Email,
				Username: payload.Username,
				Password: payload.Password,
				Role:     payload.Role,
			},
		)
		grpcCode := status.Code(err)

		if grpcCode != codes.OK {
			switch grpcCode {
			case codes.AlreadyExists:
				logger.Errorf("user already exists: %s", err.Error())
				return nil, &error_utils.CustomErr{
					HttpCode: 400,
					GrpcCode: grpcCode,
					Message:  err.Error(),
					Detail:   err.Error(),
				}
			default:
				logger.Errorf("error creating user: %s", err.Error())
				return nil, &error_utils.CustomErr{
					HttpCode: 500,
					GrpcCode: grpcCode,
					Message:  err.Error(),
					Detail:   err.Error(),
				}
			}
		}

		if createUserResp == nil {
			logger.Errorf("create user resp is nil")
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				GrpcCode: grpcCode,
				Message:  "internal server error",
				Detail:   "create user resp is nil",
			}
		}

		parsedUserUUID, err = uuid.Parse(createUserResp.Uuid)
		if err != nil {
			logger.Errorf("error parsing user uuid: %s", err.Error())
			return nil, err
		}

		userEmail = createUserResp.Email
		userUsername = createUserResp.Username
		userRole = createUserResp.Role
	}

	newAuthor := &model.Author{
		UUID:      uuid.New(),
		UserUUID:  parsedUserUUID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		BirthDate: payload.BirthDate,
		Bio:       payload.Bio,
	}

	// validate
	err = newAuthor.Validate()
	if err != nil {
		logger.Errorf("author validation error: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
			Detail:   err.Error(),
		}
	}

	// create
	err = u.authorRepo.Create(newAuthor)
	if err != nil {
		logger.Errorf("error creating author: %s", err.Error())
		return nil, err
	}

	respData := &dto.CreateNewAuthorRespData{
		UUID:      newAuthor.UUID,
		CreatedAt: newAuthor.CreatedAt,
		UpdatedAt: newAuthor.UpdatedAt,
		UserUUID:  newAuthor.UserUUID,
		FirstName: newAuthor.FirstName,
		LastName:  newAuthor.LastName,
		BirthDate: newAuthor.BirthDate,
		Bio:       newAuthor.Bio,
		Email:     userEmail,
		Username:  userUsername,
		Role:      userRole,
	}

	return respData, nil
}

func (u *AuthorUcase) EditAuthor(
	ctx *gin.Context,
	authorUUID string,
	payload dto.EditAuthorReq,
) (*dto.EditAuthorRespData, error) {
	// handle authorUUID me
	if authorUUID == "me" {
		currentUserRaw, ok := ctx.Get("currentUser")
		if !ok {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				Message:  "internal server error",
				Detail:   "current user not found",
			}
		}

		currentUser, ok := currentUserRaw.(dto.CurrentUser)
		if !ok {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				Message:  "internal server error",
				Detail:   "current user missmatched",
			}
		}

		myAuthor, err := u.authorRepo.GetByUserUUID(currentUser.UUID)
		if err != nil {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				Message:  "internal server error",
				Detail:   err.Error(),
			}
		}

		authorUUID = myAuthor.UUID.String()
	}

	// get existing author
	author, err := u.authorRepo.GetByUUID(authorUUID)
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 404,
			Message:  "user not found",
			Detail:   err.Error(),
		}
	}

	// update user through auth service grpc
	updateUserReqPayload := auth_pb.UpdateUserReq{Uuid: author.UserUUID.String()}
	if payload.Username != nil {
		updateUserReqPayload.Username = *payload.Username
	} else {
		updateUserReqPayload.UsernameNull = true
	}

	if payload.Email != nil {
		updateUserReqPayload.Email = *payload.Email
	} else {
		updateUserReqPayload.EmailNull = true
	}

	if payload.Password != nil {
		updateUserReqPayload.Password = *payload.Password
	} else {
		updateUserReqPayload.PasswordNull = true
	}

	if payload.Role != nil {
		updateUserReqPayload.Role = *payload.Role
	} else {
		updateUserReqPayload.RoleNull = true
	}

	updateUserResp, err := u.authGrpcServiceClient.UpdateUser(ctx, &updateUserReqPayload)
	grpcCode := status.Code(err)
	if grpcCode != codes.OK || err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err.Error(),
		}
	}

	if updateUserResp == nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   "update user response is nil",
		}
	}

	// prepare update author
	if payload.FirstName != nil {
		author.FirstName = *payload.FirstName
	}
	if payload.LastName != nil {
		author.LastName = *payload.LastName
	}
	if payload.BirthDate != nil {
		author.BirthDate = payload.BirthDate
	}
	if payload.Bio != nil {
		author.Bio = payload.Bio
	}

	// validate
	err = author.Validate()
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
			Detail:   err.Error(),
		}
	}

	// update
	err = u.authorRepo.Update(author)
	if err != nil {
		return nil, err
	}

	respData := &dto.EditAuthorRespData{
		UUID:      author.UUID,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
		UserUUID:  author.UserUUID,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		BirthDate: author.BirthDate,
		Bio:       author.Bio,
		Email:     updateUserResp.Email,
		Username:  updateUserResp.Username,
		Role:      updateUserResp.Role,
	}

	return respData, nil
}

func (u *AuthorUcase) DeleteAuthor(ctx *gin.Context, authorUUID string) (*dto.DeleteAuthorRespData, error) {

	author, err := u.authorRepo.GetByUUID(authorUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// delete user through auth service
	deleteUserResp, err := u.authGrpcServiceClient.DeleteUser(ctx, &auth_pb.DeleteUserReq{Uuid: authorUUID})
	grpcCode := status.Code(err)
	if grpcCode != codes.OK || err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err.Error(),
		}
	}

	if deleteUserResp == nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   "delete user response is nil",
		}
	}

	// delete author
	err = u.authorRepo.Delete(author.UUID.String())
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}
	return &dto.DeleteAuthorRespData{
		UUID:      author.UUID,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
		UserUUID:  author.UserUUID,
		Email:     deleteUserResp.Email,
		Username:  deleteUserResp.Username,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		BirthDate: author.BirthDate,
		Bio:       author.Bio,
		Role:      deleteUserResp.Role,
	}, nil
}

func (u *AuthorUcase) GetAuthorDetail(ctx *gin.Context, authorUUID string) (*dto.GetAuthorDetailRespData, error) {
	// handle authorUUID me
	if authorUUID == "me" {
		currentUserRaw, ok := ctx.Get("currentUser")
		if !ok {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				Message:  "internal server error",
				Detail:   "current user not found",
			}
		}

		currentUser, ok := currentUserRaw.(dto.CurrentUser)
		if !ok {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				Message:  "internal server error",
				Detail:   "current user missmatched",
			}
		}

		myAuthor, err := u.authorRepo.GetByUserUUID(currentUser.UUID)
		if err != nil {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				Message:  "internal server error",
				Detail:   err.Error(),
			}
		}

		authorUUID = myAuthor.UUID.String()
	}

	author, err := u.authorRepo.GetByUUID(authorUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// get user through auth service
	getUserResp, err := u.authGrpcServiceClient.GetUserByUUID(
		ctx,
		&auth_pb.GetUserByUUIDRequest{
			Uuid: author.UserUUID.String(),
		},
	)
	grpcCode := status.Code(err)
	if grpcCode != codes.OK {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: grpcCode,
			Message:  err.Error(),
			Detail:   err.Error(),
		}
	}

	if getUserResp == nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: grpcCode,
			Message:  "internal server error",
			Detail:   "user resp is nil",
		}
	}

	resp, err := u.bookGrpcServiceClient.GetBookTotalByAuthorUUID(
		ctx,
		&book_pb.GetBookTotalByAuthorUUIDReq{
			AuthorUuid: author.UUID.String(),
		},
	)
	code := status.Code(err)
	if code != codes.OK || err != nil {
		logger.Errorf("failed to get book total by author uuid: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	respData := &dto.GetAuthorDetailRespData{
		UUID:      author.UUID,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
		UserUUID:  author.UserUUID,
		Email:     getUserResp.Email,
		Username:  getUserResp.Username,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		BirthDate: author.BirthDate,
		Bio:       author.Bio,
		Role:      getUserResp.Role,
		BookTotal: resp.BookTotal,
	}

	return respData, nil
}

func (u *AuthorUcase) GetList(
	ctx *gin.Context, query dto.GetAuthorListReq,
) ([]dto.GetAuthorListRespDataItem, int64, error) {
	// handle query by
	queryBy := query.QueryBy
	if query.QueryBy == "any" {
		queryBy = ""
	}

	data, err := u.authorRepo.GetList(ctx, dto.AuthorRepo_GetListParams{
		Query:     query.Query,
		QueryBy:   queryBy,
		Page:      query.Page,
		Limit:     query.Limit,
		SortOrder: query.SortOrder,
		SortBy:    query.SortBy,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := u.authorRepo.CountGetList(ctx, dto.AuthorRepo_GetListParams{
		Query:   query.Query,
		QueryBy: query.QueryBy,
	})
	if err != nil {
		return nil, 0, err
	}

	resp, err := u.bookGrpcServiceClient.BulkGetBookTotalByAuthorUUIDs(
		ctx,
		&book_pb.BulkGetBookTotalByAuthorUUIDsReq{
			AuthorUuids: []string{},
		},
	)
	code := status.Code(err)
	if code != codes.OK || err != nil {
		logger.Warningf("failed to get book total by author uuids: %v", err)
	}

	bookTotalMapByAuthorUUID := make(map[string]int64)

	if resp != nil {
		if resp.Data != nil {
			for _, item := range resp.Data {
				if item != nil {
					if item.AuthorUuid == "" {
						logger.Warningf("failed to get book total; author uuid is empty; skip")
						continue
					}
					bookTotalMapByAuthorUUID[item.AuthorUuid] = item.BookTotal
				}
			}
		}
	}

	respItems := make([]dto.GetAuthorListRespDataItem, 0)
	for _, v := range data {
		bookTotal, ok := bookTotalMapByAuthorUUID[v.UUID.String()]
		if !ok {
			logger.Warningf("book total not found for author uuid: %s; set to 0", v.UUID.String())
			bookTotal = 0
		}
		respItems = append(respItems, dto.GetAuthorListRespDataItem{
			UUID:      v.UUID.String(),
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			BirthDate: v.BirthDate,
			Bio:       v.Bio,
			BookTotal: bookTotal,
		})
	}

	return respItems, count, nil
}

func (u *AuthorUcase) GetAuthorByUserUUID(
	ctx context.Context, userUUID string,
) (*dto.GetAuthorByUserUUIDRespData, error) {
	author, err := u.authorRepo.GetByUserUUID(userUUID)
	if err != nil {
		if err.Error() == "not found" {
			logger.Errorf("author not found: %s", userUUID)
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "author not found",
				Detail:   err.Error(),
			}
		}
		logger.Errorf("error getting author by user uuid: %s", userUUID)
		return nil, err
	}

	return &dto.GetAuthorByUserUUIDRespData{
		UUID:      author.UUID,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		BirthDate: author.BirthDate,
		Bio:       author.Bio,
	}, nil
}
