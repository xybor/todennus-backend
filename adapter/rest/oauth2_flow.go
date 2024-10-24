package rest

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/abstraction"
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

	r.Get("/consent", a.GetConsentPage())
	r.Post("/consent", a.UpdateConsent())
}

// @Summary OAuth2 Authorization Endpoint
// @Description The authorization endpoint is used to interact with the resource owner and obtain an authorization grant.
// @Description This is the entry point for starting an OAuth2 flow, such as Authorization Code or Implicit.
// @Tags OAuth2
// @Param response_type query string true "The type of response requested, typically 'code' or 'token'."
// @Param client_id query string true "The client ID of the application making the authorization request."
// @Param redirect_uri query string true "The URI to which the response will be sent after the authorization."
// @Param scope query string false "The scope of the access request. It defines the level of access the application is requesting."
// @Param state query string false "An opaque value used by the client to maintain state between the request and callback."
// @Success 303 "Redirect to client application with authorization code or error"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Router /oauth2/authorize [get]
func (a *OAuth2Adapter) Authorize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2AuthorizeRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Authorize(ctx, req.To())
		if err != nil {
			if url, err := dto.NewOAuth2AuthorizeRedirectURIWithError(ctx, req, err); err != nil {
				response.HandleError(ctx, w, err)
			} else {
				response.Redirect(ctx, w, r, url, http.StatusSeeOther)
			}

			return
		}

		redirectURI, err := dto.NewOAuth2AuthorizeRedirectURI(req, resp)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		response.Redirect(ctx, w, r, redirectURI, http.StatusSeeOther)
	}
}

// @Summary OAuth2 Token Endpoint
// @Description The token endpoint is used to exchange an authorization code, client credentials, or refresh token for an access token and optionally a refresh token.
// @Description This is part of the OAuth2 flow to grant access tokens to clients.
// @Tags OAuth2
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "The OAuth2 grant type (authorization_code, client_credentials, refresh_token)"
// @Param code formData string false "The authorization code received from the authorize endpoint (required for authorization_code grant type)"
// @Param redirect_uri formData string false "The redirect URI used in the authorization request (required for authorization_code grant type)"
// @Param client_id formData string true "The client ID of the application"
// @Param client_secret formData string true "The client secret of the application"
// @Param refresh_token formData string false "The refresh token (required for refresh_token grant type)"
// @Param scope formData string false "The scope of the access request (optional, space-separated)"
// @Success 200 {object} dto.OAuth2TokenResponse "Successfully generated access token"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Router /oauth2/token [post]
func (a *OAuth2Adapter) Token() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2TokenRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Token(ctx, req.To())
		response.NewResponseHandler(ctx, dto.NewOAuth2TokenResponse(resp), err).
			Map(http.StatusBadRequest,
				usecase.ErrRequestInvalid, usecase.ErrClientInvalid,
				usecase.ErrScopeInvalid, usecase.ErrTokenInvalidGrant,
			).
			WriteHTTPResponseWithoutWrap(ctx, w)
	}
}

// @Summary Authentication Callback Endpoint
// @Description This endpoint is called by the IdP after it validated the user.
// @Description It notifies to the server about the authentication result (success or failure) and the inforamtion of user.
// @Tags OAuth2
// @Accept json
// @Produce json
// @Param body body dto.OAuth2AuthenticationCallbackRequest true "Authentication result"
// @Success 200 {object} dto.OAuth2AuthenticationCallbackResponse "Successfully accept the result"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 401 {object} standard.SwaggerUnauthorizedErrorResponse "Unauthorized"
// @Failure 401 {object} standard.SwaggerNotFoundErrorResponse "Not found"
// @Router /auth/callback [post]
func (a *OAuth2Adapter) AuthenticationCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2AuthenticationCallbackRequest](r)
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
		response.NewResponseHandler(ctx, dto.NewOAuth2AuthenticationCallbackResponse(resp), err).
			Map(http.StatusBadRequest, usecase.ErrRequestInvalid).
			Map(http.StatusNotFound, usecase.ErrNotFound).
			Map(http.StatusUnauthorized, usecase.ErrUnauthenticated).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Session Update Endpoint
// @Description The user will be redirected to this endpoint by the IdP after it sends the authentication result to the server. <br>
// @Description This endpoint updates the user session state to `authenticated`, `unauthenticated`, or `failed authentication`.
// @Tags OAuth2
// @Param authentication_id query string true "Authentication id"
// @Success 303 "Redirect back to oauth2 authorization endpoint"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Router /session/update [get]
func (a *OAuth2Adapter) SessionUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2SessionUpdateRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.SessionUpdate(ctx, req.To())
		response.NewResponseHandler(ctx, dto.NewOAuth2SessionUpdateRedirectURI(resp), err).
			Map(http.StatusBadRequest, usecase.ErrRequestInvalid).
			Redirect(ctx, w, r, http.StatusSeeOther)
	}
}

// @Summary Consent page
// @Description This endpoint serves a consent page when the server needs the user consent for client.
// @Tags OAuth2
// @Produce text/html
// @Param authorization_id query string true "Authorization ID"
// @Success 200 {string} string "Consent page rendered successfully"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Router /oauth2/consent [get]
func (a *OAuth2Adapter) GetConsentPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2GetConsentPageRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.GetConsent(ctx, req.To())
		if err != nil {
			if errors.Is(err, usecase.ErrRequestInvalid) {
				response.WriteError(ctx, w, http.StatusBadRequest, err)
			} else {
				response.WriteError(ctx, w, http.StatusInternalServerError, err)
			}
			return
		}

		tmpl, err := template.ParseFiles("template/consent.html")
		if err != nil {
			response.WriteError(ctx, w, http.StatusInternalServerError,
				usecase.ErrServer.Hide(err, "failed-to-parse-template"))
			return
		}

		if err = tmpl.Execute(w, dto.NewOAuth2GetConsentPageResponse(resp)); err != nil {
			response.WriteError(ctx, w, http.StatusInternalServerError,
				usecase.ErrServer.Hide(err, "failed-to-render-template"))
		}
	}
}

// @Summary Update consent
// @Description This endpoint updates the consent result of user then redirect back to the oauth2 authorization endpoint.
// @Tags OAuth2
// @Param authorization_id query string true "Authorization ID"
// @Param consent formData string false "The consent result (accepted or denied)"
// @Param scope formData string false "The accepted scopes of user (usually less than the requested scope)."
// @Success 303 "Redirect back to oauth2 authorization endpoint"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Router /oauth2/consent [post]
func (a *OAuth2Adapter) UpdateConsent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.OAuth2UpdateConsentRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.UpdateConsent(ctx, req.To())
		response.NewResponseHandler(ctx, dto.NewOAuth2ConsentUpdateRedirectURI(resp), err).
			Map(http.StatusBadRequest, usecase.ErrRequestInvalid, usecase.ErrAuthorizationAccessDenied).
			Redirect(ctx, w, r, http.StatusSeeOther)
	}
}
