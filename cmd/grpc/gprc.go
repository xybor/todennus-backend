package grpc

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"github.com/xybor/todennus-backend/adapter/grpc"
	"github.com/xybor/todennus-backend/wiring"
	"github.com/xybor/x/xcontext"
)

var Command = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		envPaths, err := cmd.Flags().GetStringArray("env")
		if err != nil {
			panic(err)
		}

		system, ctx, err := wiring.InitializeSystem(envPaths...)
		if err != nil {
			panic(err)
		}

		address := fmt.Sprintf("%s:%d", system.Config.Variable.Server.Host, system.Config.Variable.Server.Port)
		app := grpc.App(system.Config, system.Infras, system.Usecases)

		listener, err := net.Listen("tcp", address)
		if err != nil {
			panic(err)
		}

		xcontext.Logger(ctx).Info("gRPC server started", "address", address)
		if err := app.Serve(listener); err != nil {
			panic(err)
		}
	},
}
