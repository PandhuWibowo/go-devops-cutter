# go-devops-cutter

![License](https://img.shields.io/badge/license-MIT-green)

## 📝 Description

Go-Devops-Cutter is a versatile tool built in Go, designed to streamline and automate various DevOps tasks. It offers a powerful combination of features, including a robust API for seamless integration with existing infrastructure, a built-in database for efficient data management, and a command-line interface (CLI) for easy interaction and scripting. Whether you need to automate deployments, manage infrastructure, or monitor system performance, Go-Devops-Cutter provides a comprehensive solution for modern DevOps workflows.

## ✨ Features

- 🌐 Api
- 🗄️ Database
- 💻 Cli


## 🛠️ Tech Stack

- 🐹 Go


## 📦 Key Dependencies

```
(: latest
```

## 🚀 Run Commands

- **all**: `make all`
- **deps**: `make deps`
- **build**: `make build`
- **build-api**: `make build-api`
- **build-worker**: `make build-worker`
- **build-cli**: `make build-cli`
- **install-cli**: `make install-cli`
- **run-api**: `make run-api`
- **run-worker**: `make run-worker`
- **docker-up**: `make docker-up`
- **docker-down**: `make docker-down`
- **test**: `make test`
- **clean**: `make clean`
- **help**: `make help`
- **Run**: `go run .`
- **Build**: `go build`


## 📁 Project Structure

```
.
├── LICENSE
├── Makefile
├── cmd
│   ├── api
│   │   └── main.go
│   ├── cutter
│   │   └── main.go
│   └── worker
│       └── main.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── auth.go
│   │   ├── backup.go
│   │   └── routes.go
│   ├── cli
│   │   ├── commands
│   │   │   ├── config.go
│   │   │   └── db.go
│   │   └── config
│   │       └── config.go
│   ├── database
│   │   ├── database.go
│   │   └── models.go
│   └── queue
│       └── handler.go
└── pkg
    └── client
        └── client.go
```

## 💻 CLI Commands

### Available Commands

- **`cutter config`** - Manage CLI configuration
  - `cutter config list` - Show current configuration
  - `cutter config path` - Show config file location

- **`cutter db`** - Direct database operations
  - `cutter db backup` - Backup database directly to local machine
    - Supports PostgreSQL and MySQL
    - Auto-compresses with gzip
    - Can use SSH jump host
  - `cutter db list` - List backup files in current directory

### CLI Usage Examples

```bash
# Build the CLI
make build-cli

# Direct PostgreSQL backup
./build/cutter db backup --type postgres --host localhost --port 5432 \
  --username myuser --password mypass --database mydb

# Backup via SSH jump host
./build/cutter db backup --type postgres --host 10.0.1.10 --port 5432 \
  --username myuser --password mypass --database mydb \
  --ssh-jump user@jumphost.com

# List backup files
./build/cutter db list

# Show configuration
./build/cutter config list

# Install CLI system-wide
make install-cli
cutter --help
```

## 🛠️ Development Setup

### Go Setup
1. Install Go (v1.18+ recommended)
2. Install dependencies: `go mod download`
3. Run the project: `go run .`


## 👥 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/PandhuWibowo/go-devops-cutter.git`
3. **Create** a new branch: `git checkout -b feature/your-feature`
4. **Commit** your changes: `git commit -am 'Add some feature'`
5. **Push** to your branch: `git push origin feature/your-feature`
6. **Open** a pull request

Please ensure your code follows the project's style guidelines and includes tests where applicable.

## 📜 License

This project is licensed under the MIT License.

---

Made with ❤️ by [Pandhu Wibowo](https://github.com/PandhuWibowo)