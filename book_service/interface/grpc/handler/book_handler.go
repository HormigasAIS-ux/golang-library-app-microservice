package handler

import (
	"book_service/domain/dto"
	book_grpc "book_service/interface/grpc/genproto/book"
	ucase "book_service/usecase"
	error_utils "book_service/utils/error"
	"context"

	"github.com/op/go-logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookServiceHandler struct {
	book_grpc.UnimplementedBookServiceServer
	bookUcase ucase.IBookUcase
}

var logger = logging.MustGetLogger("main")

func NewAuthorServiceHandler(bookUcase ucase.IBookUcase) *BookServiceHandler {
	handler := &BookServiceHandler{bookUcase: bookUcase}
	return handler
}

func (r *BookServiceHandler) GetBookTotalByAuthorUUID(
	ctx context.Context,
	in *book_grpc.GetBookTotalByAuthorUUIDReq,
) (*book_grpc.GetBookTotalByAuthorUUIDResp, error) {
	logger.Debugf("incoming request: %v", in)

	raw, err := r.bookUcase.GetBookTotalByAuthorUUID(ctx, in.AuthorUuid)
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &book_grpc.GetBookTotalByAuthorUUIDResp{
		BookTotal: raw,
	}
	return resp, nil
}

func (r *BookServiceHandler) BulkGetBookTotalByAuthorUUIDs(
	ctx context.Context,
	in *book_grpc.BulkGetBookTotalByAuthorUUIDsReq,
) (*book_grpc.BulkGetBookTotalByAuthorUUIDsResp, error) {
	logger.Debugf("incoming request: %v", in)

	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(in.AuthorUuids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "author uuids are required")
	}

	payloadDto := dto.BulkGetBookTotalByAuthorUUIDsReq{
		AuthorUUIDs: in.AuthorUuids,
	}
	raw, err := r.bookUcase.BulkGetBookTotalByAuthorUUIDs(ctx, payloadDto)
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &book_grpc.BulkGetBookTotalByAuthorUUIDsResp{}
	for _, item := range raw {
		resp.Data = append(resp.Data, &book_grpc.BulkGetBookTotalByAuthorUUIDsResp_Data{
			AuthorUuid: item.AuthorUUID,
			BookTotal:  item.Total,
		})
	}
	return resp, nil
}
