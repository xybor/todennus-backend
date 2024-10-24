package conversion

import (
	pbdto "github.com/xybor/todennus-backend/adapter/grpc/gen/dto"
	"github.com/xybor/todennus-backend/adapter/grpc/gen/dto/resource"
	ucdto "github.com/xybor/todennus-backend/usecase/dto"
)

func NewUsecaseUserValidateRequest(req *pbdto.UserValidateRequest) *ucdto.UserValidateCredentialsRequest {
	return &ucdto.UserValidateCredentialsRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

func NewPbUserValidateResponse(resp *ucdto.UserValidateCredentialsResponse) *pbdto.UserValidateResponse {
	return &pbdto.UserValidateResponse{
		User: &resource.User{
			Id:          resp.User.ID.Int64(),
			Username:    resp.User.Username,
			DisplayName: resp.User.DisplayName,
			Role:        resp.User.Role.String(),
		},
	}
}
