package wiring

import (
	"context"
	"fmt"

	config "github.com/xybor/todennus-config"
)

type System struct {
	config.Config
	Infras
	Domains
	Databases
	Repositories
	Usecases
}

func InitializeSystem(paths ...string) (System, context.Context, error) {
	config, err := config.Load(sources(paths)...)
	if err != nil {
		return System{}, nil, fmt.Errorf("failed to load variable and secrets, err=%w", err)
	}

	infras, err := InitializeInfras(config)
	if err != nil {
		return System{}, nil, fmt.Errorf("failed to initialize context, err=%w", err)
	}

	ctx := context.Background()
	ctx = WithInfras(ctx, infras)

	domains, err := InitializeDomains(ctx, config, infras)
	if err != nil {
		return System{}, nil, fmt.Errorf("failed to initialize domains, err=%w", err)
	}

	databases, err := InitializeDatabases(ctx, config)
	if err != nil {
		return System{}, nil, fmt.Errorf("failed to initialize databases, err=%w", err)
	}

	repositories, err := InitializeRepositories(ctx, config, databases)
	if err != nil {
		return System{}, nil, fmt.Errorf("failed to initialize repositories, err=%w", err)
	}

	usecases, err := InitializeUsecases(ctx, config, infras, databases, domains, repositories)
	if err != nil {
		return System{}, nil, fmt.Errorf("failed to initialize usecases, err=%w", err)
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
