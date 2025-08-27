# Cognize Server

A powerful CRM and prospect management system built with Go, featuring a RESTful API for managing contacts, lists, and tags with Google OAuth authentication.

## 🚀 Features

- **User Management**: Google OAuth integration for secure authentication
- **Prospect Management**: Create, update, delete, and organize prospect cards
- **List Organization**: Kanban-style lists for organizing prospects
- **Tagging System**: Colorful tags for categorizing and filtering prospects
- **Drag & Drop**: Move prospects between lists with automatic ordering
- **Bulk Import**: API for importing multiple prospects at once
- **API Keys**: Generate secure API keys for external integrations
- **Real-time Logging**: Structured logging with Axiom integration

## 🏗️ Architecture

- **Backend**: Go with Gin web framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens with Google OAuth
- **Logging**: Zap structured logging with Axiom
- **Containerization**: Docker and Docker Compose support

## 📋 Requirements

- Go 1.24.6 or higher
- PostgreSQL 17+
- Google OAuth credentials
- (Optional) Axiom account for logging

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/Cognize-AI/server-cognize.git
cd server-cognize
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory:

```bash
# Server Configuration
PORT=4000
ENVIRONMENT=dev

# Database Configuration
DB_STRING=postgres://username:password@localhost:5432/cognize?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here

# Google OAuth Configuration
GOOGLE_OAUTH_CLIENT_ID=your-google-oauth-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback

# Encryption for API Keys
ENC_SECRET=your-32-character-encryption-key

# Optional: Axiom Logging
AXIOM_TOKEN=your-axiom-token
AXIOM_ORG=your-axiom-org
AXIOM_DATASET=your-axiom-dataset

# Optional: SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### 3. Using Docker Compose (Recommended)

```bash
# Start PostgreSQL and the application
docker-compose up -d

# The API will be available at http://localhost:8080
```

### 4. Manual Setup

```bash
# Install dependencies
go mod download

# Start PostgreSQL (using Docker)
docker run --name cognize-postgres \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=cognize \
  -p 5432:5432 -d postgres:17

# Run the application
go run main.go

# The API will be available at http://localhost:4000
```

## 📖 Documentation

- [API Documentation](docs/api.md) - Complete API reference with examples
- [Development Setup](docs/setup.md) - Detailed development environment setup
- [Deployment Guide](docs/deployment.md) - Production deployment instructions
- [Architecture Overview](docs/architecture.md) - Technical architecture details

## 🔌 API Overview

### Authentication Endpoints
- `GET /oauth/google/redirect-uri` - Get Google OAuth URL
- `GET /oauth/google/callback` - Handle OAuth callback

### User Endpoints
- `GET /user/me` - Get current user profile

### List Endpoints
- `GET /list/create-default` - Create default lists
- `GET /list/all` - Get all lists with cards

### Card Endpoints
- `POST /card/create` - Create a new prospect card
- `PUT /card/:id` - Update a prospect card
- `DELETE /card/:id` - Delete a prospect card
- `POST /card/move` - Move card between lists

### Tag Endpoints
- `POST /tag/create` - Create a new tag
- `GET /tag/` - Get all tags
- `PUT /tag/` - Edit a tag
- `DELETE /tag/:id` - Delete a tag
- `POST /tag/add-to-card` - Add tag to card
- `POST /tag/remove-from-card` - Remove tag from card

### API Key Endpoints
- `GET /key/api` - Generate API key
- `POST /api/bulk-prospect` - Bulk import prospects (requires API key)

## 🛠️ Development

### Hot Reload Development

```bash
# Install Air for hot reloading
go install github.com/air-verse/air@latest

# Start development server with hot reload
air
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Code Formatting

```bash
# Format code
go fmt ./...

# Run linter (install golangci-lint first)
golangci-lint run
```

## 🐳 Docker

### Building the Image

```bash
docker build -t cognize-server .
```

### Running with Docker

```bash
docker run -p 4000:4000 --env-file .env cognize-server
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [Cognize Client](https://github.com/Cognize-AI/client-cognize) - React frontend for this API

## 📞 Support

For support, email support@cognize-ai.com or create an issue in this repository.

---

Built with ❤️ by the Cognize AI team