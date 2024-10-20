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
	userID := xcontext.RequestUserID(ctx)
	if userID == 0 {
		return dto.OAuth2ClientCreateResponseDTO{}, xerror.WrapDebug(ErrUnauthorized)
	}

	requiredScope := scope.New(domain.Actions.Write.Create, domain.Resources.Client)
	if !xcontext.Scope(ctx).Contains(requiredScope) {
		return dto.OAuth2ClientCreateResponseDTO{}, xerror.WrapDebug(ErrForbidden).
			WithMessage("insufficient scope %s", requiredScope.String())
	}

	client, secret, err := usecase.oauth2ClientDomain.CreateClient(userID, req.Name, req.IsConfidential)
	if err != nil {
		return dto.OAuth2ClientCreateResponseDTO{}, wrapDomainError(err)
	}

	err = usecase.oauth2ClientRepo.Create(ctx, client)
	if err != nil {
		return dto.OAuth2ClientCreateResponseDTO{}, wrapNonDomainError(xerror.ServerityCritical, err)
	}

	return dto.NewOAuth2ClientCreateResponseDTO(ctx, client, secret), nil
}

func (usecase *OAuth2ClientUsecase) CreateByAdmin(
	ctx context.Context,
	req dto.OAuth2ClientCreateFirstRequestDTO,
) (dto.OAuth2ClientCreateByAdminResponseDTO, error) {
	if !usecase.isNoClient {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.WrapDebug(ErrRequestInvalid)
	}

	if err := usecase.firstClientLock.Lock(ctx); err != nil {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}
	defer usecase.firstClientLock.Unlock(ctx)

	count, err := usecase.oauth2ClientRepo.Count(ctx)
	if err != nil {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	if count > 0 {
		usecase.isNoClient = false
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.WrapDebug(ErrRequestInvalid)
	}

	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.WrapDebug(ErrUserNotFound)
		}

		return dto.OAuth2ClientCreateByAdminResponseDTO{}, wrapNonDomainError(xerror.ServerityCritical, err)
	}

	ok, err := usecase.userDomain.Validate(user.HashedPass, req.Password)
	if err != nil {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, wrapNonDomainError(xerror.ServerityCritical, err)
	}

	if !ok {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.WrapDebug(ErrUnauthorized).
			WithMessage("invalid password")
	}

	if user.Role != domain.UserRoleAdmin {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, xerror.WrapDebug(ErrForbidden).
			WithMessage("require admin")
	}

	client, secret, err := usecase.oauth2ClientDomain.CreateClient(user.ID, req.Name, true)
	if err != nil {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, wrapDomainError(err)
	}

	err = usecase.oauth2ClientRepo.Create(ctx, client)
	if err != nil {
		return dto.OAuth2ClientCreateByAdminResponseDTO{}, wrapNonDomainError(xerror.ServerityCritical, err)
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
			return dto.OAuth2ClientGetResponseDTO{}, xerror.WrapDebug(ErrClientNotFound)
		}

		return dto.OAuth2ClientGetResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.NewOAuth2ClientGetResponse(ctx, client), nil
}
