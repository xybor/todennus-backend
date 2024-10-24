package grpc

import (
	api "github.com/xybor/todennus-backend/adapter/grpc/gen"
	"github.com/xybor/todennus-backend/adapter/grpc/interceptor"
	"github.com/xybor/todennus-backend/wiring"
	config "github.com/xybor/todennus-config"
	"google.golang.org/grpc"
)

func App(config *config.Config, infras *wiring.Infras, usecases *wiring.Usecases) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnarySetupContext(config, infras)),
		grpc.UnaryInterceptor(interceptor.UnaryAuthenticate(infras.TokenEngine)),
	)

	api.RegisterUserServer(s, NewUserServer(usecases.UserUsecase))

	return s
}
