package main

import (
	"github.com/spf13/cobra"
	"github.com/xybor/todennus-backend/cmd/grpc"
	"github.com/xybor/todennus-backend/cmd/rest"
	"github.com/xybor/todennus-backend/cmd/swagger"
)

var rootCommand = &cobra.Command{
	Use:   "todennus",
	Short: "todennus is an Identity, OpenID Connect, and OAuth2 provider",
}

func main() {
	rootCommand.PersistentFlags().StringArray("env", []string{".env"}, "environment file paths")
	rootCommand.AddCommand(rest.Command)
	rootCommand.AddCommand(grpc.Command)
	rootCommand.AddCommand(swagger.Command)

	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
