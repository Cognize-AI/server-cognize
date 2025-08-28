# Cognize API Documentation

Complete API reference for the Cognize CRM server.

## Base URL

```
http://localhost:4000  # Development
https://your-domain.com  # Production
```

## Authentication

Most endpoints require authentication. There are two types:

1. **JWT Token Authentication** - For user-specific operations
2. **API Key Authentication** - For bulk operations

### JWT Authentication

Include the JWT token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### API Key Authentication

Include the API key in the Authorization header:

```http
Authorization: <your-api-key>
```

## Endpoints

### System

#### Health Check

Get system status and welcome message.

```http
GET /
```

**Response:**
```
Welcome to cognize
```

---

### Authentication

#### Get Google OAuth Redirect URL

Get the URL to redirect users to Google for authentication.

```http
GET /oauth/google/redirect-uri
```

**Response:**
```json
{
  "data": {
    "url": "https://accounts.google.com/oauth/authorize?..."
  }
}
```

#### Handle Google OAuth Callback

Process the OAuth callback from Google (used by frontend).

```http
GET /oauth/google/callback?code=<auth-code>&state=<state>
```

**Response:**
Sets a cookie with the JWT token and redirects to the frontend.

---

### User Management

#### Get Current User

Get the profile of the currently authenticated user.

```http
GET /user/me
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "profile_picture": "https://example.com/avatar.jpg"
  }
}
```

---

### List Management

#### Create Default Lists

Create default lists for a new user (Prospects, Contacted, Qualified, etc.).

```http
GET /list/create-default
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "data": {
    "lists": [
      {
        "id": 1,
        "name": "Prospects",
        "color": "#3B82F6",
        "list_order": 1.0,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### Get All Lists

Get all lists belonging to the authenticated user, including their cards.

```http
GET /list/all
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "data": {
    "lists": [
      {
        "id": 1,
        "name": "Prospects",
        "color": "#3B82F6",
        "list_order": 1.0,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "cards": [
          {
            "id": 1,
            "name": "John Smith",
            "designation": "Software Engineer",
            "email": "john.smith@example.com",
            "phone": "+1234567890",
            "image_url": "https://example.com/photo.jpg",
            "tags": [
              {
                "id": 1,
                "name": "Hot Lead",
                "color": "#EF4444"
              }
            ]
          }
        ]
      }
    ]
  }
}
```

---

### Card Management

#### Create Card

Create a new prospect card.

```http
POST /card/create
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Jane Doe",
  "designation": "Product Manager",
  "email": "jane.doe@example.com",
  "phone": "+1987654321",
  "image_url": "https://example.com/jane.jpg",
  "list_id": 1
}
```

**Response:**
```json
{
  "data": {
    "id": 2
  }
}
```

#### Update Card

Update an existing prospect card.

```http
PUT /card/{id}
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Parameters:**
- `id` (path, required): Card ID

**Request Body:**
```json
{
  "name": "Jane Smith",
  "designation": "Senior Product Manager",
  "email": "jane.smith@example.com",
  "phone": "+1987654321",
  "image_url": "https://example.com/jane-updated.jpg"
}
```

**Response:**
```json
{
  "data": "ok"
}
```

#### Delete Card

Delete a prospect card.

```http
DELETE /card/{id}
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
  "data": "ok"
}
```

#### Move Card

Move a card between lists or reorder within a list.

```http
POST /card/move
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "curr_card": 1,
  "list_id": 2,
  "prev_card": 0,
  "next_card": 3
}
```

**Parameters:**
- `curr_card` (required): ID of the card being moved
- `list_id` (required): Target list ID
- `prev_card` (optional): ID of the card that should come before (0 if first)
- `next_card` (optional): ID of the card that should come after (0 if last)

**Response:**
```json
{
  "data": "ok"
}
```

---

### Tag Management

#### Create Tag

Create a new tag.

```http
POST /tag/create
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Hot Lead",
  "color": "#EF4444"
}
```

**Response:**
```json
{
  "data": {
    "id": 1
  }
}
```

#### Get All Tags

Get all tags belonging to the authenticated user.

```http
GET /tag/
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "data": {
    "tags": [
      {
        "id": 1,
        "name": "Hot Lead",
        "color": "#EF4444"
      },
      {
        "id": 2,
        "name": "Cold Lead",
        "color": "#6B7280"
      }
    ]
  }
}
```

#### Edit Tag

Update an existing tag.

```http
PUT /tag/
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "id": 1,
  "name": "Very Hot Lead"
}
```

**Response:**
```json
{
  "data": {
    "id": 1
  }
}
```

#### Delete Tag

Delete a tag.

```http
DELETE /tag/{id}
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Parameters:**
- `id` (path, required): Tag ID

**Response:**
```json
{
  "data": "ok"
}
```

#### Add Tag to Card

Associate a tag with a card.

```http
POST /tag/add-to-card
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "tag_id": 1,
  "card_id": 1
}
```

**Response:**
```json
{
  "data": "ok"
}
```

#### Remove Tag from Card

Remove a tag association from a card.

```http
POST /tag/remove-from-card
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "tag_id": 1,
  "card_id": 1
}
```

**Response:**
```json
{
  "data": "ok"
}
```

---

### API Key Management

#### Generate API Key

Generate a new API key for external integrations.

```http
GET /key/api
```

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "data": {
    "api_key": "ak_1234567890abcdef",
    "name": "API Key 1"
  }
}
```

#### Bulk Create Prospects

Import multiple prospects at once using an API key.

```http
POST /api/bulk-prospect
```

**Headers:**
```http
Authorization: <api-key>
Content-Type: application/json
```

**Request Body:**
```json
{
  "list_id": 1,
  "prospects": [
    {
      "name": "Alice Johnson",
      "designation": "Data Scientist",
      "email": "alice@example.com",
      "phone": "+1111111111",
      "image_url": "https://example.com/alice.jpg"
    },
    {
      "name": "Bob Wilson",
      "designation": "DevOps Engineer",
      "email": "bob@example.com",
      "phone": "+2222222222",
      "image_url": "https://example.com/bob.jpg"
    }
  ]
}
```

**Response:**
```json
{
  "data": {}
}
```

---

## Error Responses

All endpoints return errors in this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### Common HTTP Status Codes

- `200 OK` - Success
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Common Error Messages

- `"Unauthorized"` - Missing or invalid JWT token
- `"Authorization header required"` - Missing Authorization header
- `"card not found for user"` - Card doesn't exist or doesn't belong to user
- `"tag doesnt exists"` - Tag doesn't exist or doesn't belong to user
- `"list not found for card_id: X"` - Associated list not found

---

## Rate Limiting

Currently, there are no rate limits implemented, but this may change in future versions.

## Versioning

The API is currently unversioned. Future versions will include version numbers in the URL path.

## CORS Policy

The API accepts requests from:
- `http://localhost:3000` (development)
- `https://client-cognize.vercel.app` (production frontend)

---

For more information, see the [main documentation](../README.md).