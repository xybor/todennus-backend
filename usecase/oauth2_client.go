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

type OAuth2ClientUsecase struct {
	oauth2ClientDomain abstraction.OAuth2ClientDomain

	oauth2ClientRepo abstraction.OAuth2ClientRepository
}

func NewOAuth2ClientUsecase(
	oauth2ClientDomain abstraction.OAuth2ClientDomain,
	oauth2ClientRepo abstraction.OAuth2ClientRepository,
) *OAuth2ClientUsecase {
	return &OAuth2ClientUsecase{
		oauth2ClientDomain: oauth2ClientDomain,
		oauth2ClientRepo:   oauth2ClientRepo,
	}
}

func (usecase *OAuth2ClientUsecase) CreateClient(
	ctx context.Context,
	req dto.OAuth2ClientCreateRequestDTO,
) (dto.OAuth2ClientCreateResponseDTO, error) {
	userID := xcontext.RequestUserID(ctx)
	if userID == 0 {
		return dto.OAuth2ClientCreateResponseDTO{}, xerror.WrapDebug(ErrUnauthorized)
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

func (usecase *OAuth2ClientUsecase) GetClient(
	ctx context.Context,
	req dto.OAuth2ClientGetRequestDTO,
) (dto.OAuth2ClientGetResponseDTO, error) {
	client, err := usecase.oauth2ClientRepo.GetByID(ctx, req.ClientID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return dto.OAuth2ClientGetResponseDTO{}, xerror.WrapDebug(ErrClientNotFound)
		}

		return dto.OAuth2ClientGetResponseDTO{}, wrapNonDomainError(xerror.ServerityWarn, err)
	}

	return dto.NewOAuth2ClientGetResponse(ctx, client), nil
}
