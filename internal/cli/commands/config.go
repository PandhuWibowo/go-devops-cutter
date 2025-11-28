package commands

import (
	"fmt"

	"github.com/PandhuWibowo/go-devops-cutter/internal/cli/config"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}

	cmd.AddCommand(newConfigListCmd())
	cmd.AddCommand(newConfigPathCmd())

	return cmd
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			fmt.Println("Current Configuration:")
			fmt.Printf("  Server:   %s\n", cfg.Server)
			fmt.Printf("  Username: %s\n", cfg.Username)

			if cfg.AccessToken != "" {
				tokenLen := len(cfg.AccessToken)
				if tokenLen > 20 {
					tokenLen = 20
				}
				fmt.Printf("  Token:    %s...\n", cfg.AccessToken[:tokenLen])
			}

			return nil
		},
	}
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show config file location",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.GetConfigPath())
		},
	}
}
