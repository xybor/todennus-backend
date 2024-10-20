package usecase

import (
	"context"
	"errors"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/usecase/abstraction"
	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/lock"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
)

type UserUsecase struct {
	adminLocker       lock.Locker
	shouldCreateAdmin bool

	userDomain abstraction.UserDomain
	userRepo   abstraction.UserRepository
}

func NewUserUsecase(
	locker lock.Locker,
	userRepo abstraction.UserRepository,
	userDomain abstraction.UserDomain,
) *UserUsecase {
	return &UserUsecase{
		adminLocker:       locker,
		shouldCreateAdmin: true,
		userRepo:          userRepo,
		userDomain:        userDomain,
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

	created, err := uc.createAdmin(ctx, &user)
	if err != nil {
		return dto.UserRegisterResponseDTO{}, err
	}

	if !created {
		if err = uc.userRepo.Create(ctx, user); err != nil {
			return dto.UserRegisterResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
		}
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

	user, err := usecase.userRepo.GetByID(ctx, req.UserID.Int64())
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

func (uc *UserUsecase) createAdmin(
	ctx context.Context,
	user *domain.User,
) (bool, error) {
	if !uc.shouldCreateAdmin {
		return false, nil
	}

	err := lock.Func(uc.adminLocker, ctx, func() error {
		adminCount, err := uc.userRepo.CountByRole(ctx, domain.UserRoleAdmin)
		if err != nil {
			return wrapNonDomainError(xerror.ServerityCritical, err)
		}

		if adminCount > 0 {
			uc.shouldCreateAdmin = false
			return nil
		}

		user.Role = domain.UserRoleAdmin
		err = uc.userRepo.Create(ctx, *user)
		if err != nil {
			return wrapNonDomainError(xerror.ServerityWarn, err)
		}

		uc.shouldCreateAdmin = false
		return nil
	})

	return !uc.shouldCreateAdmin, err
}
