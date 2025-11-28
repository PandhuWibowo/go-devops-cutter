package commands

import (
	"fmt"
	"syscall"

	"github.com/PandhuWibowo/go-devops-cutter/internal/cli/config"
	"github.com/PandhuWibowo/go-devops-cutter/pkg/client"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewLoginCmd() *cobra.Command {
	var server string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with the DevOps server",
		Long:  `Login to the DevOps Cutter server and save authentication token`,
		Example: `  cutter login --server http://localhost:8080
  cutter login --server https://devops.company.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(server)
		},
	}

	cmd.Flags().StringVarP(&server, "server", "s", "http://localhost:8080", "Server URL")
	cmd.MarkFlagRequired("server")

	return cmd
}

func runLogin(server string) error {
	fmt.Print("Username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %v", err)
	}
	password := string(passwordBytes)
	fmt.Println()

	apiClient := client.NewClient(server, "")
	resp, err := apiClient.Login(username, password)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	cfg := &config.Config{
		Server:      server,
		AccessToken: resp.Token,
		Username:    resp.Username,
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Printf("✓ Successfully logged in as %s\n", resp.Username)
	return nil
}
