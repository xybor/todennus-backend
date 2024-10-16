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
	r.Use(builtinMiddleware.RequestID)
	r.Use(builtinMiddleware.RealIP)
	r.Use(middleware.WithInfras(infras))
	r.Use(middleware.RoundTripTime)
	r.Use(middleware.Authentication(infras.TokenEngine))

	r.Route("/users", NewUserAdapter(usecases.UserUsecase).Router)
	r.Route("/oauth2", NewOAuth2Adapter(usecases.OAuth2Usecase).Router)
	r.Route("/oauth2_clients", NewOAuth2ClientAdapter(usecases.OAuth2ClientUsecase).Router)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
