package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/pkg/xhttp"
	"github.com/xybor/todennus-backend/usecase"
)

type OAuth2ClientAdapter struct {
	oauth2ClientUsecase abstraction.OAuth2ClientUsecase
}

func NewOAuth2ClientAdapter(oauth2ClientUsecase abstraction.OAuth2ClientUsecase) *OAuth2ClientAdapter {
	return &OAuth2ClientAdapter{
		oauth2ClientUsecase: oauth2ClientUsecase,
	}
}

func (a *OAuth2ClientAdapter) Router(r chi.Router) {
	r.Post("/", a.Create())
	r.Get("/{client_id}", a.Get())
}

func (a *OAuth2ClientAdapter) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseRequest[dto.OAuth2ClientGetRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.oauth2ClientUsecase.GetClient(ctx, req.To())
		response.NewResponseHandler(dto.NewOAuth2ClientGetResponseDTO(resp), err).
			Map(http.StatusNotFound, usecase.ErrClientNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *OAuth2ClientAdapter) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseRequest[dto.OAuth2ClientCreateRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.oauth2ClientUsecase.CreateClient(ctx, req.To())
		response.NewResponseHandler(dto.NewOauth2ClientCreateResponseDTO(resp), err).
			Map(http.StatusBadRequest, domain.ErrClientNameInvalid).
			WriteHTTPResponse(ctx, w)
	}
}
