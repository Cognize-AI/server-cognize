# Development Setup Guide

This guide will help you set up the Cognize server for local development.

## Prerequisites

### Required Software

1. **Go 1.24.6 or higher**
   ```bash
   # Check your Go version
   go version
   
   # Install Go from https://golang.org/doc/install if needed
   ```

2. **PostgreSQL 17+**
   ```bash
   # Option 1: Install locally
   # macOS
   brew install postgresql@17
   brew services start postgresql@17
   
   # Ubuntu/Debian
   sudo apt update
   sudo apt install postgresql-17 postgresql-contrib
   sudo systemctl start postgresql
   
   # Option 2: Use Docker (Recommended)
   docker run --name cognize-postgres \
     -e POSTGRES_USER=root \
     -e POSTGRES_PASSWORD=root \
     -e POSTGRES_DB=cognize \
     -p 5432:5432 -d postgres:17
   ```

3. **Git**
   ```bash
   git --version
   ```

### Optional Development Tools

1. **Air (Hot Reload)**
   ```bash
   go install github.com/air-verse/air@latest
   ```

2. **golangci-lint (Code Linting)**
   ```bash
   # macOS
   brew install golangci-lint
   
   # Linux
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
   
   # Windows
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. **Docker & Docker Compose**
   ```bash
   docker --version
   docker-compose --version
   ```

## Project Setup

### 1. Clone the Repository

```bash
git clone https://github.com/Cognize-AI/server-cognize.git
cd server-cognize
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Configuration

Create a `.env` file in the project root:

```bash
cp .env.example .env  # If example exists, otherwise create manually
```

#### Required Environment Variables

```bash
# .env file
PORT=4000
ENVIRONMENT=dev
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-here-make-it-long-and-random
ENC_SECRET=your-32-character-encryption-key
GOOGLE_OAUTH_CLIENT_ID=your-google-oauth-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback
```

#### Optional Environment Variables

```bash
# Axiom Logging (for production-like logging)
AXIOM_TOKEN=your-axiom-token
AXIOM_ORG=your-axiom-org
AXIOM_DATASET=cognize-logs

# SMTP Configuration (for email features)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### 4. Google OAuth Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Go to "Credentials" and create an OAuth 2.0 Client ID
5. Set the authorized redirect URI to: `http://localhost:3000/auth/callback`
6. Copy the Client ID and Client Secret to your `.env` file

### 5. Database Setup

#### Option 1: Using Docker (Recommended)

```bash
# Start PostgreSQL container
docker run --name cognize-postgres \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=cognize \
  -p 5432:5432 -d postgres:17

# Verify the container is running
docker ps
```

#### Option 2: Local PostgreSQL

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database and user
CREATE DATABASE cognize;
CREATE USER cognize_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE cognize TO cognize_user;
\q

# Update your DB_STRING in .env accordingly
DB_STRING=postgres://cognize_user:your_password@localhost:5432/cognize?sslmode=disable
```

### 6. Verify Setup

```bash
# Test database connection and run migrations
go run main.go

# You should see output like:
# Logger initialized
# DB connection established
# DB sync completed
# [GIN-debug] Listening and serving HTTP on 0.0.0.0:4000
```

## Development Workflow

### Running the Development Server

#### With Hot Reload (Recommended)

```bash
# Start the server with hot reload
air

# Or if air is not in PATH
$(go env GOPATH)/bin/air
```

The server will automatically restart when you make changes to `.go` files.

#### Without Hot Reload

```bash
go run main.go
```

### Code Quality

#### Format Code

```bash
# Format all Go files
go fmt ./...

# Organize imports
goimports -w .
```

#### Run Linter

```bash
# Run golangci-lint
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

#### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/card/...
```

### Database Operations

#### Reset Database

```bash
# Stop the application first, then:
docker exec -it cognize-postgres psql -U root -d cognize -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Restart the application to run migrations
go run main.go
```

#### View Database

```bash
# Connect to the database
docker exec -it cognize-postgres psql -U root -d cognize

# List tables
\dt

# View table structure
\d users
\d cards
\d lists
\d tags
\d keys
```

### Debugging

#### Enable Debug Mode

In your `.env` file:
```bash
ENVIRONMENT=dev
```

This enables:
- Detailed error messages
- SQL query logging
- Debug-level logs

#### Using Go Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Start with debugger
dlv debug main.go

# Set breakpoints and debug interactively
```

#### View Logs

```bash
# View application logs
tail -f build-errors.log

# View container logs
docker logs cognize-postgres
```

## Common Development Tasks

### Adding a New Endpoint

1. Define request/response types in the appropriate package (e.g., `internal/card/card.go`)
2. Add the service method to the interface and implement it
3. Add the handler method
4. Register the route in `router/router.go`
5. Add tests
6. Update API documentation

### Adding Database Migrations

GORM handles automatic migrations, but for complex changes:

1. Create migration files if needed
2. Update model structs in `models/`
3. The migration runs automatically on startup

### Environment-Specific Configuration

#### Development
```bash
ENVIRONMENT=dev
# Uses zap.NewDevelopment() for logging
# Enables debug features
```

#### Production
```bash
ENVIRONMENT=prod
# Uses Axiom logging
# Optimized for performance
```

## Troubleshooting

### Common Issues

#### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check connection string format
DB_STRING=postgres://username:password@host:port/database?sslmode=disable
```

#### Port Already in Use
```bash
# Find process using port 4000
lsof -i :4000

# Kill the process
kill -9 <PID>
```

#### JWT Token Issues
```bash
# Make sure JWT_SECRET is set and long enough
# Minimum recommended length: 32 characters
```

#### Google OAuth Issues
- Verify Client ID and Secret are correct
- Check redirect URI matches exactly
- Ensure Google+ API is enabled

#### Hot Reload Not Working
```bash
# Check if air is installed correctly
air -v

# Verify .air.toml configuration
# Make sure exclude patterns are correct
```

### Getting Help

1. Check the [API documentation](api.md) for endpoint details
2. Look at existing code for patterns and examples
3. Check the GitHub issues for known problems
4. Create a new issue if you find a bug

## Next Steps

- Read the [Architecture Guide](architecture.md) to understand the codebase structure
- Check the [Deployment Guide](deployment.md) for production deployment
- Review the [Contributing Guidelines](../CONTRIBUTING.md) before making changes

---

Happy coding! 🚀