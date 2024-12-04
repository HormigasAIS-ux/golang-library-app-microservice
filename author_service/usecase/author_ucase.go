package ucase

import (
	"author_service/domain/dto"
	"author_service/repository"
	"context"
)

type AuthorUcase struct {
	authorRepo repository.IAuthorRepo
}

type IAuthorUcase interface {
	GetList(
		ctx context.Context, query dto.GetAuthorListReq,
	) ([]dto.GetAuthorListRespDataItem, int64, error)
}

func NewAuthorUcase(authorRepo repository.IAuthorRepo) IAuthorUcase {
	return &AuthorUcase{
		authorRepo: authorRepo,
	}
}

func (u *AuthorUcase) GetList(
	ctx context.Context, query dto.GetAuthorListReq,
) ([]dto.GetAuthorListRespDataItem, int64, error) {
	// handle query by
	if query.QueryBy == "any" {
		query.QueryBy = ""
	}

	data, err := u.authorRepo.GetList(ctx, dto.AuthorRepo_GetListParams{
		Query:     query.Query,
		QueryBy:   query.QueryBy,
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

	respItems := make([]dto.GetAuthorListRespDataItem, 0)
	for _, v := range data {
		respItems = append(respItems, dto.GetAuthorListRespDataItem{
			UUID:      v.UUID.String(),
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			BirthDate: v.BirthDate,
			Bio:       v.Bio,
		})
	}

	return respItems, count, nil
}
