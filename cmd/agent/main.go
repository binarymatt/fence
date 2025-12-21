package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/binarymatt/fence/pkg/server"
)

func main() {

	cmd := &cli.Command{
		Name:   "fence",
		Usage:  "metamorphosis event hub",
		Action: action,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "",
				Sources: cli.EnvVars("FENCE_CONFIG"),
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("failed during run", "error", err)
		os.Exit(1)
	}
}
func action(ctx context.Context, cmd *cli.Command) error {
	configPath := cmd.String("config")
	cfg, err := server.LoadConfig(configPath)
	if err != nil {
		return err
	}
	srv, err := server.New(ctx, cfg)
	if err != nil {
		return err
	}
	return srv.Run(ctx)

}
