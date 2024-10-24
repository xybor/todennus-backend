package dto

import (
	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/adapter/rest/dto/resource"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/xerror"
)

func ParseUserID(meID snowflake.ID, s string) (snowflake.ID, error) {
	if s == "@me" {
		return meID, nil
	}

	return snowflake.ParseString(s)
}

// Register
type UserRegisterRequest struct {
	Username string `json:"username" example:"huykingsofm"`
	Password string `json:"password" example:"s3Cr3tP@ssW0rD"`
}

func (req UserRegisterRequest) To() *dto.UserRegisterRequest {
	return &dto.UserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserRegisterResponse struct {
	*resource.User
}

func NewUserRegisterResponse(resp *dto.UserRegisterResponse) *UserRegisterResponse {
	if resp == nil {
		return nil
	}

	return &UserRegisterResponse{
		User: resource.NewUser(resp.User),
	}
}

// GetByID
type UserGetByIDRequest struct {
	UserID string `param:"user_id"`
}

func (req UserGetByIDRequest) To(meID snowflake.ID) (*dto.UserGetByIDRequest, error) {
	userID, err := ParseUserID(meID, req.UserID)
	if err != nil {
		return nil, xerror.Enrich(usecase.ErrRequestInvalid, "user id is invalid").
			Hide(err, "failed-to-parse-user-id", "uid", req.UserID)
	}

	return &dto.UserGetByIDRequest{UserID: userID}, nil
}

type UserGetByIDResponse struct {
	*resource.User
}

func NewUserGetByIDResponse(resp *dto.UserGetByIDResponse) *UserGetByIDResponse {
	if resp == nil {
		return nil
	}

	return &UserGetByIDResponse{
		User: resource.NewUser(resp.User),
	}
}

// GetByUsername
type UserGetByUsernameRequest struct {
	Username string `param:"username"`
}

func (req UserGetByUsernameRequest) To() *dto.UserGetByUsernameRequest {
	return &dto.UserGetByUsernameRequest{
		Username: req.Username,
	}
}

type UserGetByUsernameResponse struct {
	*resource.User
}

func NewUserGetByUsernameResponse(resp *dto.UserGetByUsernameResponse) *UserGetByUsernameResponse {
	if resp == nil {
		return nil
	}

	return &UserGetByUsernameResponse{
		User: resource.NewUser(resp.User),
	}
}

// Validate
type UserValidateRequest struct {
	Username string `json:"username" example:"huykingsofm"`
	Password string `json:"password" example:"s3Cr3tP@ssW0rD"`
}

func (req UserValidateRequest) To() *dto.UserValidateCredentialsRequest {
	return &dto.UserValidateCredentialsRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserValidateResponse struct {
	*resource.User
}

func NewUserValidateResponse(resp *dto.UserValidateCredentialsResponse) *UserValidateResponse {
	if resp == nil {
		return nil
	}

	return &UserValidateResponse{
		User: resource.NewUser(resp.User),
	}
}
