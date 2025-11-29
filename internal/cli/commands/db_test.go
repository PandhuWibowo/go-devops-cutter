package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewDBCmd(t *testing.T) {
	cmd := NewDBCmd()

	if cmd.Use != "db" {
		t.Errorf("Expected Use 'db', got '%s'", cmd.Use)
	}

	if cmd.Short != "Direct database operations" {
		t.Errorf("Expected Short description, got '%s'", cmd.Short)
	}

	// Check that subcommands are added
	if len(cmd.Commands()) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(cmd.Commands()))
	}

	// Verify subcommands exist
	hasBackup := false
	hasList := false
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "backup" {
			hasBackup = true
		}
		if subcmd.Use == "list" {
			hasList = true
		}
	}

	if !hasBackup {
		t.Error("Expected 'backup' subcommand to exist")
	}
	if !hasList {
		t.Error("Expected 'list' subcommand to exist")
	}
}

func TestNewDBBackupCmd(t *testing.T) {
	cmd := newDBBackupCmd()

	if cmd.Use != "backup" {
		t.Errorf("Expected Use 'backup', got '%s'", cmd.Use)
	}

	// Check required flags
	flags := cmd.Flags()

	// Test that flags exist
	expectedFlags := []string{
		"type", "host", "port", "username", "password",
		"database", "output", "compress", "ssh-jump",
	}

	for _, flagName := range expectedFlags {
		flag := flags.Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to exist", flagName)
		}
	}

	// Test default values
	typeFlag := flags.Lookup("type")
	if typeFlag.DefValue != "postgres" {
		t.Errorf("Expected default type 'postgres', got '%s'", typeFlag.DefValue)
	}

	hostFlag := flags.Lookup("host")
	if hostFlag.DefValue != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", hostFlag.DefValue)
	}

	portFlag := flags.Lookup("port")
	if portFlag.DefValue != "5432" {
		t.Errorf("Expected default port '5432', got '%s'", portFlag.DefValue)
	}

	compressFlag := flags.Lookup("compress")
	if compressFlag.DefValue != "true" {
		t.Errorf("Expected default compress 'true', got '%s'", compressFlag.DefValue)
	}
}

func TestNewDBListCmd(t *testing.T) {
	cmd := newDBListCmd()

	if cmd.Use != "list" {
		t.Errorf("Expected Use 'list', got '%s'", cmd.Use)
	}

	if cmd.Short != "List backup files in current directory" {
		t.Errorf("Expected Short description, got '%s'", cmd.Short)
	}
}

func TestDBListCmdNoFiles(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to temp directory
	os.Chdir(tempDir)

	cmd := newDBListCmd()
	err := cmd.RunE(cmd, []string{})

	if err != nil {
		t.Errorf("Expected no error when no files exist, got %v", err)
	}
}

func TestDBListCmdWithFiles(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Create some test backup files
	testFiles := []string{
		"backup1.sql",
		"backup2.sql.gz",
		"backup3.sql",
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Change to temp directory
	os.Chdir(tempDir)

	cmd := newDBListCmd()
	err := cmd.RunE(cmd, []string{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRunDBBackupInvalidType(t *testing.T) {
	err := runDBBackup("invalid", "localhost", 5432, "user", "pass", "testdb", "", true, "")

	if err == nil {
		t.Error("Expected error for invalid database type, got nil")
	}

	expectedError := "unsupported database type: invalid"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDBBackupCmdValidation(t *testing.T) {
	cmd := newDBBackupCmd()

	// Test that database flag is required
	cmd.SetArgs([]string{"--username", "test"})
	err := cmd.Execute()

	if err == nil {
		t.Error("Expected error when required flags are missing")
	}
}

func TestDBBackupFlagsBinding(t *testing.T) {
	cmd := newDBBackupCmd()
	flags := cmd.Flags()

	tests := []struct {
		name         string
		flagName     string
		expectedType string
	}{
		{"Type flag", "type", "string"},
		{"Host flag", "host", "string"},
		{"Port flag", "port", "int"},
		{"Username flag", "username", "string"},
		{"Password flag", "password", "string"},
		{"Database flag", "database", "string"},
		{"Output flag", "output", "string"},
		{"Compress flag", "compress", "bool"},
		{"SSH Jump flag", "ssh-jump", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.flagName)
			if flag == nil {
				t.Fatalf("Flag '%s' not found", tt.flagName)
			}

			if flag.Value.Type() != tt.expectedType {
				t.Errorf("Expected type '%s', got '%s'", tt.expectedType, flag.Value.Type())
			}
		})
	}
}

func TestDBCommandStructure(t *testing.T) {
	rootCmd := &cobra.Command{Use: "cutter"}
	dbCmd := NewDBCmd()
	rootCmd.AddCommand(dbCmd)

	// Test that the command tree is properly structured
	if dbCmd.Parent() != rootCmd {
		t.Error("DB command should have root command as parent")
	}

	// Test that backup and list are children of db
	backupCmd := dbCmd.Commands()[0]
	if backupCmd.Parent() != dbCmd {
		t.Error("Backup command should have db command as parent")
	}
}
