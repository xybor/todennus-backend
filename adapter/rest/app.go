package rest

import (
	"net/http"

	_ "github.com/xybor/todennus-backend/docs"

	"github.com/go-chi/chi/v5"
	builtinMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/xybor/todennus-backend/adapter/rest/middleware"
	"github.com/xybor/todennus-backend/wiring"
	config "github.com/xybor/todennus-config"
)

// @title Todennus API Endpoints
// @version 1.0
// @description This is Todennus - An Open ID Connect and OAuth2 Provider
func App(
	config *config.Config,
	infras *wiring.Infras,
	usecases *wiring.Usecases,
) chi.Router {
	r := chi.NewRouter()

	r.Use(builtinMiddleware.Recoverer)
	r.Use(builtinMiddleware.RealIP)
	r.Use(middleware.WithRequestID())
	r.Use(middleware.WithInfras(infras))
	r.Use(middleware.Timer(config))
	r.Use(middleware.Timeout(config))
	r.Use(middleware.Authentication(infras.TokenEngine))
	r.Use(middleware.WithSession(infras.SessionManager))

	r.Get("/specs/*", httpSwagger.WrapHandler)

	userAdapter := NewUserAdapter(usecases.UserUsecase)
	oauth2FlowAdapter := NewOAuth2Adapter(usecases.OAuth2Usecase)
	oauth2ClientAdapter := NewOAuth2ClientAdapter(usecases.OAuth2ClientUsecase)

	r.Get("/session/update", oauth2FlowAdapter.SessionUpdate())
	r.Post("/auth/callback", oauth2FlowAdapter.AuthenticationCallback())

	r.Route("/users", userAdapter.Router)
	r.Route("/oauth2", oauth2FlowAdapter.OAuth2Router)
	r.Route("/oauth2_clients", oauth2ClientAdapter.Router)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })

	return r
}
