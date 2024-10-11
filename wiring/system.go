package wiring

import (
	"context"
	"fmt"

	"github.com/xybor/todennus-backend/config"
)

type System struct {
	config.Config
	Infras
	Domains
	Databases
	Repositories
	Usecases
}

func InitializeSystem(envPaths []string, iniPaths []string) (System, context.Context, error) {
	config, err := config.Load(sources(envPaths), sources(iniPaths))
	if err != nil {
		return System{}, nil, fmt.Errorf("cannot load variable and secrets, err=%w", err)
	}

	infras, err := InitializeInfras(config)
	if err != nil {
		return System{}, nil, fmt.Errorf("cannot initialize context, err=%w", err)
	}

	ctx := context.Background()
	ctx = WithInfras(ctx, infras)

	domains, err := InitializeDomains(ctx, config, infras)
	if err != nil {
		return System{}, nil, fmt.Errorf("cannot initialize domains, err=%w", err)
	}

	databases, err := InitializeDatabases(ctx, config)
	if err != nil {
		return System{}, nil, fmt.Errorf("cannot initialize databases, err=%w", err)
	}

	repositories, err := InitializeRepositories(ctx, databases)
	if err != nil {
		return System{}, nil, fmt.Errorf("cannot initialize repositories, err=%w", err)
	}

	usecases, err := InitializeUsecases(ctx, infras, domains, repositories)
	if err != nil {
		return System{}, nil, fmt.Errorf("cannot initialize usecases, err=%w", err)
	}

	return System{
		Config:       config,
		Infras:       infras,
		Databases:    databases,
		Repositories: repositories,
		Domains:      domains,
		Usecases:     usecases,
	}, ctx, nil
}

func sources(paths []string) []string {
	sources := []string{}
	for i := range paths {
		if len(paths[i]) > 0 {
			sources = append(sources, paths[i])
		}
	}

	return sources
}
