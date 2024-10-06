package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/todennus-backend/pkg/xcontext"
	"github.com/xybor/todennus-backend/pkg/xhttp"
)

type OAuth2Adapter struct {
	oauth2Usecase abstraction.OAuth2Usecase
}

func NewOAuth2Adapter(oauth2Usecase abstraction.OAuth2Usecase) *OAuth2Adapter {
	return &OAuth2Adapter{oauth2Usecase: oauth2Usecase}
}

func (a *OAuth2Adapter) Router(r chi.Router) {
	r.Post("/token", a.Token())
}

func (a *OAuth2Adapter) Token() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseRequest[dto.OAuth2TokenRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Token(ctx, req.To())
		if err != nil {
			if code, errResp := dto.OAuth2TokenErrorResponseFrom(err); code != 0 {
				response.Write(ctx, w, code, errResp)
			} else {
				xcontext.Logger(ctx).Warn(err.Error())
				response.WriteErrorMsg(ctx, w, http.StatusInternalServerError, "Internal server error")
			}

			return
		}

		response.Write(ctx, w, http.StatusOK, dto.OAuth2TokenResponseFrom(resp))
	}
}
