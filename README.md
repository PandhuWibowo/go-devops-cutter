# 🔪 Go DevOps Cutter

> Self-service DevOps automation tool for developers

**Go DevOps Cutter** adalah tool open-source yang membantu developer melakukan self-service untuk task DevOps yang sering dilakukan, seperti database backup, log export, dan database access management, tanpa perlu menunggu tim DevOps.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## ✨ Features

- 🗄️ **Database Backup On-Demand** - Create, monitor, dan download database backups
- 📝 **Log Export** - Search dan export logs dari berbagai services (coming soon)
- 🔐 **Database Access Management** - Request temporary database access dengan RBAC (coming soon)
- 🔔 **Telegram Notifications** - Real-time notifications untuk task completion
- 🎯 **CLI Tool** - Easy-to-use command line interface
- 🔄 **Background Processing** - Async job processing dengan Redis Queue
- 🔒 **JWT Authentication** - Secure API authentication

## 🏗️ Architecture
```
┌─────────────┐
│   CLI Tool  │
│  (cutter)   │
└──────┬──────┘
       │
       │ HTTP/REST
       │
┌──────▼──────────────────────┐
│      API Server (Gin)       │
│  - Authentication           │
│  - Request Validation       │
│  - Job Enqueueing          │
└──────┬──────────────────────┘
       │
       │ Redis Queue
       │
┌──────▼──────────────────────┐
│   Background Worker         │
│  - Backup Execution         │
│  - Log Export               │
│  - Notifications            │
└──────┬──────────────────────┘
       │
   ┌───┴────┬─────────┐
   │        │         │
┌──▼───┐ ┌─▼──┐ ┌────▼────┐
│ PG   │ │Redis│ │ Storage │
│ SQL  │ │     │ │ (OBS/S3)│
└──────┘ └─────┘ └─────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Installation

#### 1. Clone Repository
```bash
git clone https://github.com/PandhuWibowo/go-devops-cutter.git
cd go-devops-cutter
```

#### 2. Install Dependencies
```bash
make deps
```

#### 3. Start Services
```bash
# Start PostgreSQL & Redis
make docker-up

# Start API Server
make run-api

# Start Worker (in another terminal)
make run-worker
```

#### 4. Install CLI
```bash
make build-cli
sudo make install-cli
```

### Quick Test
```bash
# Create first user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123"
  }'

# Login via CLI
cutter login --server http://localhost:8080

# Create backup
cutter backup create --db production-db

# List backups
cutter backup list
```

## 📖 Usage

### CLI Commands

#### Authentication
```bash
# Login to server
cutter login --server http://localhost:8080

# View config
cutter config list
```

#### Backup Management
```bash
# Create backup
cutter backup create --db production-db \
  --format sql \
  --compress \
  --backup-type full

# List all backups
cutter backup list

# Check backup status
cutter backup status 1

# Download backup (coming soon)
cutter backup download 1 -o backup.sql.gz
```

#### Database Access (Coming Soon)
```bash
# Request database access
cutter db access request \
  --db staging-db \
  --level read \
  --duration 7d

# List active access
cutter db access list

# Revoke access
cutter db access revoke 1
```

#### Log Export (Coming Soon)
```bash
# Search logs
cutter logs search \
  --service api-gateway \
  --from 2024-01-01 \
  --to 2024-01-31 \
  --keyword "error"

# Export logs
cutter logs export \
  --service api-gateway \
  --from 2024-01-01 \
  --to 2024-01-31 \
  -o logs.zip
```

### API Endpoints

#### Authentication
```bash
# Register
POST /api/v1/auth/register
{
  "username": "user",
  "email": "user@example.com",
  "password": "password"
}

# Login
POST /api/v1/auth/login
{
  "username": "user",
  "password": "password"
}
```

#### Backups
```bash
# Create backup
POST /api/v1/backups
Authorization: Bearer <token>
{
  "database_name": "production-db",
  "format": "sql",
  "compression": true,
  "backup_type": "full"
}

# List backups
GET /api/v1/backups
Authorization: Bearer <token>

# Get backup status
GET /api/v1/backups/:id
Authorization: Bearer <token>
```

## 🛠️ Development

### Project Structure
```
go-devops-cutter/
├── cmd/
│   ├── api/              # API server entrypoint
│   ├── worker/           # Background worker entrypoint
│   └── cutter/           # CLI tool entrypoint
├── internal/
│   ├── api/              # API handlers & routes
│   ├── auth/             # Authentication logic
│   ├── backup/           # Backup executors
│   ├── cli/              # CLI commands
│   │   ├── commands/     # CLI command implementations
│   │   └── config/       # CLI configuration
│   ├── database/         # Database models & migrations
│   ├── notification/     # Notification services
│   ├── queue/            # Background job handlers
│   └── storage/          # Object storage integration
├── pkg/
│   └── client/           # API client library
├── configs/              # Configuration files
├── migrations/           # Database migrations
├── docker-compose.yml    # Docker compose for dev
├── Makefile             # Build automation
└── go.mod
```

### Available Make Commands
```bash
make deps          # Download dependencies
make build         # Build all binaries
make build-api     # Build API server only
make build-worker  # Build worker only
make build-cli     # Build CLI only
make install-cli   # Install CLI to /usr/local/bin
make run-api       # Run API server locally
make run-worker    # Run worker locally
make docker-up     # Start all services with Docker
make docker-down   # Stop all services
make test          # Run tests
make clean         # Clean build artifacts
```

### Running Tests
```bash
# Run all tests
make test

# Run with coverage
go test -v -cover ./...

# Run specific package
go test -v ./internal/api/...
```

### Building for Production
```bash
# Build all binaries
make build

# Build for multiple platforms
GOOS=linux GOARCH=amd64 make build
GOOS=darwin GOARCH=arm64 make build
GOOS=windows GOARCH=amd64 make build
```

## 🔧 Configuration

### Environment Variables
```bash
# API Server
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/devops_cutter
REDIS_ADDR=localhost:6379
JWT_SECRET=your-secret-key

# Object Storage (Huawei OBS)
OBS_ENDPOINT=https://obs.ap-southeast-3.myhuaweicloud.com
OBS_BUCKET=devops-backups
OBS_ACCESS_KEY=your-access-key
OBS_SECRET_KEY=your-secret-key

# Telegram Notifications
TELEGRAM_BOT_TOKEN=your-bot-token
```

### Database Configuration

Edit `configs/config.yaml`:
```yaml
databases:
  - name: production-db
    type: postgres
    deployment_type: container  # vm, rds, container
    container_name: postgres-prod
    docker_host: unix:///var/run/docker.sock
    credentials:
      username: backup_user
      password: ${PROD_DB_PASSWORD}
```

## 🎯 Roadmap

- [x] CLI Tool
- [x] Database Backup (Basic)
- [x] JWT Authentication
- [x] Background Job Processing
- [ ] Real Backup Implementation (SSH, pg_dump, mysqldump)
- [ ] Object Storage Integration (OBS/S3)
- [ ] Backup Download
- [ ] Telegram Notifications
- [ ] Database Access Management
- [ ] Log Export & Search
- [ ] Web UI Dashboard
- [ ] Kubernetes Support
- [ ] Multi-cloud Support

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Pandhu Wibowo**

- GitHub: [@PandhuWibowo](https://github.com/PandhuWibowo)

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Asynq](https://github.com/hibiken/asynq) - Background job processing

## 📸 Screenshots

### CLI Usage
```bash
$ cutter backup create --db production-db
Creating backup for database: production-db
✓ Backup request created (ID: 1)

$ cutter backup list
ID    Database             Status          Size         Created
--------------------------------------------------------------------------------
1     production-db        completed       5.2 GB       2025-11-28 23:13
```

---

Made with ❤️ by [Pandhu Wibowo](https://github.com/PandhuWibowo)