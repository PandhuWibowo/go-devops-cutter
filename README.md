# go-devops-cutter

![License](https://img.shields.io/badge/license-MIT-green)
![Go Version](https://img.shields.io/badge/Go-1.24-blue)
![Version](https://img.shields.io/badge/version-0.1.0-orange)

## 📝 Description

Go-Devops-Cutter is a lightweight, stateless DevOps toolkit built in Go that streamlines database backup operations. It provides a simple, reliable way to backup PostgreSQL and MySQL databases directly to your local machine using Docker-based database clients.

Perfect for DevOps engineers who need quick, reliable database backup solutions without installing database clients, managing state, or dealing with complex orchestration tools.

## ✨ Features

- 🗄️ **Direct Database Backup** - Backup PostgreSQL and MySQL databases to local machine
- 🐳 **Docker-based Clients** - No need to install database clients locally
- 🔒 **SSH Jump Host Support** - Secure access to databases behind firewalls
- 📦 **Auto Compression** - Built-in gzip compression for backups
- 📋 **Backup Management** - List and track backup files
- ⚡ **Fast & Lightweight** - Single binary, no external dependencies
- 🛠️ **Simple CLI** - Easy to use command-line interface
- 🚫 **Stateless Architecture** - No database required, no state management
- 🏥 **Health Check API** - Simple HTTP health endpoint for monitoring

## 🛠️ Tech Stack

- **Language**: Go 1.24
- **CLI Framework**: Cobra
- **Web Framework**: Gin (for API server)
- **Containerization**: Docker (for database clients)

## 📦 Key Dependencies

```go
github.com/spf13/cobra            v1.10.1    // CLI framework
github.com/gin-gonic/gin          v1.10.0    // HTTP web framework
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or higher
- Docker (for database backup operations)
- Make

### Installation

```bash
# Clone the repository
git clone https://github.com/PandhuWibowo/go-devops-cutter.git
cd go-devops-cutter

# Install dependencies
make deps

# Build all binaries (CLI + API server)
make build

# Or build individually
make build-cli    # Build CLI only
make build-api    # Build API server only

# Install CLI system-wide (optional)
make install-cli
```

## 💻 CLI Commands

### Available Commands

#### Database Operations

**`cutter db backup`** - Backup database directly to local machine

**Flags:**
- `--type` - Database type: `postgres` or `mysql` (default: postgres)
- `--host` - Database host (default: localhost)
- `--port` - Database port (default: 5432)
- `--username` - Database username (required)
- `--password` - Database password
- `--database` - Database name (required)
- `--output` - Output file path (auto-generated if not specified)
- `--compress` - Compress with gzip (default: true)
- `--ssh-jump` - SSH jump host (format: user@host)

**`cutter db list`** - List backup files in current directory

### Usage Examples

```bash
# PostgreSQL backup to local machine
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

# Backup via SSH jump host (for databases behind firewall)
cutter db backup \
  --type postgres \
  --host 10.0.1.50 \
  --port 5432 \
  --username app \
  --password pass \
  --database internal_db \
  --ssh-jump devops@jumphost.company.com

# List all backup files
cutter db list
```

## 🌐 API Server

The project includes a lightweight, stateless API server for health monitoring.

### Running the API Server

```bash
# Run directly
make run-api

# Or run the binary
./build/devops-cutter-api

# Custom port (default: 8080)
PORT=3000 ./build/devops-cutter-api
```

### API Endpoints

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
├── LICENSE
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   ├── api/
│   │   └── main.go              # API server entry point
│   └── cutter/
│       └── main.go              # CLI entry point
├── internal/
│   └── cli/
│       ├── commands/
│       │   └── db.go            # Database backup commands
│       └── ui/                  # CLI UI components
├── pkg/
│   └── client/
│       └── client.go            # HTTP client utilities
└── build/                       # Build artifacts (generated)
    ├── cutter                   # CLI binary
    └── devops-cutter-api        # API server binary
```

## 🔧 Development

### Make Commands

```bash
# Development
make deps          # Download and tidy dependencies
make build         # Build all binaries (API + CLI)
make build-cli     # Build CLI only
make build-api     # Build API server only
make run-api       # Run API server locally
make test          # Run tests
make clean         # Clean build artifacts

# Installation
make install-cli   # Install CLI to /usr/local/bin
```

### Development Setup

1. **Install Go 1.24+**
   ```bash
   # Check your Go version
   go version
   ```

2. **Install Docker**
   ```bash
   # The CLI uses Docker to run database clients
   docker --version
   ```

3. **Install Dependencies**
   ```bash
   make deps
   ```

4. **Build and Test**
   ```bash
   # Build the CLI
   make build-cli

   # Test it out
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
```

## 🚀 Deployment

### CLI Deployment

```bash
# Build for production
make build-cli

# Deploy binary to remote server
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

# Run with systemd or supervisor
# Example systemd service:
# [Unit]
# Description=DevOps Cutter API
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

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/PandhuWibowo/go-devops-cutter.git`
3. **Create** a new branch: `git checkout -b feature/your-feature`
4. **Make** your changes and add tests
5. **Test** your changes: `make test`
6. **Commit** your changes: `git commit -am 'Add some feature'`
7. **Push** to your branch: `git push origin feature/your-feature`
8. **Open** a pull request

Please ensure your code:
- Follows Go best practices and idioms
- Includes appropriate tests
- Has clear commit messages
- Updates documentation as needed

## 📜 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - Powerful CLI framework for Go
- [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- [Docker](https://www.docker.com/) - Containerization platform for database clients

---

Made with ❤️ by [Pandhu Wibowo](https://github.com/PandhuWibowo)
