package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/x"
	"github.com/xybor/x/logging"
	"github.com/xybor/x/xcontext"
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

		req, err := x.ParseHTTPRequest[dto.OAuth2TokenRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.oauth2Usecase.Token(ctx, req.To())
		if err != nil {
			if code, errResp := dto.NewOAuth2TokenErrorResponseDTO(err); code != 0 {
				response.Write(ctx, w, code, errResp)
			} else {
				response.WriteErrorMsg(ctx, w, http.StatusInternalServerError, "Internal server error")
			}

			logging.LogError(xcontext.Logger(ctx), err)
			return
		}

		response.Write(ctx, w, http.StatusOK, dto.NewOAuth2TokenResponseDTO(resp))
	}
}
