package usecase

import (
	"context"
	"errors"

	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/pkg/xerror"
	"github.com/xybor/todennus-backend/usecase/abstraction"
	"github.com/xybor/todennus-backend/usecase/dto"
)

type UserUsecase struct {
	userDomain abstraction.UserDomain
	userRepo   abstraction.UserRepository
}

func NewUserUsecase(userRepo abstraction.UserRepository, userDomain abstraction.UserDomain) *UserUsecase {
	return &UserUsecase{
		userRepo:   userRepo,
		userDomain: userDomain,
	}
}

func (uc *UserUsecase) Register(
	ctx context.Context,
	req dto.UserRegisterRequestDTO,
) (dto.UserRegisterResponseDTO, error) {
	_, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return dto.UserRegisterResponseDTO{}, xerror.WrapDebug(ErrUsernameExisted).
			WithMessage("username %s has already taken before", req.Username)
	}

	if !errors.Is(err, database.ErrRecordNotFound) {
		return dto.UserRegisterResponseDTO{}, wrapNonDomainError(xerror.ServerityCritical, err)
	}

	user, err := uc.userDomain.Create(req.Username, req.Password)
	if err != nil {
		return dto.UserRegisterResponseDTO{}, wrapDomainError(err)
	}

	err = uc.userRepo.Save(ctx, user)
	if err != nil {
		return dto.UserRegisterResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.UserRegisterResponseDTO{User: user}, nil
}

func (uc *UserUsecase) Validate(
	ctx context.Context,
	req dto.UserValidateRequestDTO,
) (dto.UserValidateResponseDTO, error) {
	user, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.UserValidateResponseDTO{}, xerror.Debug("not found username")
		}

		return dto.UserValidateResponseDTO{}, wrapNonDomainError(xerror.ServerityCritical, err)
	}

	ok, err := uc.userDomain.Validate(user.HashedPass, req.Password)
	if err != nil {
		return dto.UserValidateResponseDTO{}, wrapDomainError(err)
	}

	if !ok {
		return dto.UserValidateResponseDTO{}, xerror.Debug("wrong password")
	}

	return dto.UserValidateResponseDTO{User: user}, nil
}
