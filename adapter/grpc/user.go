package grpc

import (
	"context"

	"github.com/xybor/todennus-backend/adapter/abstraction"
	"github.com/xybor/todennus-backend/adapter/grpc/conversion"
	pb "github.com/xybor/todennus-backend/adapter/grpc/gen"
	pbdto "github.com/xybor/todennus-backend/adapter/grpc/gen/dto"
)

var _ pb.UserServer = (*UserServer)(nil)

type UserServer struct {
	pb.UnimplementedUserServer

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
	if err != nil {
		return nil, conversion.ReduceError(ctx, err)
	}

	return conversion.NewPbUserValidateResponse(resp), nil
}
