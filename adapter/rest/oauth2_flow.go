package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/xhttp"
)

type OAuth2Adapter struct {
	oauth2Usecase abstraction.OAuth2Usecase
}

func NewOAuth2Adapter(oauth2Usecase abstraction.OAuth2Usecase) *OAuth2Adapter {
	return &OAuth2Adapter{oauth2Usecase: oauth2Usecase}
}

func (a *OAuth2Adapter) OAuth2Router(r chi.Router) {
	r.Get("/authorize", a.Authorize())
	r.Post("/token", a.Token())
}

func (a *OAuth2Adapter) Authorize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2AuthorizeRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Authorize(ctx, req.To())
		if err != nil {
			if url, err := dto.NewOAuth2AuthorizeRedirectURIWithError(ctx, &req, err); err != nil {
				response.HandleError(ctx, w, err)
			} else {
				response.Redirect(ctx, w, r, url, http.StatusSeeOther)
			}

			return
		}

		redirectURI, err := dto.NewOAuth2AuthorizeRedirectURI(&req, resp)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		response.Redirect(ctx, w, r, redirectURI, http.StatusSeeOther)
	}
}

func (a *OAuth2Adapter) Token() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2TokenRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Token(ctx, req.To())
		response.NewResponseHandler(dto.NewOAuth2TokenResponseDTO, resp, err).
			Map(http.StatusBadRequest,
				usecase.ErrRequestInvalid, usecase.ErrClientInvalid,
				usecase.ErrScopeInvalid, usecase.ErrTokenInvalidGrant,
			).
			WriteHTTPResponseWithoutWrap(ctx, w)
	}
}

func (a *OAuth2Adapter) AuthenticationCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2AuthenticationCallbackRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		usecaseReq, err := req.To()
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.AuthenticationCallback(ctx, usecaseReq)
		response.NewResponseHandler(dto.NewOAuth2AuthenticationCallbackResponseDTO, resp, err).
			Map(http.StatusUnauthorized, usecase.ErrIdPInvalid).
			Map(http.StatusBadRequest, usecase.ErrUserNotFound, usecase.ErrRequestInvalid).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *OAuth2Adapter) SessionUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2SessionUpdateRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.SessionUpdate(ctx, req.To())
		response.NewResponseHandler(dto.NewOAuth2SessionUpdateRedirectURI, resp, err).
			Redirect(ctx, w, r, http.StatusSeeOther)
	}
}
