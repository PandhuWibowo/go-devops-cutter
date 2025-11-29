# go-devops-cutter

![License](https://img.shields.io/badge/license-MIT-green)
![Go Version](https://img.shields.io/badge/Go-1.24-blue)
![Version](https://img.shields.io/badge/version-0.1.0-orange)

## 📝 Description

Go-Devops-Cutter is a lightweight CLI tool built in Go for backing up PostgreSQL and MySQL databases to your local machine. It uses Docker-based database clients, so you don't need to install database tools locally.

Perfect for DevOps engineers who need quick, reliable database backups without installing database clients or managing complex tooling.

## ✨ Features

- 🗄️ **Database Backup** - Backup PostgreSQL and MySQL databases to local machine
- 🐳 **Docker-based Clients** - No need to install `pg_dump` or `mysqldump` locally
- 📦 **Auto Compression** - Built-in gzip compression for backups
- 📋 **Backup Listing** - List backup files in current directory
- ⚡ **Fast & Lightweight** - Single binary, minimal dependencies
- 🛠️ **Simple CLI** - Easy to use command-line interface

## 🛠️ Tech Stack

- **Language**: Go 1.24
- **CLI Framework**: Cobra
- **Containerization**: Docker (for database clients)

## 📦 Key Dependencies

```go
github.com/spf13/cobra            v1.10.1    // CLI framework
github.com/gin-gonic/gin          v1.11.0    // HTTP framework (for health API)
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or higher
- Docker (for running database clients)
- Make

### Installation

```bash
# Clone the repository
git clone https://github.com/PandhuWibowo/go-devops-cutter.git
cd go-devops-cutter

# Install dependencies
make deps

# Build the CLI
make build-cli

# Install CLI system-wide (optional)
make install-cli
```

## 💻 CLI Commands

### Database Operations

**`cutter db backup`** - Backup database to local machine

**Required Flags:**
- `--database` - Database name
- `--username` - Database username

**Optional Flags:**
- `--type` - Database type: `postgres` or `mysql` (default: postgres)
- `--host` - Database host (default: localhost)
- `--port` - Database port (default: 5432 for postgres, 3306 for mysql)
- `--password` - Database password (prompted if not provided)
- `--output` - Output file path (auto-generated if not specified)
- `--compress` - Compress with gzip (default: true)

**`cutter db list`** - List backup files (*.sql*) in current directory

### Usage Examples

```bash
# PostgreSQL backup
cutter db backup \
  --type postgres \
  --host localhost \
  --port 5432 \
  --username dbuser \
  --password dbpass \
  --database myapp

# MySQL backup with custom output
cutter db backup \
  --type mysql \
  --host 192.168.1.100 \
  --port 3306 \
  --username root \
  --password secret \
  --database production \
  --output ~/backups/prod_backup.sql.gz

# List all backup files
cutter db list
```

## 🌐 Health Check API

The project includes a minimal health check API server for monitoring.

### Running the API Server

```bash
# Build and run
make build-api
./build/devops-cutter-api

# Or run directly
make run-api

# Custom port (default: 8080)
PORT=3000 ./build/devops-cutter-api
```

### API Endpoint

**GET /health** - Health check endpoint

Response:
```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

Example:
```bash
curl http://localhost:8080/health
```

## 📁 Project Structure

```
.
├── cmd/
│   ├── api/
│   │   ├── main.go              # Health check API server
│   │   └── main_test.go         # API tests
│   └── cutter/
│       └── main.go              # CLI entry point
├── internal/
│   └── cli/
│       └── commands/
│           ├── db.go            # Database backup commands
│           └── db_test.go       # Command tests
├── pkg/
│   └── client/
│       ├── client.go            # HTTP client utilities
│       └── client_test.go       # Client tests
├── build/                       # Build artifacts (generated)
│   ├── cutter                   # CLI binary
│   └── devops-cutter-api        # API binary
├── Makefile
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## 🔧 Development

### Available Commands

- **`cutter db`** - Direct database operations
  - `cutter db backup` - Backup database directly to local machine
    - Supports PostgreSQL and MySQL
    - Auto-compresses with gzip
    - Can use SSH jump host
  - `cutter db list` - List backup files in current directory

### CLI Usage Examples

```bash
# Development
make deps          # Download and tidy dependencies
make build         # Build all binaries (API + CLI)
make build-cli     # Build CLI only
make build-api     # Build API server only
make run-api       # Run API server locally
make test          # Run tests
make clean         # Clean build artifacts

# Direct PostgreSQL backup
./build/cutter db backup --type postgres --host localhost --port 5432 \
  --username myuser --password mypass --database mydb

# Backup via SSH jump host
./build/cutter db backup --type postgres --host 10.0.1.10 --port 5432 \
  --username myuser --password mypass --database mydb \
  --ssh-jump user@jumphost.com

# List backup files
./build/cutter db list

# Install CLI system-wide
make install-cli
cutter --help
```

### Development Setup

1. **Install Go 1.24+**
   ```bash
   go version
   ```

2. **Install Docker**
   ```bash
   docker --version
   ```

3. **Install Dependencies**
   ```bash
   make deps
   ```

4. **Build and Test**
   ```bash
   make build-cli
   ./build/cutter --help
   ./build/cutter db --help
   ```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific package tests
go test -v ./internal/cli/commands
go test -v ./cmd/api
go test -v ./pkg/client
```

## 🚀 Deployment

### CLI Deployment

```bash
# Build for production
make build-cli

# Deploy to remote server
scp build/cutter user@server:/usr/local/bin/

# Or install locally
sudo cp build/cutter /usr/local/bin/
```

### API Server Deployment

```bash
# Build for production
make build-api

# Deploy to server
scp build/devops-cutter-api user@server:/opt/devops-cutter/

# Run with systemd (example)
# /etc/systemd/system/devops-cutter-api.service:
# [Unit]
# Description=DevOps Cutter Health Check API
#
# [Service]
# ExecStart=/opt/devops-cutter/devops-cutter-api
# Environment=PORT=8080
# Restart=always
#
# [Install]
# WantedBy=multi-user.target
```

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a new branch: `git checkout -b feature/your-feature`
3. Make your changes and add tests
4. Run tests: `make test`
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to your branch: `git push origin feature/your-feature`
7. Open a pull request

Please ensure your code:
- Follows Go best practices
- Includes appropriate tests
- Has clear commit messages
- Updates documentation as needed

## 📜 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go
- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [Docker](https://www.docker.com/) - Container platform for database clients

---

Made with ❤️ by [Pandhu Wibowo](https://github.com/PandhuWibowo)
