# Go PDF Service

A standalone microservice built in Go that generates PDF reports for students by consuming data from the Node.js backend API. This service is part of the Tailormind School Management System.

## 📋 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Quick Start](#quick-start)
- [Available Commands](#available-commands)
- [Testing](#testing)
- [API Endpoints](#api-endpoints)
- [Folder Structure](#folder-structure)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)

## 🔍 Overview

The Go PDF Service provides:
- **PDF Report Generation**: Creates formatted student reports with complete information
- **API Integration**: Fetches student data from the Node.js backend
- **Hot Reloading**: Development environment with automatic code reloading
- **Comprehensive Testing**: Unit and integration tests
- **Containerized Development**: Docker-based development environment

## 📦 Prerequisites

Before starting, ensure you have the following installed:

- **Docker** (version 20.0+)
- **Docker Compose** (version 2.0+)
- **Make** (for running commands)
- **Go** (version 1.21+) - for local development
- **curl** - for API testing

## 🚀 Installation & Setup

### 1. Clone and Navigate

```bash
# Navigate to the go-service directory
cd go-service
```

### 2. Start Development Environment

```bash
# Start all required services (PostgreSQL, Node.js Backend, Go PDF Service)
make go-dev-up-d

# Wait for services to initialize (approximately 30-60 seconds)
```

### 3. Verify Installation

```bash
# Check service status
make status

# Test health endpoint
make health-check
```

## ⚡ Quick Start

### Start Development Services

```bash
# Option 1: Start go-dev services (PostgreSQL + Node.js + Go Service)
make go-dev-up-d

# Option 2: Start all services including frontend
make all-up-d
```

### Test PDF Generation

```bash
# 1. Check if services are running
make health-check

# 2. Generate a PDF report (returns JSON with file info)
make generate-report

# 3. Download the PDF file directly
make report-download
```

### Quick Test Workflow

```bash
# Run health check and generate report in one command
make quick-test
```

## 🛠️ Available Commands

### Docker Services

| Command | Description |
|---------|-------------|
| `make go-dev-up` | Start go-dev services (interactive mode) |
| `make go-dev-up-d` | Start go-dev services (detached mode) |
| `make all-up` | Start all services including frontend |
| `make all-up-d` | Start all services (detached mode) |
| `make go-dev-down` | Stop go-dev services |
| `make all-down` | Stop all services and remove volumes |

### Development

| Command | Description |
|---------|-------------|
| `make build` | Build Go application locally |
| `make run` | Run application locally |
| `make clean` | Clean build artifacts |
| `make logs` | Show go-pdf-service logs |
| `make logs-all` | Show all services logs |
| `make status` | Show container status |

### API Testing

| Command | Description |
|---------|-------------|
| `make health-check` | Test health endpoint |
| `make generate-report` | Generate PDF report (JSON response) |
| `make report-download` | Download PDF report file |
| `make quick-test` | Quick health + report test |

### Testing

| Command | Description |
|---------|-------------|
| `make run-test-services` | Start services needed for testing |
| `make test` | Run all unit tests |
| `make test-verbose` | Run tests with verbose output |

## 🧪 Testing

### Prerequisites for Testing

Before running tests, you need the required services:

```bash
# Option 1: Use test services profile
make run-test-services

# Option 2: Use full development environment
make go-dev-up-d
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with detailed output
make test-verbose

# Run tests locally (without Docker)
go test ./...
```

### Test Categories

- **Unit Tests**: Test individual components and functions
- **Integration Tests**: Test API endpoints and external integrations
- **Service Tests**: Test PDF generation and Node.js API communication

## 📡 API Endpoints

### Health Check
```bash
GET /api/v1/health
```
**Example:**
```bash
curl http://localhost:8080/api/v1/health
```

### Generate Student Report
```bash
GET /api/v1/students/{id}/report
```
**Example:**
```bash
# Get report info (JSON)
curl http://localhost:8080/api/v1/students/1/report

# Download PDF file
curl -o report.pdf "http://localhost:8080/api/v1/students/1/report?download=true"
```

## 📁 Folder Structure

```
go-service/
├── api/                          # API layer
│   ├── router.go                 # Main router setup
│   └── v1/                       # Version 1 API
│       ├── routes.go             # Route handlers
│       └── v1_router.go          # V1 router configuration
├── cmd/                          # Application entry points
│   └── server.go                 # Main application server
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   │   └── config.go             # Config loading and validation
│   ├── models/                   # Data models
│   │   └── student.go            # Student model definitions
│   └── service/                  # Business logic
│       ├── pdf_service.go        # PDF generation service
│       └── pdf_service_test.go   # Service tests
├── reports/                      # Generated PDF reports (gitignored)
├── testdata/                     # Test data files
├── tmp/                          # Temporary build files (gitignored)
├── .air.toml                     # Hot reload configuration
├── .gitignore                    # Git ignore rules
├── config.env                    # Environment configuration
├── docker-compose-dev.yaml       # Development Docker setup
├── Dockerfile                    # Container build instructions
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── Makefile                      # Build and development commands
└── README.md                     # This documentation
```

### Key Directories Explained

- **`api/`**: Contains HTTP handlers and routing logic
- **`cmd/`**: Application entry points and main functions
- **`internal/`**: Private code that cannot be imported by other projects
- **`internal/config/`**: Configuration loading from environment variables
- **`internal/models/`**: Data structures and API response models
- **`internal/service/`**: Core business logic for PDF generation
- **`reports/`**: Output directory for generated PDF files

## 🔧 Environment Variables

The service can be configured using the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `0.0.0.0` | Server host |
| `NODEJS_API_URL` | `http://backend:5007` | Node.js backend URL |
| `NODEJS_API_TIMEOUT` | `30` | API timeout in seconds |
| `PDF_OUTPUT_DIR` | `./reports` | PDF output directory |
| `LOG_LEVEL` | `info` | Logging level |
| `AUTH_TOKEN` | - | Authentication token for Node.js API |

## 🔧 Troubleshooting

### Common Issues

#### Services Not Starting
```bash
# Check service status
make status

# View logs for debugging
make logs-all

# Restart services
make go-dev-down && make go-dev-up-d
```

#### API Connection Issues
```bash
# Verify Node.js backend is running
curl http://localhost:5007/api/v1/health

# Check Docker network connectivity
docker network ls
```

#### PDF Generation Fails
```bash
# Check service logs
make logs

# Verify reports directory exists
ls -la reports/

# Test with different student ID
curl http://localhost:8080/api/v1/students/2/report
```

#### Test Failures
```bash
# Ensure test services are running
make run-test-services

# Wait for services to be fully ready
sleep 30

# Run tests with verbose output
make test-verbose
```

### Port Conflicts

If you encounter port conflicts, you can modify the ports in `docker-compose-dev.yaml`:

```yaml
ports:
  - "8081:8080"  # Change external port to 8081
```

### Permission Issues

If you encounter Docker permission issues:

```bash
# Add your user to docker group
sudo usermod -aG docker $USER

# Or run with sudo (as configured in Makefile)
sudo docker-compose -f docker-compose-dev.yaml ps
```

## 📞 Support

For issues and questions:

1. Check the logs: `make logs`
2. Verify service status: `make status`
3. Review environment variables in `config.env`
4. Check Docker and Docker Compose versions
5. Ensure all prerequisites are installed

## 🎯 Development Workflow

1. **Start Services**: `make go-dev-up-d`
2. **Verify Health**: `make health-check`
3. **Test PDF Generation**: `make generate-report`
4. **Run Tests**: `make test`
5. **View Logs**: `make logs`
6. **Stop Services**: `make go-dev-down`

---

**Note**: This service requires the Node.js backend and PostgreSQL database to be running for full functionality. The Makefile commands handle this automatically when using the `go-dev` or `all` profiles. 