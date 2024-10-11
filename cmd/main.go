package main

import (
	"github.com/spf13/cobra"
	"github.com/xybor/todennus-backend/cmd/rest"
)

var iniPaths []string
var envPaths []string

var rootCommand = &cobra.Command{
	Use:   "todennus",
	Short: "todennus is an Identity, OpenID Connect, and OAuth2 provider",
}

func main() {
	rootCommand.PersistentFlags().StringArrayVar(&iniPaths, "ini", []string{"config/default.ini"}, "INI file paths")
	rootCommand.PersistentFlags().StringArrayVar(&envPaths, "env", []string{"config/.env"}, "ENV file paths")

	rootCommand.AddCommand(rest.Command)

	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
