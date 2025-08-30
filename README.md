# Cognize API Server

A Go-based REST API server for a CRM/lead management system that helps users organize and manage their contacts, leads, and business relationships.

## Features

- **User Management**: Authentication via Google OAuth with JWT tokens
- **Contact Management**: Create, update, and organize contacts (cards) with detailed information
- **List Organization**: Organize contacts into customizable lists (like sales pipelines)
- **Tagging System**: Categorize contacts with colored tags
- **Activity Tracking**: Log and track interactions with contacts
- **Custom Fields**: Add custom contact and company information
- **API Key Management**: Generate API keys for external integrations
- **Bulk Operations**: Import multiple contacts via API

## Tech Stack

- **Backend**: Go with Gin web framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens + Google OAuth 2.0
- **Logging**: Uber Zap logger with Axiom integration
- **Deployment**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.24.6 or higher
- PostgreSQL 17+
- Docker & Docker Compose (optional)

### Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/Cognize-AI/server-cognize.git
cd server-cognize
```

2. Create a `.env` file in the root directory:
```env
PORT=4000
DB_STRING=postgres://root:root@localhost:5432/cognize?sslmode=disable
JWT_SECRET=your-jwt-secret-key
GOOGLE_OAUTH_CLIENT_ID=your-google-oauth-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback
ENVIRONMENT=dev
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
AXIOM_TOKEN=your-axiom-token
AXIOM_ORG=your-axiom-org
AXIOM_DATASET=cognize-logs
ENC_SECRET=your-encryption-secret
```

### Running with Docker (Recommended)

1. Start the services:
```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Cognize API server on port 8080

### Running Locally

1. Start PostgreSQL database:
```bash
docker run --name cognize-postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=cognize -p 5432:5432 -d postgres:17
```

2. Install dependencies and run:
```bash
go mod tidy
go run main.go
```

The server will start on `http://localhost:4000`

## API Documentation

### Authentication

The API uses JWT tokens for authentication. Most endpoints require the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

For API integrations, use the `Cognize-API-Key` header:

```
Cognize-API-Key: <your-api-key>
```

### Base URL

```
http://localhost:4000
```

### Core Endpoints

#### User Management
- `GET /user/me` - Get current user profile

#### OAuth
- `GET /oauth/google/redirect-uri` - Get Google OAuth redirect URL
- `GET /oauth/google/callback` - Handle Google OAuth callback

#### Lists
- `GET /list/create-default` - Create default lists for new users
- `GET /list/all` - Get all lists for authenticated user

#### Cards (Contacts)
- `POST /card/create` - Create a new contact
- `GET /card/:id` - Get contact details by ID
- `PUT /card/:id` - Update basic contact information
- `PUT /card/details/:id` - Update detailed contact information
- `DELETE /card/:id` - Delete a contact
- `POST /card/move` - Move contact between lists

#### Tags
- `POST /tag/create` - Create a new tag
- `GET /tag/` - Get all tags
- `POST /tag/add-to-card` - Add tag to a contact
- `POST /tag/remove-from-card` - Remove tag from a contact
- `PUT /tag/` - Update tag
- `DELETE /tag/:id` - Delete tag

#### Activities
- `POST /activity/create` - Create activity for a contact
- `PUT /activity/:id` - Update activity
- `DELETE /activity/:id` - Delete activity

#### API Keys
- `GET /key/api` - Generate API key for external integrations

#### Bulk Operations
- `POST /api/bulk-prospect` - Bulk import contacts (requires API key)

#### Custom Fields
- `POST /field/field-definitions` - Create custom field definitions
- `POST /field/field-value` - Add custom field values

### Example API Calls

#### Create a Contact
```bash
curl -X POST http://localhost:4000/card/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "designation": "Software Engineer",
    "email": "john@example.com",
    "phone": "+1234567890",
    "list_id": 1
  }'
```

#### Get Contact Details
```bash
curl -X GET http://localhost:4000/card/123 \
  -H "Authorization: Bearer <token>"
```

#### Bulk Import Contacts
```bash
curl -X POST http://localhost:4000/api/bulk-prospect \
  -H "Cognize-API-Key: <api-key>" \
  -H "Content-Type: application/json" \
  -d '{
    "list_id": 1,
    "prospects": [
      {
        "name": "Jane Smith",
        "designation": "Product Manager",
        "email": "jane@example.com",
        "phone": "+1987654321"
      }
    ]
  }'
```

## Database Schema

### Core Models

#### User
- `id` (uint) - Primary key
- `name` (string) - User's full name
- `email` (string) - Email address (unique)
- `password` (string) - Hashed password (for future use)
- `profile_picture` (string) - Profile picture URL

#### List
- `id` (uint) - Primary key
- `name` (string) - List name (e.g., "New Leads", "Qualified")
- `color` (string) - Hex color code for UI
- `user_id` (uint) - Foreign key to User
- `list_order` (decimal) - Display order

#### Card (Contact)
- `id` (uint) - Primary key
- `name` (string) - Contact name
- `designation` (string) - Job title
- `email` (string) - Email address
- `phone` (string) - Phone number
- `image_url` (string) - Profile picture URL
- `location` (string) - Personal location
- `list_id` (uint) - Foreign key to List
- `card_order` (decimal) - Position within list
- `company_name` (string) - Company name
- `company_role` (string) - Role at company
- `company_location` (string) - Company location
- `company_phone` (string) - Company phone
- `company_email` (string) - Company email

#### Tag
- `id` (uint) - Primary key
- `name` (string) - Tag name
- `color` (string) - Hex color code
- `user_id` (uint) - Foreign key to User

#### Activity
- `id` (uint) - Primary key
- `content` (string) - Activity description
- `card_id` (uint) - Foreign key to Card

## Development

### Project Structure

```
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   ├── config.go          # Environment configuration
│   ├── dbConfig.go        # Database configuration
│   └── oauthConfig.go     # OAuth configuration
├── internal/              # Internal packages
│   ├── activity/          # Activity management
│   ├── card/              # Contact management
│   ├── field/             # Custom fields
│   ├── keys/              # API key management
│   ├── list/              # List management
│   ├── oauth/             # OAuth handlers
│   ├── tag/               # Tag management
│   └── user/              # User management
├── middleware/            # HTTP middleware
├── models/                # Database models
├── router/                # Route definitions
├── logger/                # Logging configuration
└── util/                  # Utility functions
```

### Building

```bash
# Build for current platform
go build -o cognize .

# Build for Linux (production)
CGO_ENABLED=0 GOOS=linux go build -o cognize .
```

### Testing

```bash
go test ./...
```

### Code Style

This project follows standard Go conventions. Use `gofmt` and `go vet` for code formatting and static analysis.

## Deployment

### Docker Deployment

1. Build and deploy with Docker Compose:
```bash
docker-compose up -d
```

### Manual Deployment

1. Build the binary:
```bash
CGO_ENABLED=0 GOOS=linux go build -o cognize .
```

2. Set up PostgreSQL database

3. Configure environment variables

4. Run the binary:
```bash
./cognize
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | Server port | Yes |
| `DB_STRING` | PostgreSQL connection string | Yes |
| `JWT_SECRET` | JWT signing secret | Yes |
| `GOOGLE_OAUTH_CLIENT_ID` | Google OAuth client ID | Yes |
| `GOOGLE_OAUTH_CLIENT_SECRET` | Google OAuth client secret | Yes |
| `GOOGLE_OAUTH_REDIRECT_URL` | OAuth callback URL | Yes |
| `ENVIRONMENT` | Environment (dev/prod) | Yes |
| `SMTP_HOST` | SMTP server host | No |
| `SMTP_PORT` | SMTP server port | No |
| `SMTP_USERNAME` | SMTP username | No |
| `SMTP_PASSWORD` | SMTP password | No |
| `AXIOM_TOKEN` | Axiom logging token | No |
| `AXIOM_ORG` | Axiom organization | No |
| `AXIOM_DATASET` | Axiom dataset name | No |
| `ENC_SECRET` | Encryption secret | No |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is proprietary software owned by Cognize-AI.

## Support

For support, please contact the development team or create an issue in the repository.