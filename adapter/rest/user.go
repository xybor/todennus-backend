package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/response"
	"github.com/xybor/todennus-backend/domain"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xhttp"
)

type UserRESTAdapter struct {
	userUsecase abstraction.UserUsecase
}

func NewUserAdapter(userUsecase abstraction.UserUsecase) *UserRESTAdapter {
	return &UserRESTAdapter{userUsecase: userUsecase}
}

func (a *UserRESTAdapter) Router(r chi.Router) {
	r.Post("/", a.Register())

	r.Get("/{user_id}", a.GetByID())
	r.Get("/username/{username}", a.GetByUsername())
}

func (a *UserRESTAdapter) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request, err := xhttp.ParseHTTPRequest[dto.UserRegisterRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		user, err := a.userUsecase.Register(ctx, request.To())
		response.NewResponseHandler(dto.NewUserRegisterResponseDTO(user), err).
			Map(http.StatusConflict, usecase.ErrUsernameExisted).
			Map(http.StatusBadRequest, domain.ErrUsernameInvalid, domain.ErrPasswordInvalid).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *UserRESTAdapter) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByIDRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		ucReq, err := req.To(ctx)
		if err != nil {
			xcontext.Logger(ctx).Debug("failed-to-convert-req", "err", err)
			response.WriteError(ctx, w, http.StatusBadRequest, err)
			return
		}

		resp, err := a.userUsecase.GetByID(ctx, ucReq)
		response.NewResponseHandler(dto.NewUserGetByIDResponseDTO(resp), err).
			Map(http.StatusNotFound, usecase.ErrUserNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *UserRESTAdapter) GetByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByUsernameRequestDTO](r)
		if err != nil {
			response.HandleParseError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.GetByUsername(ctx, req.To())
		response.NewResponseHandler(dto.NewUserGetByUsernameResponseDTO(resp), err).
			Map(http.StatusNotFound, usecase.ErrUserNotFound).
			WriteHTTPResponse(ctx, w)
	}
}
