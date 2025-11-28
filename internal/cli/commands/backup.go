package commands

import (
	"fmt"

	"github.com/PandhuWibowo/go-devops-cutter/internal/cli/config"
	"github.com/PandhuWibowo/go-devops-cutter/pkg/client"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Manage database backups",
	}

	cmd.AddCommand(newBackupCreateCmd())
	cmd.AddCommand(newBackupListCmd())
	cmd.AddCommand(newBackupStatusCmd())

	return cmd
}

func newBackupCreateCmd() *cobra.Command {
	var (
		database   string
		format     string
		compress   bool
		backupType string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new database backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackupCreate(database, format, compress, backupType)
		},
	}

	cmd.Flags().StringVar(&database, "db", "", "Database name")
	cmd.Flags().StringVar(&format, "format", "sql", "Backup format")
	cmd.Flags().BoolVar(&compress, "compress", true, "Compress backup")
	cmd.Flags().StringVar(&backupType, "backup-type", "full", "Backup type")
	cmd.MarkFlagRequired("db")

	return cmd
}

func runBackupCreate(database, format string, compress bool, backupType string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if !cfg.IsAuthenticated() {
		return fmt.Errorf("not authenticated. Run 'cutter login' first")
	}

	apiClient := client.NewClient(cfg.Server, cfg.AccessToken)
	fmt.Printf("Creating backup for database: %s\n", database)

	resp, err := apiClient.CreateBackup(client.BackupCreateRequest{
		DatabaseName: database,
		Format:       format,
		Compression:  compress,
		BackupType:   backupType,
	})
	if err != nil {
		return err
	}

	fmt.Printf("✓ Backup request created (ID: %d)\n", resp.RequestID)
	return nil
}

func newBackupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all backup requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			if !cfg.IsAuthenticated() {
				return fmt.Errorf("not authenticated")
			}
			apiClient := client.NewClient(cfg.Server, cfg.AccessToken)
			backups, _ := apiClient.ListBackups()

			if len(backups) == 0 {
				fmt.Println("No backups found")
				return nil
			}

			fmt.Printf("%-5s %-20s %-15s %-20s\n", "ID", "Database", "Status", "Created")
			fmt.Println("----------------------------------------------------------------")
			for _, b := range backups {
				fmt.Printf("%-5d %-20s %-15s %-20s\n",
					b.ID, b.DatabaseName, b.Status, b.CreatedAt.Format("2006-01-02 15:04"))
			}
			return nil
		},
	}
}

func newBackupStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <id>",
		Short: "Check backup status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			if !cfg.IsAuthenticated() {
				return fmt.Errorf("not authenticated")
			}

			var id uint
			fmt.Sscanf(args[0], "%d", &id)

			apiClient := client.NewClient(cfg.Server, cfg.AccessToken)
			status, err := apiClient.GetBackupStatus(id)
			if err != nil {
				return err
			}

			fmt.Printf("Backup Request #%d\n", status.ID)
			fmt.Printf("  Database: %s\n", status.DatabaseName)
			fmt.Printf("  Status:   %s\n", status.Status)
			return nil
		},
	}
}
