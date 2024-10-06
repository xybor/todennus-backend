package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xybor/todennus-backend/adapter/rest"
	"github.com/xybor/todennus-backend/config"
	"github.com/xybor/todennus-backend/pkg/xcontext"
	"github.com/xybor/todennus-backend/wiring"
)

var iniPaths []string
var envPaths []string

var rootCommand = &cobra.Command{
	Use:   "todennus",
	Short: "todennus is an authentication provider which supports OAuth2",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.Load(sources(envPaths), sources(iniPaths))
		if err != nil {
			fmt.Printf("Cannot load variable and secrets, err=%v\n", err)
			return
		}

		infras, err := wiring.InitializeInfras(config)
		if err != nil {
			fmt.Printf("Cannot initialize context, err=%v\n", err)
			return
		}

		ctx := context.Background()
		ctx = wiring.WithInfras(ctx, infras)

		domains, err := wiring.InitializeDomains(ctx, config, infras)
		if err != nil {
			fmt.Printf("Cannot initialize domains, err=%v\n", err)
			return
		}

		databases, err := wiring.InitializeDatabases(ctx, config)
		if err != nil {
			fmt.Printf("Cannot initialize databases, err=%v\n", err)
			return
		}

		repositories, err := wiring.InitializeRepositories(ctx, databases)
		if err != nil {
			fmt.Printf("Cannot initialize repositories, err=%v\n", err)
			return
		}

		usecases, err := wiring.InitializeUsecases(ctx, infras, domains, repositories)
		if err != nil {
			fmt.Printf("Cannot initialize usecases, err=%v\n", err)
			return
		}

		address := fmt.Sprintf("%s:%d", config.Variable.Server.Host, config.Variable.Server.Port)
		app := rest.App(infras, usecases)

		xcontext.Logger(ctx).Info("Server started", "address", address)
		if err := http.ListenAndServe(address, app); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCommand.Flags().StringArrayVar(&iniPaths, "ini", []string{"config/default.ini"}, "INI file paths")
	rootCommand.Flags().StringArrayVar(&envPaths, "env", []string{"config/.env"}, "ENV file paths")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
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
