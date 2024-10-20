package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/logging"
	"github.com/xybor/x/xcontext"
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
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Authorize(ctx, req.To())
		if err != nil {
			if url, err := dto.NewOAuth2AuthorizeRedirectURIWithError(ctx, req, err); err != nil {
				response.WriteAndWarnError(ctx, w, http.StatusInternalServerError, err)
			} else {
				response.Redirect(ctx, w, r, url, http.StatusSeeOther)
			}

			return
		}

		redirectURI, err := dto.NewOAuth2AuthorizeRedirectURI(req, resp)
		if err != nil {
			response.WriteAndWarnError(ctx, w, http.StatusInternalServerError, err)
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
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Token(ctx, req.To())
		if err != nil {
			if code, errResp := dto.NewOAuth2TokenErrorResponseDTO(err); code != 0 {
				response.Write(ctx, w, code, errResp)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			logging.LogError(xcontext.Logger(ctx), err)
			return
		}

		response.Write(ctx, w, http.StatusOK, dto.NewOAuth2TokenResponseDTO(resp))
	}
}

func (a *OAuth2Adapter) AuthenticationCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2AuthenticationCallbackRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		usecaseReq, err := req.To()
		if err != nil {
			xcontext.Logger(ctx).Debug("failed-to-parse-req", "err", err)
			response.WriteError(ctx, w, http.StatusBadRequest, err)
			return
		}

		resp, err := a.oauth2Usecase.AuthenticationCallback(ctx, usecaseReq)
		response.NewResponseHandler(dto.NewOAuth2AuthenticationCallbackResponseDTO(resp), err).
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
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.SessionUpdate(ctx, req.To())
		if err != nil {
			response.WriteAndWarnError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		response.Redirect(ctx, w, r, dto.NewOAuth2SessionUpdateRedirectURI(resp), http.StatusSeeOther)
	}
}
