# Development Setup Guide

This guide will help you set up a local development environment for the Cognize API server.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go 1.24.6 or higher** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 17+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Git** - [Download Git](https://git-scm.com/downloads)
- **Docker & Docker Compose** (optional but recommended) - [Download Docker](https://www.docker.com/get-started)

## Environment Setup

### 1. Clone the Repository

```bash
git clone https://github.com/Cognize-AI/server-cognize.git
cd server-cognize
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
PORT=4000
ENVIRONMENT=dev

# Database Configuration
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Google OAuth Configuration
GOOGLE_OAUTH_CLIENT_ID=your-google-oauth-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback

# Email Configuration (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Logging Configuration (Optional)
AXIOM_TOKEN=your-axiom-token
AXIOM_ORG=your-axiom-org
AXIOM_DATASET=cognize-logs

# Encryption (Optional)
ENC_SECRET=your-encryption-secret-key
```

### 4. Set Up Google OAuth

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Create OAuth 2.0 credentials:
   - Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client IDs"
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:3000/auth/callback`
5. Copy the Client ID and Client Secret to your `.env` file

## Database Setup

### Option 1: Using Docker (Recommended)

Start a PostgreSQL container:

```bash
docker run --name cognize-postgres \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=cognize \
  -p 5432:5432 \
  -d postgres:17
```

### Option 2: Local PostgreSQL Installation

1. Install PostgreSQL on your system
2. Create a database:
   ```sql
   CREATE DATABASE cognize;
   CREATE USER root WITH PASSWORD 'root';
   GRANT ALL PRIVILEGES ON DATABASE cognize TO root;
   ```

### Database Migration

The application will automatically create the necessary tables on startup using GORM auto-migration.

## Running the Application

### Development Mode

```bash
# Run with live reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Or run directly
go run main.go
```

The server will start on `http://localhost:4000`

### Production Build

```bash
# Build the binary
go build -o cognize .

# Run the binary
./cognize
```

### Using Docker Compose

For a complete development environment with database:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Cognize API server on port 8080

## Project Structure

```
server-cognize/
├── main.go                 # Application entry point
├── .env                    # Environment variables (create this)
├── .air.toml              # Air configuration for live reload
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # Docker image configuration
├── go.mod                 # Go module file
├── go.sum                 # Go dependencies
├── config/                # Configuration management
│   ├── config.go          # Environment configuration
│   ├── dbConfig.go        # Database configuration
│   └── oauthConfig.go     # OAuth configuration
├── internal/              # Internal application packages
│   ├── activity/          # Activity management
│   │   ├── activity.go
│   │   ├── activity_handler.go
│   │   └── activity_service.go
│   ├── card/              # Contact/Card management
│   │   ├── card.go
│   │   ├── card_handler.go
│   │   └── card_service.go
│   ├── field/             # Custom fields
│   ├── keys/              # API key management
│   ├── list/              # List management
│   ├── oauth/             # OAuth handlers
│   ├── tag/               # Tag management
│   └── user/              # User management
├── middleware/            # HTTP middleware
│   └── auth.go           # Authentication middleware
├── models/                # Database models
│   ├── user.go
│   ├── card.go
│   ├── list.go
│   ├── tag.go
│   └── activity.go
├── router/                # Route definitions
│   └── router.go
├── logger/                # Logging configuration
├── db/                    # Database utilities
├── util/                  # Utility functions
└── docs/                  # Documentation
    └── API_REFERENCE.md
```

## Development Workflow

### Code Style

This project follows standard Go conventions:

```bash
# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Run static analysis (optional)
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Building

```bash
# Build for current platform
go build -o cognize .

# Build for Linux (production)
CGO_ENABLED=0 GOOS=linux go build -o cognize .

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o cognize.exe .
```

## Debugging

### Using Delve (Go Debugger)

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug

# Debug specific test
dlv test ./internal/card/
```

### Logging

The application uses Uber Zap for structured logging. Logs are output to:
- Console (development)
- Axiom (production, if configured)

Log levels:
- `DEBUG` - Detailed debugging information
- `INFO` - General information
- `WARN` - Warning messages
- `ERROR` - Error messages

### Database Debugging

To inspect the database:

```bash
# Connect to PostgreSQL
docker exec -it cognize-postgres psql -U root -d cognize

# List tables
\dt

# Describe table structure
\d cards
```

## Common Issues and Solutions

### Port Already in Use

If port 4000 is already in use:

```bash
# Find the process using the port
lsof -i :4000

# Kill the process
kill -9 <PID>

# Or change the port in .env
PORT=4001
```

### Database Connection Issues

1. Ensure PostgreSQL is running:
   ```bash
   docker ps # Check if container is running
   ```

2. Check database credentials in `.env`

3. Verify database exists:
   ```bash
   docker exec -it cognize-postgres psql -U root -c "\l"
   ```

### OAuth Issues

1. Verify Google OAuth credentials are correct
2. Check redirect URI matches exactly (including trailing slashes)
3. Ensure Google+ API is enabled in Google Cloud Console

### Module Issues

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy up dependencies
go mod tidy
```

## Development Tools

### Recommended IDE Setup

**VS Code Extensions:**
- Go (official Go extension)
- REST Client (for testing API endpoints)
- PostgreSQL (for database management)
- Docker (for container management)

**GoLand/IntelliJ IDEA:**
- Built-in Go support
- Database tools
- Docker integration

### Useful Tools

```bash
# Air for live reload
go install github.com/cosmtrek/air@latest

# Database migration tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# API testing
brew install httpie  # or use curl/Postman
```

## Testing API Endpoints

### Using curl

```bash
# Health check
curl http://localhost:4000/

# Get OAuth redirect URL
curl http://localhost:4000/oauth/google/redirect-uri
```

### Using HTTPie

```bash
# Create a card (requires authentication)
http POST localhost:4000/card/create \
  Authorization:"Bearer <token>" \
  name="John Doe" \
  email="john@example.com" \
  list_id:=1
```

### Using Postman

1. Import the API endpoints
2. Set up environment variables for base URL and tokens
3. Create collections for different endpoint groups

## Contributing

1. Create a feature branch from `main`
2. Make your changes following the coding standards
3. Write tests for new functionality
4. Run tests and ensure they pass
5. Submit a pull request

### Commit Message Convention

```
type(scope): description

feat(card): add bulk import functionality
fix(auth): resolve JWT token expiration issue
docs(api): update endpoint documentation
```

## Production Considerations

When preparing for production:

1. **Environment Variables**: Use secure values for all secrets
2. **Database**: Use a managed PostgreSQL service
3. **Logging**: Configure Axiom or another logging service
4. **Monitoring**: Set up health checks and monitoring
5. **Security**: Review authentication and authorization
6. **Performance**: Optimize database queries and add caching if needed

## Getting Help

- Check the [API Reference](./docs/API_REFERENCE.md)
- Review existing issues in the repository
- Contact the development team for specific questions