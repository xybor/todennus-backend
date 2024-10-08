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

type UserRESTAdapter struct {
	userUsecase abstraction.UserUsecase
}

func NewUserAdapter(userUsecase abstraction.UserUsecase) *UserRESTAdapter {
	return &UserRESTAdapter{userUsecase: userUsecase}
}

func (a *UserRESTAdapter) Router(r chi.Router) {
	r.Post("/", a.Register())
}

func (a *UserRESTAdapter) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request, err := xhttp.ParseRequest[dto.UserRegisterRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		user, err := a.userUsecase.Register(ctx, request.To())
		response.NewResponseHandler(dto.NewUserRegisterResponseDTO(user), err).
			Map(http.StatusConflict, usecase.ErrUsernameExisted).
			Map(http.StatusBadRequest, domain.ErrUsernameInvalid, domain.ErrPasswordInvalid).
			Map(http.StatusInternalServerError).
			WriteHTTPResponse(ctx, w)
	}
}
