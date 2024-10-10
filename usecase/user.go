package usecase

import (
	"context"
	"errors"

	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/pkg/xcontext"
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

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return dto.UserRegisterResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.NewUserRegisterResponseDTO(ctx, user), nil
}

func (usecase *UserUsecase) GetByID(
	ctx context.Context,
	req dto.UserGetByIDRequestDTO,
) (dto.UserGetByIDResponseDTO, error) {
	if xcontext.RequestUserID(ctx) == 0 {
		return dto.UserGetByIDResponseDTO{}, xerror.WrapDebug(ErrUnauthorized)
	}

	user, err := usecase.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.UserGetByIDResponseDTO{}, xerror.WrapDebug(ErrUserNotFound).
				WithMessage("not found user with id %d", req.UserID)
		}

		return dto.UserGetByIDResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.NewUserGetByIDResponseDTO(ctx, user), nil
}

func (usecase *UserUsecase) GetByUsername(
	ctx context.Context,
	req dto.UserGetByUsernameRequestDTO,
) (dto.UserGetByUsernameResponseDTO, error) {
	if xcontext.RequestUserID(ctx) == 0 {
		return dto.UserGetByUsernameResponseDTO{}, xerror.WrapDebug(ErrUnauthorized)
	}

	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.UserGetByUsernameResponseDTO{}, xerror.WrapDebug(ErrUserNotFound).
				WithMessage("not found user %s", req.Username)
		}

		return dto.UserGetByUsernameResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.NewUserGetByUsernameResponseDTO(ctx, user), nil
}
