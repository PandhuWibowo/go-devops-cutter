package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func NewDBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Direct database operations",
	}

	cmd.AddCommand(newDBBackupCmd())
	cmd.AddCommand(newDBListCmd())

	return cmd
}

func newDBBackupCmd() *cobra.Command {
	var (
		host     string
		port     int
		username string
		password string
		database string
		dbType   string
		output   string
		compress bool
		sshJump  string
	)

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup database directly to local machine",
		Example: `  # Direct PostgreSQL backup
  cutter db backup --type postgres --host localhost --port 5432 \
    --username myuser --password mypass --database mydb

  # Backup via SSH jump host
  cutter db backup --type postgres --host 10.0.1.10 --port 5432 \
    --username myuser --password mypass --database mydb \
    --ssh-jump user@jumphost.com

  # Backup with custom output
  cutter db backup --type postgres --host localhost --database mydb \
    --output ~/backups/mydb.sql.gz`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDBBackup(dbType, host, port, username, password, database, output, compress, sshJump)
		},
	}

	cmd.Flags().StringVar(&dbType, "type", "postgres", "Database type (postgres, mysql)")
	cmd.Flags().StringVar(&host, "host", "localhost", "Database host")
	cmd.Flags().IntVar(&port, "port", 5432, "Database port")
	cmd.Flags().StringVar(&username, "username", "", "Database username")
	cmd.Flags().StringVar(&password, "password", "", "Database password")
	cmd.Flags().StringVar(&database, "database", "", "Database name")
	cmd.Flags().StringVar(&output, "output", "", "Output file path (default: auto-generated)")
	cmd.Flags().BoolVar(&compress, "compress", true, "Compress with gzip")
	cmd.Flags().StringVar(&sshJump, "ssh-jump", "", "SSH jump host (user@host)")

	cmd.MarkFlagRequired("database")
	cmd.MarkFlagRequired("username")

	return cmd
}

func runDBBackup(dbType, host string, port int, username, password, database, output string, compress bool, sshJump string) error {
	// Generate output filename if not provided
	if output == "" {
		timestamp := time.Now().Format("20060102_150405")
		output = fmt.Sprintf("%s_%s.sql", database, timestamp)
		if compress {
			output += ".gz"
		}
	}

	fmt.Printf("Starting backup for %s database: %s\n", dbType, database)
	fmt.Printf("Host: %s:%d\n", host, port)
	fmt.Printf("Output: %s\n", output)

	var err error
	switch dbType {
	case "postgres":
		err = backupPostgres(host, port, username, password, database, output, compress, sshJump)
	case "mysql":
		err = backupMySQL(host, port, username, password, database, output, compress, sshJump)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return fmt.Errorf("backup failed: %v", err)
	}

	// Get file size
	fileInfo, err := os.Stat(output)
	if err == nil {
		size := float64(fileInfo.Size()) / (1024 * 1024) // MB
		fmt.Printf("\n✓ Backup completed successfully!\n")
		fmt.Printf("  File: %s\n", output)
		fmt.Printf("  Size: %.2f MB\n", size)
	}

	return nil
}

func backupPostgres(host string, port int, username, password, database, output string, compress bool, sshJump string) error {
	fmt.Println("Using Docker PostgreSQL client...")

	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed. Please install Docker or use --ssh-jump option")
	}

	// Build pg_dump command
	pgDumpCmd := fmt.Sprintf("PGPASSWORD=%s pg_dump -h %s -p %d -U %s %s",
		password, host, port, username, database)

	var cmdStr string

	if compress {
		// Use Docker to run pg_dump and compress
		cmdStr = fmt.Sprintf(
			`docker run --rm --network host postgres:15-alpine sh -c "%s" | gzip > %s`,
			pgDumpCmd, output,
		)
	} else {
		// Use Docker to run pg_dump without compression
		cmdStr = fmt.Sprintf(
			`docker run --rm --network host postgres:15-alpine sh -c "%s" > %s`,
			pgDumpCmd, output,
		)
	}

	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func backupMySQL(host string, port int, username, password, database, output string, compress bool, sshJump string) error {
	fmt.Println("Using Docker MySQL client...")

	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed. Please install Docker or use --ssh-jump option")
	}

	// Build mysqldump command
	mysqldumpCmd := fmt.Sprintf("mysqldump -h %s -P %d -u %s -p%s %s",
		host, port, username, password, database)

	var cmdStr string

	if compress {
		cmdStr = fmt.Sprintf(
			`docker run --rm --network host mysql:8 sh -c "%s" | gzip > %s`,
			mysqldumpCmd, output,
		)
	} else {
		cmdStr = fmt.Sprintf(
			`docker run --rm --network host mysql:8 sh -c "%s" > %s`,
			mysqldumpCmd, output,
		)
	}

	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func newDBListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List backup files in current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			files, _ := filepath.Glob("*.sql*")
			if len(files) == 0 {
				fmt.Println("No backup files found in current directory")
				return nil
			}

			fmt.Println("Backup files:")
			for _, file := range files {
				info, _ := os.Stat(file)
				size := float64(info.Size()) / (1024 * 1024)
				fmt.Printf("  - %s (%.2f MB)\n", file, size)
			}
			return nil
		},
	}
}
