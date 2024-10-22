package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/rest/abstraction"
	"github.com/xybor/todennus-backend/adapter/rest/dto"
	"github.com/xybor/todennus-backend/adapter/rest/middleware"
	"github.com/xybor/todennus-backend/adapter/rest/response"
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
	r.Post("/validate", a.Validate())

	r.Get("/{user_id}", middleware.RequireAuthentication(a.GetByID()))
	r.Get("/username/{username}", middleware.RequireAuthentication(a.GetByUsername()))
}

func (a *UserRESTAdapter) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request, err := xhttp.ParseHTTPRequest[dto.UserRegisterRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		user, err := a.userUsecase.Register(ctx, request.To())
		response.NewResponseHandler(ctx, dto.NewUserRegisterResponseDTO, user, err).
			Map(http.StatusConflict, usecase.ErrUsernameExisted).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *UserRESTAdapter) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByIDRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		ucReq, err := req.To(xcontext.RequestUserID(ctx))
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.GetByID(ctx, ucReq)
		response.NewResponseHandler(ctx, dto.NewUserGetByIDResponseDTO, resp, err).
			Map(http.StatusNotFound, usecase.ErrUserNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *UserRESTAdapter) GetByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByUsernameRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.GetByUsername(ctx, req.To())
		response.NewResponseHandler(ctx, dto.NewUserGetByUsernameResponseDTO, resp, err).
			Map(http.StatusNotFound, usecase.ErrUserNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

func (a *UserRESTAdapter) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserValidateRequestDTO](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.ValidateCredentials(ctx, req.To())
		response.NewResponseHandler(ctx, dto.NewUserValidateResponseDTO, resp, err).
			Map(http.StatusUnauthorized, usecase.ErrCredentialsInvalid).
			WriteHTTPResponse(ctx, w)
	}
}
