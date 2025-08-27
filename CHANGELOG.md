# Changelog

All notable changes to the Cognize Server project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive documentation including README, API docs, setup guide, deployment guide, and architecture overview
- Contributing guidelines for new developers
- Environment configuration template (.env.example)
- Enhanced .gitignore with comprehensive patterns

### Changed
- Improved project structure documentation
- Updated repository with better developer onboarding experience

## [1.0.0] - Initial Release

### Added
- User authentication with Google OAuth integration
- JWT token-based authorization system
- Prospect card management (CRUD operations)
- Kanban-style list organization system
- Tag system for categorizing prospects
- Card ordering with automatic rebalancing
- API key generation for external integrations
- Bulk prospect import functionality
- Structured logging with Axiom integration
- PostgreSQL database with GORM ORM
- Docker containerization support
- CORS configuration for frontend integration
- Middleware for authentication and request handling

### Features
- **User Management**: Google OAuth authentication, user profile management
- **List Management**: Create default lists, organize prospects in Kanban columns
- **Card Management**: Full CRUD operations, move between lists, automatic ordering
- **Tag System**: Create, edit, delete tags with colors, associate with cards
- **API Keys**: Generate encrypted API keys for external integrations
- **Bulk Operations**: Import multiple prospects via API
- **Security**: JWT tokens, API key encryption, user data isolation
- **Logging**: Structured logging with development and production modes
- **Database**: PostgreSQL with automated migrations and relationship management

### Endpoints
- `GET /` - Health check
- `GET /oauth/google/redirect-uri` - Get Google OAuth URL
- `GET /oauth/google/callback` - Handle OAuth callback
- `GET /user/me` - Get current user profile
- `GET /list/create-default` - Create default lists
- `GET /list/all` - Get all lists with cards
- `POST /card/create` - Create new prospect card
- `PUT /card/:id` - Update prospect card
- `DELETE /card/:id` - Delete prospect card
- `POST /card/move` - Move card between lists
- `POST /tag/create` - Create new tag
- `GET /tag/` - Get all tags
- `PUT /tag/` - Edit tag
- `DELETE /tag/:id` - Delete tag
- `POST /tag/add-to-card` - Add tag to card
- `POST /tag/remove-from-card` - Remove tag from card
- `GET /key/api` - Generate API key
- `POST /api/bulk-prospect` - Bulk import prospects

---

## Release Notes Format

### [Version] - Date

#### Added
- New features and capabilities

#### Changed  
- Changes to existing functionality

#### Deprecated
- Features marked for future removal

#### Removed
- Features that have been removed

#### Fixed
- Bug fixes and corrections

#### Security
- Security-related changes and improvements