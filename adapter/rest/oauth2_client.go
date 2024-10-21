package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/middleware"
	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/xhttp"
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
	r.Get("/{client_id}", middleware.RequireAuthentication(a.Get()))

	r.Post("/", middleware.RequireAuthentication(a.Create()))
	r.Post("/first", a.CreateByAdmin())
}

func (a *OAuth2ClientAdapter) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2ClientGetRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2ClientUsecase.Get(ctx, req.To())
		response.NewResponseHandler(dto.NewOAuth2ClientGetResponseDTO, resp, err).
			Map(http.StatusNotFound, usecase.ErrClientInvalid).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *OAuth2ClientAdapter) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2ClientCreateRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2ClientUsecase.Create(ctx, req.To())
		response.NewResponseHandler(dto.NewOauth2ClientCreateResponseDTO, resp, err).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *OAuth2ClientAdapter) CreateByAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2ClientCreateFirstRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2ClientUsecase.CreateByAdmin(ctx, req.To())
		response.NewResponseHandler(dto.NewOauth2ClientCreateFirstResponseDTO, resp, err).
			Map(http.StatusBadRequest, usecase.ErrUserNotFound).
			WriteHTTPResponse(ctx, w)
	}
}
