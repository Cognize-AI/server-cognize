# Configuration Guide

This document provides detailed information about configuring the Cognize API server for different environments.

## Table of Contents

- [Environment Variables](#environment-variables)
- [Configuration Files](#configuration-files)
- [Database Configuration](#database-configuration)
- [Authentication Configuration](#authentication-configuration)
- [Logging Configuration](#logging-configuration)
- [Email Configuration](#email-configuration)
- [Development vs Production](#development-vs-production)
- [Security Best Practices](#security-best-practices)

## Environment Variables

The application uses environment variables for configuration management via the [Viper](https://github.com/spf13/viper) library. Variables can be set in a `.env` file or as system environment variables.

### Core Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `PORT` | string | Yes | - | Server port number |
| `ENVIRONMENT` | string | Yes | - | Environment mode (`dev`, `staging`, `prod`) |

#### Example
```env
PORT=4000
ENVIRONMENT=dev
```

### Database Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `DB_STRING` | string | Yes | - | PostgreSQL connection string |

#### Example
```env
# Development
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable

# Production
DB_STRING=postgres://username:password@db-host:5432/cognize?sslmode=require
```

#### Connection String Format
```
postgres://[username[:password]@][host[:port]]/database[?options]
```

**Common Options:**
- `sslmode=disable` - Disable SSL (development only)
- `sslmode=require` - Require SSL (production)
- `connect_timeout=10` - Connection timeout in seconds
- `pool_max_conns=25` - Maximum connections in pool

### Authentication Configuration

#### JWT Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `JWT_SECRET` | string | Yes | - | Secret key for JWT token signing |

#### Google OAuth Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `GOOGLE_OAUTH_CLIENT_ID` | string | Yes | - | Google OAuth client ID |
| `GOOGLE_OAUTH_CLIENT_SECRET` | string | Yes | - | Google OAuth client secret |
| `GOOGLE_OAUTH_REDIRECT_URL` | string | Yes | - | OAuth callback URL |

#### Example
```env
JWT_SECRET=your-super-secret-jwt-signing-key-256-bits-minimum
GOOGLE_OAUTH_CLIENT_ID=123456789-abcdef.apps.googleusercontent.com
GOOGLE_OAUTH_CLIENT_SECRET=your-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback
```

#### JWT Security Requirements
- **Minimum length**: 256 bits (32 characters)
- **Character set**: Use alphanumeric + special characters
- **Uniqueness**: Generate a unique secret for each environment
- **Rotation**: Rotate secrets regularly in production

### Email Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `SMTP_HOST` | string | No | - | SMTP server hostname |
| `SMTP_PORT` | int | No | - | SMTP server port |
| `SMTP_USERNAME` | string | No | - | SMTP authentication username |
| `SMTP_PASSWORD` | string | No | - | SMTP authentication password |

#### Example
```env
# Gmail SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=noreply@yourdomain.com
SMTP_PASSWORD=your-app-password

# AWS SES SMTP
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USERNAME=your-ses-smtp-username
SMTP_PASSWORD=your-ses-smtp-password
```

### Logging Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `AXIOM_TOKEN` | string | No | - | Axiom API token for log ingestion |
| `AXIOM_ORG` | string | No | - | Axiom organization ID |
| `AXIOM_DATASET` | string | No | - | Axiom dataset name for logs |

#### Example
```env
AXIOM_TOKEN=xaat-12345678-1234-1234-1234-123456789abc
AXIOM_ORG=your-org-name
AXIOM_DATASET=cognize-production-logs
```

### Encryption Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `ENC_SECRET` | string | No | - | Secret key for additional encryption operations |

#### Example
```env
ENC_SECRET=another-256-bit-secret-for-encryption-operations
```

## Configuration Files

### .env File Structure

Create a `.env` file in the root directory:

```env
# ======================
# Server Configuration
# ======================
PORT=4000
ENVIRONMENT=dev

# ======================
# Database Configuration
# ======================
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable

# ======================
# Authentication
# ======================
JWT_SECRET=your-jwt-secret-key-must-be-long-and-secure-256-bits
GOOGLE_OAUTH_CLIENT_ID=your-google-oauth-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback

# ======================
# Email Configuration (Optional)
# ======================
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# ======================
# Logging Configuration (Optional)
# ======================
AXIOM_TOKEN=your-axiom-token
AXIOM_ORG=your-axiom-org
AXIOM_DATASET=cognize-logs

# ======================
# Encryption (Optional)
# ======================
ENC_SECRET=your-encryption-secret-key
```

### Environment-Specific Files

You can create different configuration files for different environments:

```
.env.development
.env.staging
.env.production
```

Load specific environment files:

```bash
# Development
cp .env.development .env

# Staging
cp .env.staging .env

# Production
cp .env.production .env
```

## Database Configuration

### Connection Pool Settings

The application uses GORM with PostgreSQL. Connection pool settings are configured in `config/dbConfig.go`:

```go
// Example configuration
db.DB().SetMaxOpenConns(25)    // Maximum open connections
db.DB().SetMaxIdleConns(25)    // Maximum idle connections
db.DB().SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
```

### Database URL Examples

#### Local Development
```env
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable
```

#### Docker Compose
```env
DB_STRING=postgres://root:root@postgres:5432/cognize?sslmode=disable
```

#### AWS RDS
```env
DB_STRING=postgres://username:password@cognize-db.cluster-xyz.us-east-1.rds.amazonaws.com:5432/cognize?sslmode=require
```

#### Google Cloud SQL
```env
DB_STRING=postgres://username:password@10.0.0.1:5432/cognize?sslmode=require
```

#### Azure Database
```env
DB_STRING=postgres://username@servername:password@servername.postgres.database.azure.com:5432/cognize?sslmode=require
```

### SSL Configuration

For production databases, always use SSL:

```env
# Require SSL
DB_STRING=postgres://user:pass@host:5432/db?sslmode=require

# Verify CA certificate
DB_STRING=postgres://user:pass@host:5432/db?sslmode=verify-ca

# Verify full certificate
DB_STRING=postgres://user:pass@host:5432/db?sslmode=verify-full
```

## Authentication Configuration

### JWT Configuration

#### Development
```env
JWT_SECRET=dev-secret-key-for-development-only-not-secure
```

#### Production
Generate a secure secret:

```bash
# Generate a secure random secret
openssl rand -hex 32

# Or using Go
go run -c "package main; import (\"crypto/rand\"; \"encoding/hex\"; \"fmt\"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(hex.EncodeToString(b)) }"
```

```env
JWT_SECRET=a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456
```

### Google OAuth Setup

1. **Create Google Cloud Project:**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create new project or select existing one

2. **Enable Google+ API:**
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google+ API" and enable it

3. **Create OAuth Credentials:**
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth 2.0 Client IDs"
   - Application type: "Web application"
   - Add authorized redirect URIs

4. **Configure Redirect URIs:**
   ```
   # Development
   http://localhost:3000/auth/callback
   
   # Production
   https://yourdomain.com/auth/callback
   ```

5. **Environment Configuration:**
   ```env
   GOOGLE_OAUTH_CLIENT_ID=123456789-abcdefghijklmnop.apps.googleusercontent.com
   GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret-from-google
   GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback
   ```

## Logging Configuration

### Console Logging (Development)

In development, logs are output to the console with colored formatting.

### Axiom Integration (Production)

[Axiom](https://axiom.co/) provides centralized logging for production:

1. **Create Axiom Account:**
   - Sign up at [axiom.co](https://axiom.co/)
   - Create organization and dataset

2. **Generate API Token:**
   - Go to Settings > API Tokens
   - Create new token with appropriate permissions

3. **Configure Environment:**
   ```env
   AXIOM_TOKEN=xaat-12345678-1234-1234-1234-123456789abc
   AXIOM_ORG=your-organization-name
   AXIOM_DATASET=cognize-production-logs
   ```

### Log Levels

Configure log levels based on environment:

- **Development**: `DEBUG` level for detailed debugging
- **Staging**: `INFO` level for general information
- **Production**: `WARN` level for warnings and errors only

### Custom Logging Providers

To use different logging providers, modify `logger/logger.go`:

```go
// Example: ELK Stack integration
// Example: CloudWatch Logs integration
// Example: Fluentd integration
```

## Email Configuration

### Gmail SMTP

1. **Enable 2-Factor Authentication** on your Google account

2. **Generate App Password:**
   - Google Account settings > Security > 2-Step Verification
   - App passwords > Generate password for "Cognize API"

3. **Configuration:**
   ```env
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-16-character-app-password
   ```

### AWS SES

1. **Set up AWS SES:**
   - Create SES identity (domain or email)
   - Verify identity
   - Create SMTP credentials

2. **Configuration:**
   ```env
   SMTP_HOST=email-smtp.us-east-1.amazonaws.com
   SMTP_PORT=587
   SMTP_USERNAME=your-ses-smtp-username
   SMTP_PASSWORD=your-ses-smtp-password
   ```

### SendGrid

1. **Create SendGrid account and API key**

2. **Configuration:**
   ```env
   SMTP_HOST=smtp.sendgrid.net
   SMTP_PORT=587
   SMTP_USERNAME=apikey
   SMTP_PASSWORD=your-sendgrid-api-key
   ```

## Development vs Production

### Development Configuration

```env
# Development .env
PORT=4000
ENVIRONMENT=dev
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable
JWT_SECRET=dev-secret-not-secure
GOOGLE_OAUTH_CLIENT_ID=dev-client-id
GOOGLE_OAUTH_CLIENT_SECRET=dev-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback
```

### Production Configuration

```env
# Production .env
PORT=4000
ENVIRONMENT=prod
DB_STRING=postgres://user:secure_pass@prod-db:5432/cognize?sslmode=require
JWT_SECRET=very-long-secure-random-production-secret-256-bits
GOOGLE_OAUTH_CLIENT_ID=prod-client-id.apps.googleusercontent.com
GOOGLE_OAUTH_CLIENT_SECRET=prod-client-secret
GOOGLE_OAUTH_REDIRECT_URL=https://api.yourdomain.com/auth/callback
AXIOM_TOKEN=prod-axiom-token
AXIOM_ORG=your-org
AXIOM_DATASET=cognize-prod-logs
```

### Configuration Validation

The application validates configuration on startup:

```go
// Example validation in config/config.go
func (c *Config) Validate() error {
    if c.PORT == "" {
        return errors.New("PORT is required")
    }
    if c.JwtSecret == "" {
        return errors.New("JWT_SECRET is required")
    }
    if len(c.JwtSecret) < 32 {
        return errors.New("JWT_SECRET must be at least 32 characters")
    }
    return nil
}
```

## Security Best Practices

### Secret Management

1. **Never commit secrets to version control:**
   ```bash
   # Add to .gitignore
   .env
   .env.*
   *.key
   *.pem
   ```

2. **Use environment-specific secrets:**
   - Different secrets for dev/staging/prod
   - Rotate secrets regularly
   - Use secret management services (AWS Secrets Manager, Azure Key Vault, etc.)

3. **Minimum secret requirements:**
   - JWT secrets: minimum 256 bits (32 characters)
   - Use cryptographically secure random generation
   - Include uppercase, lowercase, numbers, and symbols

### Access Control

1. **Database access:**
   - Use dedicated database users for applications
   - Grant minimal required permissions
   - Use connection pooling and timeouts

2. **API access:**
   - Implement rate limiting
   - Use HTTPS in production
   - Validate all inputs

### Environment Isolation

1. **Network isolation:**
   - Use VPCs/VNets for production
   - Separate environments completely
   - Restrict database access to application servers only

2. **Credential isolation:**
   - Different credentials for each environment
   - No production credentials in development
   - Use managed services for secret storage

## Configuration Troubleshooting

### Common Issues

1. **Missing Environment Variables:**
   ```bash
   # Check if variable is set
   echo $JWT_SECRET
   
   # Application startup will show validation errors
   2024/01/15 10:30:00 Error loading config: JWT_SECRET is required
   ```

2. **Database Connection Issues:**
   ```bash
   # Test database connectivity
   psql "postgres://user:pass@host:5432/db" -c "SELECT 1;"
   
   # Check application logs
   2024/01/15 10:30:00 Error connecting to database: connection refused
   ```

3. **OAuth Configuration Issues:**
   ```bash
   # Check redirect URI matches exactly
   # Verify client ID and secret
   # Ensure Google+ API is enabled
   ```

### Validation Tools

Create a configuration validation script:

```bash
#!/bin/bash
# validate-config.sh

echo "Validating configuration..."

# Check required variables
required_vars=("PORT" "DB_STRING" "JWT_SECRET" "GOOGLE_OAUTH_CLIENT_ID")

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "ERROR: $var is not set"
        exit 1
    fi
done

# Check JWT secret length
if [ ${#JWT_SECRET} -lt 32 ]; then
    echo "ERROR: JWT_SECRET must be at least 32 characters"
    exit 1
fi

# Test database connection
if ! pg_isready -d "$DB_STRING" >/dev/null 2>&1; then
    echo "ERROR: Cannot connect to database"
    exit 1
fi

echo "Configuration validation passed!"
```

### Configuration Templates

Create configuration templates for different environments:

```bash
# scripts/create-env.sh
#!/bin/bash

ENVIRONMENT=${1:-dev}

cat > .env << EOF
# Generated configuration for $ENVIRONMENT
PORT=4000
ENVIRONMENT=$ENVIRONMENT
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable
JWT_SECRET=$(openssl rand -hex 32)
GOOGLE_OAUTH_CLIENT_ID=your-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback
EOF

echo "Configuration file created for $ENVIRONMENT environment"
```

This configuration guide provides comprehensive information for setting up the Cognize API server in various environments with proper security considerations.