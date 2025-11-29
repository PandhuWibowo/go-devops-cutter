package main

import (
	"fmt"
	"os"

	"github.com/PandhuWibowo/go-devops-cutter/internal/cli/commands"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "cutter",
		Short: "DevOps Cutter - Self-Service CLI for DevOps Tasks",
		Long: `Cutter is a CLI tool that helps developers self-service common DevOps tasks:
- Create and download database backups
- Direct database connection and backup
- Search and export logs
- Request database access`,
		Version: version,
	}

	rootCmd.AddCommand(commands.NewDBCmd())

	rootCmd.SetVersionTemplate(`{{printf "cutter version %s\n" .Version}}`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
