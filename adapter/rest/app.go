package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	builtinMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/xybor/todennus-backend/adapter/rest/middleware"
	"github.com/xybor/todennus-backend/wiring"
	config "github.com/xybor/todennus-config"
)

func App(
	config config.Config,
	infras wiring.Infras,
	usecases wiring.Usecases,
) chi.Router {
	r := chi.NewRouter()

	r.Use(builtinMiddleware.Recoverer)
	r.Use(middleware.WithRequestID(config))
	r.Use(builtinMiddleware.RealIP)
	r.Use(middleware.WithInfras(config, infras))
	r.Use(middleware.Timer)
	r.Use(middleware.Authentication(infras.TokenEngine))
	r.Use(middleware.WithSession(infras.SessionManager))

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
