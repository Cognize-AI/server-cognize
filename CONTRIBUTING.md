# Contributing to Cognize Server

Thank you for your interest in contributing to Cognize Server! This guide will help you get started with contributing to our CRM and prospect management system.

## 🚀 Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Set up the development environment** (see [Setup Guide](docs/setup.md))
4. **Create a feature branch** for your changes
5. **Make your changes** and test them thoroughly
6. **Submit a pull request** with a clear description

## 📋 Prerequisites

Before contributing, make sure you have:

- Go 1.24.6 or higher installed
- PostgreSQL 17+ running (Docker recommended)
- Basic understanding of Go, Gin framework, and REST APIs
- Familiarity with Git and GitHub workflows

## 🛠️ Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/server-cognize.git
cd server-cognize

# Add the original repository as upstream
git remote add upstream https://github.com/Cognize-AI/server-cognize.git
```

### 2. Set Up Environment

```bash
# Copy environment template
cp .env.example .env

# Install dependencies
go mod download

# Start PostgreSQL (using Docker)
docker run --name cognize-postgres-dev \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=cognize \
  -p 5432:5432 -d postgres:17

# Start the development server
go run main.go
```

### 3. Verify Setup

```bash
# Test the API
curl http://localhost:4000/
# Should return: "Welcome to cognize"
```

## 📝 Code Style and Standards

### Go Code Style

We follow standard Go conventions:

```go
// ✅ Good: Clear function names and proper error handling
func (s *service) CreateCard(ctx context.Context, req CreateCardReq, user models.User) (*CreateCardResp, error) {
    if req.Name == "" {
        return nil, errors.New("card name is required")
    }
    
    card := models.Card{
        Name:        req.Name,
        Designation: req.Designation,
        Email:       req.Email,
        Phone:       req.Phone,
        ImageURL:    req.ImageURL,
        ListID:      req.ListID,
        CardOrder:   1.0,
    }
    
    if err := s.DB.Create(&card).Error; err != nil {
        return nil, fmt.Errorf("failed to create card: %w", err)
    }
    
    return &CreateCardResp{ID: card.ID}, nil
}
```

### Code Formatting

```bash
# Format code before committing
go fmt ./...

# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

### Naming Conventions

- **Variables**: `camelCase` (e.g., `userEmail`, `cardList`)
- **Constants**: `UPPER_SNAKE_CASE` (e.g., `MAX_RETRY_COUNT`)
- **Types**: `PascalCase` (e.g., `CreateCardReq`, `UserService`)
- **Packages**: `lowercase` (e.g., `card`, `user`, `oauth`)
- **Files**: `snake_case` (e.g., `card_service.go`, `user_handler.go`)

### Error Handling

```go
// ✅ Good: Wrap errors with context
if err := s.DB.Create(&card).Error; err != nil {
    return nil, fmt.Errorf("failed to create card for user %d: %w", user.ID, err)
}

// ✅ Good: Use structured logging
logger.Logger.Error("Failed to create card", 
    zap.Error(err),
    zap.Uint("user_id", user.ID),
    zap.String("card_name", req.Name))
```

## 🏗️ Architecture Patterns

### Service Layer Pattern

Each domain has a service interface and implementation:

```go
// Interface defines the contract
type Service interface {
    CreateCard(ctx context.Context, req CreateCardReq, user models.User) (*CreateCardResp, error)
    UpdateCard(ctx context.Context, req UpdateCardReq, user models.User) (*UpdateCardResp, error)
}

// Implementation contains the business logic
type service struct {
    timeout time.Duration
    DB      *gorm.DB
}

func NewService() Service {
    return &service{
        timeout: time.Duration(20) * time.Second,
        DB:      config.DB,
    }
}
```

### Handler Pattern

HTTP handlers are thin and delegate to services:

```go
func (h *Handler) CreateCard(c *gin.Context) {
    // 1. Extract user from context
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. Parse request
    var req CreateCardReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 3. Delegate to service
    resp, err := h.Service.CreateCard(c, req, user.(models.User))
    if err != nil {
        logger.Logger.Error("Error creating card", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 4. Return response
    c.JSON(http.StatusOK, gin.H{"data": resp})
}
```

### Data Access Pattern

Use GORM best practices:

```go
// ✅ Good: Use preloading to avoid N+1 queries
var cards []models.Card
s.DB.Preload("Tags").Where("list_id = ?", listID).Find(&cards)

// ✅ Good: Validate ownership
var card models.Card
s.DB.Preload("List").Where("id = ?", cardID).First(&card)
if card.List.UserID != user.ID {
    return errors.New("card not found for user")
}

// ✅ Good: Use transactions for complex operations
tx := s.DB.Begin()
if err := tx.Create(&card).Error; err != nil {
    tx.Rollback()
    return err
}
tx.Commit()
```

## 🧪 Testing Guidelines

### Writing Tests

Create tests for your services:

```go
func TestCreateCard(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    service := &service{DB: db}
    
    // Create test user and list
    user := models.User{Name: "Test User", Email: "test@example.com"}
    db.Create(&user)
    
    list := models.List{Name: "Test List", UserID: user.ID}
    db.Create(&list)
    
    // Test card creation
    req := CreateCardReq{
        Name:   "John Doe",
        Email:  "john@example.com",
        ListID: list.ID,
    }
    
    resp, err := service.CreateCard(context.Background(), req, user)
    
    // Verify results
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Greater(t, resp.ID, uint(0))
    
    // Verify in database
    var card models.Card
    db.First(&card, resp.ID)
    assert.Equal(t, "John Doe", card.Name)
    assert.Equal(t, list.ID, card.ListID)
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/card/...

# Run tests with verbose output
go test -v ./...
```

### Test Database Setup

Use a separate test database:

```go
func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    
    // Run migrations
    db.AutoMigrate(&models.User{}, &models.List{}, &models.Card{}, &models.Tag{})
    
    return db
}
```

## 📦 Adding New Features

### 1. Plan Your Feature

Before coding, consider:

- **Purpose**: What problem does this solve?
- **API Design**: What endpoints will you add/modify?
- **Database Changes**: Do you need new models or fields?
- **Breaking Changes**: Will this affect existing functionality?

### 2. Create a Feature Branch

```bash
# Make sure you're on main and up to date
git checkout main
git pull upstream main

# Create and switch to feature branch
git checkout -b feature/your-feature-name
```

### 3. Implement Your Feature

Follow this pattern for new endpoints:

1. **Define types** in the domain package (e.g., `internal/card/card.go`)
2. **Add service method** to the interface and implement it
3. **Add handler method** that uses the service
4. **Register route** in `router/router.go`
5. **Write tests** for your service and handler
6. **Update API documentation**

### 4. Example: Adding a New Endpoint

```go
// 1. Add types to internal/card/card.go
type ArchiveCardReq struct {
    ID uint `uri:"id" binding:"required"`
}

type ArchiveCardResp struct {
    ID uint `json:"id"`
}

// 2. Add to service interface
type Service interface {
    // ... existing methods
    ArchiveCard(ctx context.Context, req ArchiveCardReq, user models.User) (*ArchiveCardResp, error)
}

// 3. Implement service method
func (s *service) ArchiveCard(ctx context.Context, req ArchiveCardReq, user models.User) (*ArchiveCardResp, error) {
    var card models.Card
    s.DB.Preload("List").Where("id = ?", req.ID).First(&card)
    
    if card.List.UserID != user.ID {
        return nil, errors.New("card not found for user")
    }
    
    card.ArchivedAt = time.Now()
    s.DB.Save(&card)
    
    return &ArchiveCardResp{ID: card.ID}, nil
}

// 4. Add handler method
func (h *Handler) ArchiveCard(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var req ArchiveCardReq
    if err := c.ShouldBindUri(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.Service.ArchiveCard(c, req, user.(models.User))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": resp})
}

// 5. Register route in router/router.go
cardRouter.PATCH("/:id/archive", middleware.RequireAuth, cardHandler.ArchiveCard)
```

## 🐛 Bug Fixes

### Finding and Reporting Bugs

1. **Search existing issues** to avoid duplicates
2. **Create detailed bug reports** with:
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details
   - Error logs if applicable

### Fixing Bugs

1. **Write a test** that reproduces the bug
2. **Fix the bug** with minimal changes
3. **Verify the test passes**
4. **Check for regressions** by running the full test suite

Example bug fix:

```go
// Bug: Card ordering calculation was incorrect
// Before (buggy):
func calculateOrder(prev, next float64) float64 {
    return (prev + next) / 2  // This can cause precision issues
}

// After (fixed):
func calculateOrder(prev, next float64) float64 {
    if math.Abs(next-prev) <= 1e-9 {
        // Trigger rebalancing for very small gaps
        return -1  // Special value to indicate rebalancing needed
    }
    return (prev + next) / 2
}
```

## 📖 Documentation

### Code Documentation

Document complex functions:

```go
// RebalanceCards reorders all cards in a list to prevent precision issues
// with decimal ordering. This is called when card orders become too close
// together (within 1e-9) during move operations.
func RebalanceCards(db *gorm.DB, listID uint, userID uint) error {
    // Implementation...
}
```

### API Documentation

When adding new endpoints, update `docs/api.md`:

```markdown
#### Archive Card

Archive a prospect card (soft delete).

```http
PATCH /card/{id}/archive
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Parameters:**
- `id` (path, required): Card ID

**Response:**
```json
{
  "data": {
    "id": 1
  }
}
```
```

## 🔄 Pull Request Process

### Before Submitting

1. **Test your changes** thoroughly
2. **Update documentation** if needed
3. **Check code style** with linter
4. **Write clear commit messages**

### Commit Message Format

Use conventional commits:

```bash
# Feature commits
feat: add card archiving functionality

# Bug fix commits
fix: correct card ordering calculation precision issue

# Documentation commits
docs: update API documentation for new endpoints

# Refactor commits
refactor: simplify card service error handling
```

### Pull Request Template

Use this template for your PR description:

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have tested the changes manually

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
```

### Review Process

1. **Automated checks** must pass (formatting, tests)
2. **Code review** by maintainers
3. **Manual testing** for significant changes
4. **Merge** after approval

## 🎯 Good First Issues

Looking for ideas? Check out issues labeled `good first issue`:

- Fix typos in documentation
- Add input validation to existing endpoints
- Improve error messages
- Add unit tests for existing functions
- Update API documentation examples

## 📚 Resources

### Learning Go and Gin

- [Official Go Tour](https://tour.golang.org/)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)

### Project-Specific Resources

- [Setup Guide](docs/setup.md) - Development environment setup
- [API Documentation](docs/api.md) - Complete API reference
- [Architecture Guide](docs/architecture.md) - Technical architecture
- [Deployment Guide](docs/deployment.md) - Production deployment

## 💬 Getting Help

### Communication Channels

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Code Reviews**: For implementation questions

### Asking Questions

When asking for help:

1. **Be specific** about what you're trying to do
2. **Include relevant code** snippets
3. **Describe what you've tried** already
4. **Mention your environment** (OS, Go version, etc.)

## 🏆 Recognition

Contributors are recognized in:

- **GitHub Contributors** section
- **Release notes** for significant contributions
- **Project README** for major features

## 📄 License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).

---

Thank you for contributing to Cognize Server! Your efforts help make this CRM system better for everyone. 🚀