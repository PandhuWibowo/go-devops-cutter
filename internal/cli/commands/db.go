package commands

import (
	"fmt"
	"net"
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
	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed")
	}

	// Handle SSH jump host if provided
	var tunnel *sshTunnel
	var err error
	originalHost := host
	originalPort := port

	if sshJump != "" {
		// Create SSH tunnel
		tunnel, err = createSSHTunnel(sshJump, host, port)
		if err != nil {
			return fmt.Errorf("SSH tunnel failed: %v", err)
		}
		defer tunnel.close()

		// Update host and port to use tunnel
		host = "host.docker.internal"
		port = tunnel.localPort
	}

	fmt.Println("Using Docker PostgreSQL client...")

	// Build pg_dump command
	var pgDumpCmd string
	var dockerOpts string

	if sshJump != "" {
		// When using tunnel, connect through localhost via host.docker.internal
		pgDumpCmd = fmt.Sprintf("PGPASSWORD=%s pg_dump -h host.docker.internal -p %d -U %s %s",
			password, port, username, database)
		// Add host mapping for Linux compatibility
		dockerOpts = "--add-host=host.docker.internal:host-gateway"
	} else {
		// Direct connection using host network
		pgDumpCmd = fmt.Sprintf("PGPASSWORD=%s pg_dump -h %s -p %d -U %s %s",
			password, originalHost, originalPort, username, database)
		dockerOpts = "--network host"
	}

	var cmdStr string

	if compress {
		// Use Docker to run pg_dump and compress
		cmdStr = fmt.Sprintf(
			`docker run --rm %s postgres:15-alpine sh -c "%s" | gzip > %s`,
			dockerOpts, pgDumpCmd, output,
		)
	} else {
		// Use Docker to run pg_dump without compression
		cmdStr = fmt.Sprintf(
			`docker run --rm %s postgres:15-alpine sh -c "%s" > %s`,
			dockerOpts, pgDumpCmd, output,
		)
	}

	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func backupMySQL(host string, port int, username, password, database, output string, compress bool, sshJump string) error {
	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed")
	}

	// Handle SSH jump host if provided
	var tunnel *sshTunnel
	var err error
	originalHost := host
	originalPort := port

	if sshJump != "" {
		// Create SSH tunnel
		tunnel, err = createSSHTunnel(sshJump, host, port)
		if err != nil {
			return fmt.Errorf("SSH tunnel failed: %v", err)
		}
		defer tunnel.close()

		// Update host and port to use tunnel
		host = "host.docker.internal"
		port = tunnel.localPort
	}

	fmt.Println("Using Docker MySQL client...")

	// Build mysqldump command
	var mysqldumpCmd string
	var dockerOpts string

	if sshJump != "" {
		// When using tunnel, connect through localhost via host.docker.internal
		mysqldumpCmd = fmt.Sprintf("mysqldump -h host.docker.internal -P %d -u %s -p%s %s",
			port, username, password, database)
		// Add host mapping for Linux compatibility
		dockerOpts = "--add-host=host.docker.internal:host-gateway"
	} else {
		// Direct connection using host network
		mysqldumpCmd = fmt.Sprintf("mysqldump -h %s -P %d -u %s -p%s %s",
			originalHost, originalPort, username, password, database)
		dockerOpts = "--network host"
	}

	var cmdStr string

	if compress {
		cmdStr = fmt.Sprintf(
			`docker run --rm %s mysql:8 sh -c "%s" | gzip > %s`,
			dockerOpts, mysqldumpCmd, output,
		)
	} else {
		cmdStr = fmt.Sprintf(
			`docker run --rm %s mysql:8 sh -c "%s" > %s`,
			dockerOpts, mysqldumpCmd, output,
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

// sshTunnel represents an active SSH tunnel
type sshTunnel struct {
	cmd       *exec.Cmd
	localPort int
}

// findAvailablePort finds an available local port
func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// createSSHTunnel creates an SSH tunnel to the database through a jump host
func createSSHTunnel(sshJump, dbHost string, dbPort int) (*sshTunnel, error) {
	// Validate SSH jump host format (user@host or user@host:port)
	if sshJump == "" {
		return nil, fmt.Errorf("SSH jump host cannot be empty")
	}

	// Find available local port
	localPort, err := findAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find available port: %v", err)
	}

	// Build SSH tunnel command
	// Format: ssh -f -N -L local_port:db_host:db_port user@jumphost
	tunnelSpec := fmt.Sprintf("%d:%s:%d", localPort, dbHost, dbPort)

	fmt.Printf("Creating SSH tunnel through %s...\n", sshJump)
	fmt.Printf("  Local port %d -> %s:%d\n", localPort, dbHost, dbPort)

	cmd := exec.Command("ssh",
		"-f",                             // Run in background
		"-N",                             // Don't execute remote command
		"-o", "ExitOnForwardFailure=yes", // Exit if tunnel fails
		"-o", "StrictHostKeyChecking=no", // Don't prompt for host key
		"-L", tunnelSpec, // Local port forwarding
		sshJump, // Jump host
	)

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start SSH tunnel: %v", err)
	}

	// Wait a moment for tunnel to establish
	time.Sleep(2 * time.Second)

	tunnel := &sshTunnel{
		cmd:       cmd,
		localPort: localPort,
	}

	fmt.Println("✓ SSH tunnel established")

	return tunnel, nil
}

// close terminates the SSH tunnel
func (t *sshTunnel) close() error {
	if t.cmd == nil || t.cmd.Process == nil {
		return nil
	}

	fmt.Println("Closing SSH tunnel...")

	// Kill the SSH process
	if err := t.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill SSH tunnel: %v", err)
	}

	// Wait for process to exit
	_ = t.cmd.Wait()

	return nil
}
