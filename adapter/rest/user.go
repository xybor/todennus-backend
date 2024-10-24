package rest

import (
	"net/http"

	_ "github.com/xybor/todennus-backend/adapter/rest/standard"

	"github.com/go-chi/chi/v5"
	"github.com/xybor/todennus-backend/adapter/abstraction"
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

// @Summary Register a new user
// @Description Register a new user by providing username and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "User registration data"
// @Success 201 {object} standard.SwaggerSuccessResponse[dto.UserRegisterResponse] "User registered successfully"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 409 {object} standard.SwaggerDuplicatedErrorResponse "Duplicated"
// @Router /users [post]
func (a *UserRESTAdapter) Register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request, err := xhttp.ParseHTTPRequest[dto.UserRegisterRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		user, err := a.userUsecase.Register(ctx, request.To())
		response.NewResponseHandler(ctx, dto.NewUserRegisterResponse(user), err).
			WithDefaultCode(http.StatusCreated).
			Map(http.StatusConflict, usecase.ErrDuplicated).
			Map(http.StatusBadRequest, usecase.ErrRequestInvalid).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Get user by id
// @Description Get an user information by user id. <br>
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} standard.SwaggerSuccessResponse[dto.UserGetByIDResponse] "Get user successfully"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 401 {object} standard.SwaggerUnauthorizedErrorResponse "Invalid credentials"
// @Failure 404 {object} standard.SwaggerNotFoundErrorResponse "Not found"
// @Router /users/{user_id} [get]
func (a *UserRESTAdapter) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByIDRequest](r)
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
		response.NewResponseHandler(ctx, dto.NewUserGetByIDResponse(resp), err).
			Map(http.StatusBadRequest, usecase.ErrRequestInvalid).
			Map(http.StatusUnauthorized, usecase.ErrCredentialsInvalid).
			Map(http.StatusNotFound, usecase.ErrNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Get user by username
// @Description Get an user information by user username. <br>
// @Tags User
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} standard.SwaggerSuccessResponse[dto.UserGetByUsernameResponse] "Get user successfully"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 404 {object} standard.SwaggerNotFoundErrorResponse "Not found"
// @Router /users/username/{username} [get]
func (a *UserRESTAdapter) GetByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByUsernameRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.GetByUsername(ctx, req.To())
		response.NewResponseHandler(ctx, dto.NewUserGetByUsernameResponse(resp), err).
			Map(http.StatusBadRequest, usecase.ErrRequestInvalid).
			Map(http.StatusNotFound, usecase.ErrNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Validate user credentials
// @Description Validate the user credentials and returns the user information.
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserValidateRequest true "Validation data"
// @Success 200 {object} standard.SwaggerSuccessResponse[dto.UserValidateResponse] "Validate successfully"
// @Failure 400 {object} standard.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 401 {object} standard.SwaggerInvalidCredentialsErrorResponse "Invalid credentials"
// @Router /users/validate [post]
func (a *UserRESTAdapter) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserValidateRequest](r)
		if err != nil {
			response.HandleError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.ValidateCredentials(ctx, req.To())
		response.NewResponseHandler(ctx, dto.NewUserValidateResponse(resp), err).
			Map(http.StatusBadRequest, usecase.ErrRequestInvalid).
			Map(http.StatusUnauthorized, usecase.ErrCredentialsInvalid).
			WriteHTTPResponse(ctx, w)
	}
}
