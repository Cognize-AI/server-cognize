# Architecture Overview

This document provides a comprehensive overview of the Cognize server architecture, design patterns, and technical decisions.

## System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (React)       │◄──►│   (Go/Gin)      │◄──►│   (PostgreSQL)  │
│                 │    │                 │    │                 │
│ - Client UI     │    │ - REST API      │    │ - User Data     │
│ - Auth Tokens   │    │ - Business Logic│    │ - Cards/Lists   │
│ - State Mgmt    │    │ - Data Access   │    │ - Tags/Keys     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   External      │
                       │   Services      │
                       │                 │
                       │ - Google OAuth  │
                       │ - Axiom Logs    │
                       │ - SMTP Email    │
                       └─────────────────┘
```

## Technology Stack

### Backend Framework
- **Go 1.24.6**: Modern, fast, and efficient language
- **Gin**: Lightweight HTTP web framework with excellent performance
- **GORM**: Object-relational mapping library for database operations

### Database
- **PostgreSQL 17**: Robust relational database with excellent JSON support
- **GORM Migrations**: Automatic schema management and versioning

### Authentication & Security
- **JWT Tokens**: Stateless authentication with configurable expiration
- **Google OAuth 2.0**: Secure third-party authentication
- **AES Encryption**: API key encryption for secure storage
- **CORS**: Configurable cross-origin resource sharing

### Logging & Monitoring
- **Zap**: High-performance structured logging
- **Axiom**: Cloud-based log aggregation and analysis
- **Console Logging**: Development-friendly local logging

### Development Tools
- **Air**: Hot reloading for development
- **Docker**: Containerization for consistent environments
- **Viper**: Configuration management with environment variables

## Project Structure

```
server-cognize/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── Dockerfile              # Container image definition
├── docker-compose.yml      # Local development setup
├── .air.toml              # Hot reload configuration
├── .env                   # Environment variables (not in git)
├── .gitignore             # Git ignore patterns
│
├── config/                # Configuration management
│   ├── config.go          # Main config struct and loader
│   ├── dbConfig.go        # Database connection setup
│   └── oauthConfig.go     # OAuth configuration
│
├── middleware/            # HTTP middleware
│   └── auth.go           # JWT and API key authentication
│
├── router/               # HTTP routing
│   └── router.go         # Route definitions and CORS setup
│
├── logger/               # Logging configuration
│   └── logger.go         # Zap logger setup with Axiom
│
├── models/               # Database models
│   ├── user.go           # User model with OAuth profile
│   ├── card.go           # Prospect card model
│   ├── list.go           # List/column model for organization
│   ├── tag.go            # Tag model for categorization
│   └── keys.go           # API key model with encryption
│
├── internal/             # Business logic modules
│   ├── user/             # User management
│   │   ├── user.go       # Types and interface
│   │   ├── user_handler.go  # HTTP handlers
│   │   └── user_service.go  # Business logic
│   │
│   ├── oauth/            # OAuth authentication
│   │   ├── oauth.go      # Types and interface
│   │   ├── oauth_handler.go # OAuth flow handlers
│   │   └── oauth_service.go # OAuth business logic
│   │
│   ├── list/             # List management
│   │   ├── list.go       # Types and interface
│   │   ├── list_handler.go  # HTTP handlers
│   │   └── list_service.go  # Business logic
│   │
│   ├── card/             # Card/prospect management
│   │   ├── card.go       # Types and interface
│   │   ├── card_handler.go  # HTTP handlers
│   │   └── card_service.go  # Business logic with ordering
│   │
│   ├── tag/              # Tag management
│   │   ├── tag.go        # Types and interface
│   │   ├── tag_handler.go   # HTTP handlers
│   │   └── tag_service.go   # Business logic
│   │
│   └── keys/             # API key management
│       ├── keys.go       # Types and interface
│       ├── keys_handler.go  # HTTP handlers
│       └── keys_service.go  # Business logic
│
├── util/                 # Utility functions
│   └── ...              # Helper functions
│
└── docs/                # Documentation
    ├── api.md           # API documentation
    ├── setup.md         # Development setup
    ├── deployment.md    # Deployment guide
    └── architecture.md  # This file
```

## Design Patterns

### Clean Architecture

The application follows clean architecture principles:

1. **Separation of Concerns**: Each layer has distinct responsibilities
2. **Dependency Inversion**: Interfaces define contracts, implementations are injected
3. **Testability**: Business logic is isolated from framework concerns

### Layer Structure

```
┌─────────────────────────────────────────────────────────────┐
│                    Handlers (HTTP Layer)                    │
│  - HTTP request/response handling                           │
│  - Input validation and serialization                      │
│  - Authentication middleware integration                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Services (Business Layer)                  │
│  - Core business logic                                      │
│  - Data validation and transformation                      │
│  - Authorization and access control                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Models (Data Layer)                       │
│  - Database schema definitions                              │
│  - Data access patterns                                     │
│  - Relationship management                                  │
└─────────────────────────────────────────────────────────────┘
```

### Repository Pattern

Each service uses GORM directly, but follows repository patterns:

```go
// Service interface defines business operations
type Service interface {
    CreateCard(ctx context.Context, req CreateCardReq, user models.User) (*CreateCardResp, error)
    UpdateCard(ctx context.Context, req UpdateCardReq, user models.User) (*UpdateCardResp, error)
    // ... other methods
}

// Service implementation handles business logic
type service struct {
    timeout time.Duration
    DB      *gorm.DB
}
```

### Dependency Injection

Dependencies are injected through constructors:

```go
// In main.go
userSvc := user.NewService()
cardSvc := card.NewService()
userHandler := user.NewHandler(userSvc)
cardHandler := card.NewHandler(cardSvc)
```

## Data Models

### Core Entities

#### User
```go
type User struct {
    gorm.Model
    Name           string
    Email          string `gorm:"uniqueIndex"`
    Password       string
    ProfilePicture string
    Lists          []List `gorm:"foreignKey:UserID;references:ID"`
}
```

#### List (Kanban Columns)
```go
type List struct {
    gorm.Model
    Name      string
    Color     string
    UserID    uint    `gorm:"index"`
    ListOrder float64 `gorm:"type:decimal(20,10);index"`
    User      User    `gorm:"foreignKey:UserID;references:ID"`
    Cards     []Card  `gorm:"foreignKey:ListID;references:ID"`
}
```

#### Card (Prospects)
```go
type Card struct {
    gorm.Model
    Name        string  `gorm:"index"`
    Designation string
    Email       string  `gorm:"index"`
    Phone       string
    ImageURL    string
    ListID      uint    `gorm:"index"`
    CardOrder   float64 `gorm:"type:decimal(20,10);index"`
    List        List    `gorm:"foreignKey:ListID;references:ID"`
    Tags        []Tag   `gorm:"many2many:card_tags;"`
}
```

#### Tag
```go
type Tag struct {
    gorm.Model
    Name   string
    Color  string
    UserID uint   `gorm:"index"`
    User   User   `gorm:"foreignKey:UserID;references:ID"`
    Cards  []Card `gorm:"many2many:card_tags;"`
}
```

#### Key (API Keys)
```go
type Key struct {
    gorm.Model
    Name   string
    Value  string `gorm:"unique"`  // Encrypted
    Hash   string                  // For lookup
    UserID uint   `gorm:"index"`
    User   User   `gorm:"foreignKey:UserID;references:ID"`
}
```

### Relationships

- **User → Lists**: One-to-Many (User owns multiple lists)
- **List → Cards**: One-to-Many (List contains multiple cards)
- **User → Tags**: One-to-Many (User owns multiple tags)
- **Card ↔ Tags**: Many-to-Many (Cards can have multiple tags)
- **User → Keys**: One-to-Many (User can have multiple API keys)

## Authentication Flow

### OAuth Flow

```
1. Frontend                2. Backend               3. Google OAuth
   │                          │                        │
   │ GET /oauth/google/       │                        │
   │ redirect-uri             │                        │
   │─────────────────────────►│                        │
   │                          │                        │
   │ {url: "google-auth-url"} │                        │
   │◄─────────────────────────│                        │
   │                          │                        │
   │ User clicks URL          │                        │
   │─────────────────────────────────────────────────►│
   │                          │                        │
   │                     User authorizes                │
   │◄─────────────────────────────────────────────────│
   │                          │                        │
   │ GET /oauth/google/       │                        │
   │ callback?code=xyz        │                        │
   │─────────────────────────►│                        │
   │                          │ Exchange code for      │
   │                          │ user profile           │
   │                          │───────────────────────►│
   │                          │◄───────────────────────│
   │                          │ Create/update user     │
   │                          │ Generate JWT           │
   │ Set cookie & redirect    │                        │
   │◄─────────────────────────│                        │
```

### JWT Authentication

```go
// Middleware validates JWT tokens
func RequireAuth(c *gin.Context) {
    // 1. Extract token from cookie or Authorization header
    // 2. Validate JWT signature and expiration
    // 3. Load user from database
    // 4. Set user in context for handlers
}
```

### API Key Authentication

```go
// Middleware validates API keys for bulk operations
func RequireAPIKey(c *gin.Context) {
    // 1. Extract API key from Authorization header
    // 2. Hash the key for database lookup
    // 3. Load associated user
    // 4. Set user in context
}
```

## Business Logic

### Card Ordering System

Cards within lists use a decimal ordering system:

```go
// Cards are ordered by CardOrder field (decimal)
// When moving cards, new order is calculated between neighbors
func calculateNewOrder(prevOrder, nextOrder float64) float64 {
    if prevOrder == 0 && nextOrder == 0 {
        return 1.0  // First card
    }
    if prevOrder == 0 {
        return nextOrder / 2  // Insert at beginning
    }
    if nextOrder == 0 {
        return prevOrder + 1  // Insert at end
    }
    return (prevOrder + nextOrder) / 2  // Insert between
}
```

### Rebalancing System

When orders get too close (< 1e-9), the system rebalances:

```go
func RebalanceCards(db *gorm.DB, listID uint, userID uint) error {
    var cards []models.Card
    db.Where("list_id = ?", listID).Order("card_order ASC").Find(&cards)
    
    for i := range cards {
        cards[i].CardOrder = float64(i + 1)
    }
    
    return db.Save(&cards).Error
}
```

### Data Access Patterns

#### User Ownership Validation

All operations validate user ownership:

```go
// Example from card service
func (s *service) UpdateCard(ctx context.Context, req UpdateCardReq, user models.User) (*UpdateCardResp, error) {
    var card models.Card
    s.DB.Preload("List").Where("id = ?", req.ID).First(&card)
    
    // Validate ownership through list relationship
    if card.List.UserID != user.ID {
        return nil, errors.New("card not found for user")
    }
    
    // Proceed with update
}
```

#### Preloading Relationships

GORM preloading is used to avoid N+1 queries:

```go
// Load lists with their cards and tags
s.DB.Preload("Cards.Tags").Where("user_id = ?", user.ID).Find(&lists)
```

## Security Considerations

### Data Protection

1. **API Key Encryption**: Keys are encrypted with AES before storage
2. **Password Hashing**: User passwords are hashed (OAuth users don't have passwords)
3. **SQL Injection Prevention**: GORM provides SQL injection protection
4. **CORS Configuration**: Restricts frontend domains

### Access Control

1. **User Isolation**: All operations validate user ownership
2. **JWT Expiration**: Tokens have configurable expiration times
3. **API Key Management**: Users can generate/revoke their own keys

### Environment Security

1. **Secrets Management**: All sensitive data in environment variables
2. **Database SSL**: Production databases use SSL connections
3. **HTTPS Only**: Production deployments require HTTPS

## Performance Considerations

### Database Optimization

1. **Indexes**: Strategic indexes on foreign keys and lookup fields
2. **Connection Pooling**: GORM manages database connection pools
3. **Preloading**: Efficient relationship loading to prevent N+1 queries

### Caching Strategy

Currently no caching is implemented, but could be added:

1. **Redis Session Store**: For JWT token blacklisting
2. **Query Result Caching**: For frequently accessed data
3. **CDN**: For static assets and images

### Monitoring

1. **Structured Logging**: All operations are logged with context
2. **Health Checks**: Basic health endpoint for monitoring
3. **Error Tracking**: Errors are logged with stack traces

## Configuration Management

### Environment Variables

Configuration is managed through environment variables using Viper:

```go
type Config struct {
    PORT                    string
    DbString                string
    JwtSecret               string
    GoogleOAuthClientID     string
    GoogleOAuthClientSecret string
    // ... other fields
}
```

### Development vs Production

- **Development**: Console logging, debug features enabled
- **Production**: Axiom logging, optimized performance

## Testing Strategy

### Unit Testing

Each service should have unit tests:

```go
func TestCreateCard(t *testing.T) {
    // Setup test database
    // Create test user
    // Test card creation
    // Verify results
}
```

### Integration Testing

Test complete workflows:

```go
func TestCardMoveBetweenLists(t *testing.T) {
    // Create user, lists, and cards
    // Move card between lists
    // Verify ordering and relationships
}
```

### API Testing

Test HTTP endpoints:

```go
func TestCreateCardEndpoint(t *testing.T) {
    // Setup test server
    // Make HTTP request
    // Verify response
}
```

## Future Enhancements

### Scalability

1. **Database Sharding**: Partition data by user
2. **Read Replicas**: Separate read/write operations
3. **Caching Layer**: Redis for frequent queries
4. **Load Balancing**: Multiple application instances

### Features

1. **Real-time Updates**: WebSocket support for live collaboration
2. **File Uploads**: Image storage for prospect photos
3. **Bulk Operations**: More efficient bulk data processing
4. **Advanced Search**: Full-text search with Elasticsearch
5. **Analytics**: Usage metrics and reporting

### DevOps

1. **CI/CD Pipeline**: Automated testing and deployment
2. **Infrastructure as Code**: Terraform for cloud resources
3. **Monitoring**: Comprehensive application monitoring
4. **Backup Strategy**: Automated database backups

---

This architecture provides a solid foundation for a scalable CRM system while maintaining simplicity and developer productivity.