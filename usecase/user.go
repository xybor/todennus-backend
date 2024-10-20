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
		return dto.UserRegisterResponseDTO{}, xerror.Wrap(ErrUsernameExisted,
			"username %s has already taken before", req.Username)
	}

	if !errors.Is(err, database.ErrRecordNotFound) {
		xcontext.Logger(ctx).Critical("failed-to-get-user", "err", err, "username", req.Username)
		return dto.UserRegisterResponseDTO{}, ErrServer
	}

	user, err := uc.userDomain.Create(req.Username, req.Password)
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-new-user", "err", err)
		return dto.UserRegisterResponseDTO{}, ErrServer
	}

	created, err := uc.createAdmin(ctx, &user)
	if err != nil {
		return dto.UserRegisterResponseDTO{}, err
	}

	if !created {
		if err = uc.userRepo.Create(ctx, user); err != nil {
			xcontext.Logger(ctx).Warn("failed-to-create-user", "err", err)
			return dto.UserRegisterResponseDTO{}, ErrServer
		}
	}

	return dto.NewUserRegisterResponseDTO(ctx, user), nil
}

func (usecase *UserUsecase) GetByID(
	ctx context.Context,
	req dto.UserGetByIDRequestDTO,
) (dto.UserGetByIDResponseDTO, error) {
	user, err := usecase.userRepo.GetByID(ctx, req.UserID.Int64())
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.UserGetByIDResponseDTO{}, xerror.Wrap(ErrUserNotFound,
				"not found user with id %d", req.UserID)
		}

		xcontext.Logger(ctx).Warn("failed-to-get-user", "err", err, "uid", req.UserID)
		return dto.UserGetByIDResponseDTO{}, ErrServer
	}

	return dto.NewUserGetByIDResponseDTO(ctx, user), nil
}

func (usecase *UserUsecase) GetByUsername(
	ctx context.Context,
	req dto.UserGetByUsernameRequestDTO,
) (dto.UserGetByUsernameResponseDTO, error) {
	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.UserGetByUsernameResponseDTO{}, xerror.Wrap(ErrUserNotFound,
				"not found user with username %s", req.Username)
		}

		xcontext.Logger(ctx).Warn("failed-to-get-user", "err", err, "username", req.Username)
		return dto.UserGetByUsernameResponseDTO{}, ErrServer
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

	xcontext.Logger(ctx).Info("check-create-admin")

	err := lock.Func(uc.adminLocker, ctx, func() error {
		adminCount, err := uc.userRepo.CountByRole(ctx, domain.UserRoleAdmin)
		if err != nil {
			xcontext.Logger(ctx).Warn("failed-to-count-by-admin-role", "err", err)
			return ErrServer
		}

		if adminCount > 0 {
			xcontext.Logger(ctx).Info("cancel-create-admin")
			uc.shouldCreateAdmin = false
			return nil
		}

		xcontext.Logger(ctx).Info("create-admin", "username", user.Username, "uid", user.ID)

		user.Role = domain.UserRoleAdmin
		err = uc.userRepo.Create(ctx, *user)
		if err != nil {
			xcontext.Logger(ctx).Warn("failed-to-create-admin-user", "err", err)
			return ErrServer
		}

		uc.shouldCreateAdmin = false
		return nil
	})

	return uc.shouldCreateAdmin, err
}
