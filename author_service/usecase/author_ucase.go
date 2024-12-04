package ucase

import (
	"author_service/domain/dto"
	"author_service/domain/model"
	auth_pb "author_service/interface/grpc/genproto/auth"
	book_pb "author_service/interface/grpc/genproto/book"
	"author_service/repository"
	error_utils "author_service/utils/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

type AuthorUcase struct {
	authorRepo repository.IAuthorRepo
}

type IAuthorUcase interface {
	CreateNewAuthor(
		ctx *gin.Context,
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
}

func NewAuthorUcase(authorRepo repository.IAuthorRepo) IAuthorUcase {
	return &AuthorUcase{
		authorRepo: authorRepo,
	}
}

func (u *AuthorUcase) CreateNewAuthor(
	ctx *gin.Context,
	payload dto.CreateNewAuthorReq,
) (*dto.CreateNewAuthorRespData, error) {
	// TODO: create new user through auth service
	createdUser := auth_pb.CreateUserResp{}
	parsedUserUUID, err := uuid.Parse(createdUser.Uuid)
	if err != nil {
		return nil, err
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
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
			Detail:   err.Error(),
		}
	}

	// create
	err = u.authorRepo.Create(newAuthor)
	if err != nil {
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
		Email:     payload.Email,
		Username:  payload.Username,
		Role:      createdUser.Role,
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

	updateUserReq := auth_pb.UpdateUserReq{Uuid: author.UserUUID.String()}
	if payload.Username != nil {
		updateUserReq.Username = *payload.Username
	} else {
		updateUserReq.UsernameNull = true
	}

	if payload.Email != nil {
		updateUserReq.Email = *payload.Email
	} else {
		updateUserReq.EmailNull = true
	}

	if payload.Password != nil {
		updateUserReq.Password = *payload.Password
	} else {
		updateUserReq.PasswordNull = true
	}

	if payload.Role != nil {
		updateUserReq.Role = *payload.Role
	} else {
		updateUserReq.RoleNull = true
	}

	// TODO: edit user through auth service
	updatedUser := auth_pb.UpdateUserResp{}
	if err != nil {
		return nil, err
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
		Email:     updatedUser.Email,
		Username:  updatedUser.Username,
		Role:      updatedUser.Role,
	}

	return respData, nil
}

func (u *AuthorUcase) DeleteAuthor(ctx *gin.Context, authorUUID string) (*dto.DeleteAuthorRespData, error) {
	// TODO: delete user through auth service
	deletedUser := auth_pb.DeleteUserResp{}

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
		Email:     deletedUser.Email,
		Username:  deletedUser.Username,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		BirthDate: author.BirthDate,
		Bio:       author.Bio,
		Role:      deletedUser.Role,
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

	// TODO: get user through auth service
	user := auth_pb.GetUserByUUIDResponse{}

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

	// TODO: get book total by author uuid through book service
	bookTotalResp := book_pb.GetBookTotalByAuthorUUIDResp{}

	respData := &dto.GetAuthorDetailRespData{
		UUID:      author.UUID,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
		UserUUID:  author.UserUUID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		BirthDate: author.BirthDate,
		Bio:       author.Bio,
		Role:      user.Role,
		BookTotal: bookTotalResp.BookTotal,
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

	// TODO: get bulk book total by author uuid through book service
	bookTotalResp := book_pb.BulkGetBookTotalByAuthorUUIDsResp{}
	bookTotalByAuthorUUID := make(map[string]int64)

	if bookTotalResp.Data != nil {
		for _, item := range bookTotalResp.Data {
			bookTotalByAuthorUUID[item.AuthorUuid] = item.BookTotal
		}
	}

	respItems := make([]dto.GetAuthorListRespDataItem, 0)
	for _, v := range data {
		bookTotal, ok := bookTotalByAuthorUUID[v.UUID.String()]
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
