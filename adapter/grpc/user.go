package grpc

import (
	"context"

	"github.com/xybor/todennus-backend/adapter/abstraction"
	"github.com/xybor/todennus-backend/adapter/grpc/conversion"
	service "github.com/xybor/todennus-backend/adapter/grpc/gen"
	pbdto "github.com/xybor/todennus-backend/adapter/grpc/gen/dto"
	"github.com/xybor/todennus-backend/usecase"
	"google.golang.org/grpc/codes"
)

var _ service.UserServer = (*UserServer)(nil)

type UserServer struct {
	service.UnimplementedUserServer

	userUsecase abstraction.UserUsecase
}

func NewUserServer(userUsecase abstraction.UserUsecase) *UserServer {
	return &UserServer{
		userUsecase: userUsecase,
	}
}

func (s *UserServer) Validate(ctx context.Context, req *pbdto.UserValidateRequest) (*pbdto.UserValidateResponse, error) {
	ucreq := conversion.NewUsecaseUserValidateRequest(req)
	resp, err := s.userUsecase.ValidateCredentials(ctx, ucreq)

	return conversion.NewResponseHandler(ctx, conversion.NewPbUserValidateResponse(resp), err).
		Map(codes.InvalidArgument, usecase.ErrRequestInvalid).
		Map(codes.PermissionDenied, usecase.ErrCredentialsInvalid).
		Map(codes.NotFound, usecase.ErrNotFound).Finalize(ctx)
}
