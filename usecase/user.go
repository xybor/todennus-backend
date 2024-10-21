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
	req *dto.UserRegisterRequestDTO,
) (*dto.UserRegisterResponseDTO, error) {
	_, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, xerror.Enrich(ErrUsernameExisted, "username %s has already taken before", req.Username)
	}

	if !errors.Is(err, database.ErrRecordNotFound) {
		return nil, ErrServer.Hide(err, "failed-to-get-user")
	}

	user, err := uc.userDomain.Create(req.Username, req.Password)
	if err != nil {
		return nil, errcfg.Event(err, "failed-to-new-user").Enrich(ErrRequestInvalid).Error()
	}

	shouldCreateUser, err := uc.createAdmin(ctx, user)
	if err != nil {
		return nil, err
	}

	if shouldCreateUser {
		if err = uc.userRepo.Create(ctx, user); err != nil {
			return nil, ErrServer.Hide(err, "failed-to-create-user")
		}
	}

	return dto.NewUserRegisterResponseDTO(ctx, user), nil
}

func (usecase *UserUsecase) GetByID(
	ctx context.Context,
	req *dto.UserGetByIDRequestDTO,
) (*dto.UserGetByIDResponseDTO, error) {
	user, err := usecase.userRepo.GetByID(ctx, req.UserID.Int64())
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, xerror.Enrich(ErrUserNotFound, "not found user with id %d", req.UserID)
		}

		return nil, ErrServer.Hide(err, "failed-to-get-user", "uid", req.UserID)
	}

	return dto.NewUserGetByIDResponseDTO(ctx, user), nil
}

func (usecase *UserUsecase) GetByUsername(
	ctx context.Context,
	req *dto.UserGetByUsernameRequestDTO,
) (*dto.UserGetByUsernameResponseDTO, error) {
	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, xerror.Enrich(ErrUserNotFound, "not found user with username %s", req.Username)
		}

		return nil, ErrServer.Hide(err, "failed-to-get-user", "username", req.Username)
	}

	return dto.NewUserGetByUsernameResponseDTO(ctx, user), nil
}

func (usecase *UserUsecase) ValidateCredentials(
	ctx context.Context,
	req *dto.UserValidateCredentialsRequestDTO,
) (*dto.UserValidateCredentialsResponseDTO, error) {
	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, xerror.Enrich(ErrCredentialsInvalid, "invalid username or password")
		}

		return nil, ErrServer.Hide(err, "failed-to-get-user", "username", req.Username)
	}

	if err := usecase.userDomain.Validate(user.HashedPass, req.Password); err != nil {
		return nil, errcfg.Event(err, "failed-to-validate-user-credentials").
			EnrichWith(ErrCredentialsInvalid, "invalid username or password").
			Error()
	}

	ctx = xcontext.WithRequestUserID(ctx, user.ID)
	return dto.NewUserValidateCredentialsResponseDTO(ctx, user), nil
}

func (uc *UserUsecase) createAdmin(
	ctx context.Context,
	user *domain.User,
) (bool, error) {
	if !uc.shouldCreateAdmin {
		return true, nil
	}

	shouldCreateUser := true
	xcontext.Logger(ctx).Info("check-create-admin")

	err := lock.Func(uc.adminLocker, ctx, func() error {
		adminCount, err := uc.userRepo.CountByRole(ctx, domain.UserRoleAdmin)
		if err != nil {
			return ErrServer.Hide(err, "failed-to-count-by-admin-role")
		}

		if adminCount > 0 {
			xcontext.Logger(ctx).Info("cancel-create-admin")
			uc.shouldCreateAdmin = false
			return nil
		}

		xcontext.Logger(ctx).Info("create-admin", "username", user.Username, "uid", user.ID)

		user.Role = domain.UserRoleAdmin
		err = uc.userRepo.Create(ctx, user)
		if err != nil {
			return ErrServer.Hide(err, "failed-to-create-admin-user")
		}

		shouldCreateUser = false
		uc.shouldCreateAdmin = false
		return nil
	})

	return shouldCreateUser, err
}
