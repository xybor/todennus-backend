package usecase

import (
	"context"
	"errors"

	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/infras/database"
	"github.com/xybor/todennus-backend/usecase/abstraction"
	"github.com/xybor/todennus-backend/usecase/dto"
	"github.com/xybor/x/lock"
	"github.com/xybor/x/scope"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
)

type OAuth2ClientUsecase struct {
	isNoClient         bool
	firstClientLock    lock.Locker
	userDomain         abstraction.UserDomain
	oauth2ClientDomain abstraction.OAuth2ClientDomain

	userRepo         abstraction.UserRepository
	oauth2ClientRepo abstraction.OAuth2ClientRepository
}

func NewOAuth2ClientUsecase(
	locker lock.Locker,
	userDomain abstraction.UserDomain,
	oauth2ClientDomain abstraction.OAuth2ClientDomain,
	userRepo abstraction.UserRepository,
	oauth2ClientRepo abstraction.OAuth2ClientRepository,
) *OAuth2ClientUsecase {
	return &OAuth2ClientUsecase{
		isNoClient:         true,
		firstClientLock:    locker,
		userDomain:         userDomain,
		oauth2ClientDomain: oauth2ClientDomain,
		userRepo:           userRepo,
		oauth2ClientRepo:   oauth2ClientRepo,
	}
}

func (usecase *OAuth2ClientUsecase) Create(
	ctx context.Context,
	req dto.OAuth2ClientCreateRequestDTO,
) (dto.OAuth2ClientCreateResponseDTO, error) {
	requiredScope := scope.New(domain.Actions.Write.Create, domain.Resources.Client)
	if !xcontext.Scope(ctx).Contains(requiredScope) {
		return dto.OAuth2ClientCreateResponseDTO{}, xerror.Wrap(ErrForbidden,
			"insufficient scope, require %s", requiredScope.String())
	}

	userID := xcontext.RequestUserID(ctx)
	client, secret, err := usecase.oauth2ClientDomain.CreateClient(userID, req.Name, req.IsConfidential)
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-new-client", "err", err)
		return dto.OAuth2ClientCreateResponseDTO{}, ErrServer
	}

	if err = usecase.oauth2ClientRepo.Create(ctx, client); err != nil {
		xcontext.Logger(ctx).Warn("failed-to-create-client", "err", err)
		return dto.OAuth2ClientCreateResponseDTO{}, ErrServer
	}

	return dto.NewOAuth2ClientCreateResponseDTO(ctx, client, secret), nil
}

func (usecase *OAuth2ClientUsecase) CreateByAdmin(
	ctx context.Context,
	req dto.OAuth2ClientCreateFirstRequestDTO,
) (dto.OAuth2ClientCreateByAdminResponseDTO, error) {
	if !usecase.isNoClient {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.Wrap(ErrRequestInvalid,
			"this api is only openned for creating the first client")
	}

	if err := usecase.firstClientLock.Lock(ctx); err != nil {
		xcontext.Logger(ctx).Warn("failed-to-lock-first-client-flow", "err", err)
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, ErrServer
	}
	defer usecase.firstClientLock.Unlock(ctx)

	count, err := usecase.oauth2ClientRepo.Count(ctx)
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-count-client", "err", err)
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, ErrServer
	}

	if count > 0 {
		usecase.isNoClient = false
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.Wrap(ErrRequestInvalid,
			"this api is only openned for creating the first client")
	}

	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.Wrap(ErrForbidden,
				"invalid username or password")
		}

		xcontext.Logger(ctx).Warn("failed-to-get-user", "err", err, "username", req.Username)
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, ErrServer
	}

	ok, err := usecase.userDomain.Validate(user.HashedPass, req.Password)
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-validate-user-credentials", "err", err)
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, ErrServer
	}

	if !ok {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.Wrap(ErrForbidden,
			"invalid username or password")
	}

	if user.Role != domain.UserRoleAdmin {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.Wrap(ErrForbidden, "require admin")
	}

	client, secret, err := usecase.oauth2ClientDomain.CreateClient(user.ID, req.Name, true)
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-new-first-client", "err", err)
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, ErrServer
	}

	if err = usecase.oauth2ClientRepo.Create(ctx, client); err != nil {
		xcontext.Logger(ctx).Warn("failed-to-create-first-client", "err", err)
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, ErrServer
	}

	usecase.isNoClient = false
	return dto.NewOAuth2ClientCreateFirstResponseDTO(ctx, client, secret), nil
}

func (usecase *OAuth2ClientUsecase) Get(
	ctx context.Context,
	req dto.OAuth2ClientGetRequestDTO,
) (dto.OAuth2ClientGetResponseDTO, error) {
	client, err := usecase.oauth2ClientRepo.GetByID(ctx, req.ClientID.Int64())
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2ClientGetResponseDTO{}, xerror.Wrap(ErrClientInvalid, "not found client")
		}

		xcontext.Logger(ctx).Warn("failed-to-get-client", "err", err, "cid", req.ClientID)
		return dto.OAuth2ClientGetResponseDTO{}, ErrServer
	}

	return dto.NewOAuth2ClientGetResponse(ctx, client), nil
}
