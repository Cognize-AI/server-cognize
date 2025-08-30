# API Reference

## Overview

The Cognize API is a RESTful service for managing contacts, lists, tags, and activities in a CRM system. All endpoints return JSON responses and use standard HTTP status codes.

## Base URL

```
http://localhost:4000
```

## Authentication

### JWT Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```http
Authorization: Bearer <jwt_token>
```

### API Key Authentication

Bulk operations and external integrations use API key authentication:

```http
Cognize-API-Key: <api_key>
```

## Response Format

### Success Response

```json
{
  "data": {
    // Response data
  }
}
```

### Error Response

```json
{
  "error": "Error message"
}
```

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Endpoints

### User Management

#### Get Current User

Get the profile information for the authenticated user.

```http
GET /user/me
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "profilePicture": "https://example.com/profile.jpg"
  }
}
```

### OAuth

#### Get Google OAuth Redirect URL

Get the Google OAuth redirect URL for user authentication.

```http
GET /oauth/google/redirect-uri
```

**Response:**
```json
{
  "data": {
    "url": "https://accounts.google.com/oauth/authorize?client_id=..."
  }
}
```

#### Handle Google OAuth Callback

Process the OAuth callback from Google.

```http
GET /oauth/google/callback?code=<auth_code>
```

**Query Parameters:**
- `code` (string, required) - Authorization code from Google

**Response:**
```json
{
  "data": {
    "token": "jwt_token_here",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

### Lists

#### Create Default Lists

Create default lists for a new user (New Leads, Follow Up, Qualified, Rejected).

```http
GET /list/create-default
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "data": {
    "message": "Default lists created successfully"
  }
}
```

#### Get All Lists

Get all lists for the authenticated user.

```http
GET /list/all
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "New Leads",
      "color": "#F9BA0B",
      "user_id": 1,
      "list_order": 1.0,
      "cards": [
        {
          "id": 1,
          "name": "John Doe",
          "designation": "Software Engineer",
          "email": "john@example.com"
        }
      ]
    }
  ]
}
```

### Cards (Contacts)

#### Create Card

Create a new contact card.

```http
POST /card/create
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "name": "John Doe",
  "designation": "Software Engineer",
  "email": "john@example.com",
  "phone": "+1234567890",
  "image_url": "https://example.com/profile.jpg",
  "list_id": 1
}
```

**Response:**
```json
{
  "data": {
    "id": 123
  }
}
```

#### Get Card by ID

Get detailed information about a specific card.

```http
GET /card/{id}
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Path Parameters:**
- `id` (integer, required) - Card ID

**Response:**
```json
{
  "data": {
    "id": 1,
    "name": "John Doe",
    "designation": "Software Engineer",
    "email": "john@example.com",
    "phone": "+1234567890",
    "image_url": "https://example.com/profile.jpg",
    "location": "San Francisco, CA",
    "list_name": "New Leads",
    "list_color": "#F9BA0B",
    "company": {
      "name": "Tech Corp",
      "role": "Senior Developer",
      "location": "San Francisco, CA",
      "phone": "+1234567890",
      "email": "contact@techcorp.com"
    },
    "tags": [
      {
        "id": 1,
        "name": "frontend",
        "color": "#A78BFA"
      }
    ],
    "activity": [
      {
        "id": 1,
        "content": "Initial contact made",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "additional_contact": [],
    "additional_company": []
  }
}
```

#### Update Card (Basic)

Update basic card information.

```http
PUT /card/{id}
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Path Parameters:**
- `id` (integer, required) - Card ID

**Request Body:**
```json
{
  "name": "John Doe",
  "designation": "Senior Software Engineer",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "image_url": "https://example.com/profile.jpg"
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

#### Update Card Details

Update detailed card information including company details.

```http
PUT /card/details/{id}
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Path Parameters:**
- `id` (integer, required) - Card ID

**Request Body:**
```json
{
  "name": "John Doe",
  "designation": "Senior Software Engineer",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "image_url": "https://example.com/profile.jpg",
  "location": "San Francisco, CA",
  "company_name": "Tech Corp",
  "company_role": "Senior Developer",
  "company_location": "San Francisco, CA",
  "company_phone": "+1234567890",
  "company_email": "contact@techcorp.com"
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

#### Move Card

Move a card between lists or change its position within a list.

```http
POST /card/move
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "prev_card": 0,
  "curr_card": 1,
  "next_card": 2,
  "list_id": 2
}
```

**Response:**
```json
{
  "data": {
    "message": "Card moved successfully"
  }
}
```

#### Delete Card

Delete a contact card.

```http
DELETE /card/{id}
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Path Parameters:**
- `id` (integer, required) - Card ID

**Response:**
```json
{
  "data": {
    "id": 1
  }
}
```

#### Bulk Import Contacts

Import multiple contacts at once using an API key.

```http
POST /api/bulk-prospect
```

**Headers:**
- `Cognize-API-Key: <api_key>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "list_id": 1,
  "prospects": [
    {
      "name": "Jane Smith",
      "designation": "Product Manager",
      "email": "jane@example.com",
      "phone": "+1987654321",
      "image_url": "https://example.com/jane.jpg"
    },
    {
      "name": "Bob Johnson",
      "designation": "Designer",
      "email": "bob@example.com",
      "phone": "+1555666777"
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "message": "Contacts imported successfully"
  }
}
```

### Tags

#### Create Tag

Create a new tag.

```http
POST /tag/create
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "name": "frontend",
  "color": "#A78BFA"
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

Get all tags for the authenticated user.

```http
GET /tag/
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "frontend",
      "color": "#A78BFA"
    },
    {
      "id": 2,
      "name": "backend",
      "color": "#FCA5A5"
    }
  ]
}
```

#### Add Tag to Card

Associate a tag with a card.

```http
POST /tag/add-to-card
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

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
  "data": {
    "message": "Tag added to card successfully"
  }
}
```

#### Remove Tag from Card

Remove a tag association from a card.

```http
POST /tag/remove-from-card
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

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
  "data": {
    "message": "Tag removed from card successfully"
  }
}
```

#### Update Tag

Update an existing tag.

```http
PUT /tag/
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "id": 1,
  "name": "frontend-developer",
  "color": "#9333EA"
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
- `Authorization: Bearer <token>` (required)

**Path Parameters:**
- `id` (integer, required) - Tag ID

**Response:**
```json
{
  "data": {
    "id": 1
  }
}
```

### Activities

#### Create Activity

Create a new activity for a card.

```http
POST /activity/create
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "text": "Called and left voicemail",
  "card_id": 1
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

#### Update Activity

Update an existing activity.

```http
PUT /activity/{id}
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Path Parameters:**
- `id` (integer, required) - Activity ID

**Request Body:**
```json
{
  "text": "Called and scheduled follow-up meeting"
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

#### Delete Activity

Delete an activity.

```http
DELETE /activity/{id}
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Path Parameters:**
- `id` (integer, required) - Activity ID

**Response:**
```json
{
  "data": {
    "message": "Activity deleted successfully"
  }
}
```

### API Keys

#### Generate API Key

Generate a new API key for external integrations.

```http
GET /key/api
```

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "data": {
    "api_key": "your-generated-api-key",
    "message": "API key generated successfully"
  }
}
```

### Custom Fields

#### Create Field Definition

Create a custom field definition.

```http
POST /field/field-definitions
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "name": "LinkedIn Profile",
  "data_type": "url",
  "category": "contact"
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

#### Add Field Value

Add a value for a custom field to a card.

```http
POST /field/field-value
```

**Headers:**
- `Authorization: Bearer <token>` (required)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "field_id": 1,
  "card_id": 1,
  "value": "https://linkedin.com/in/johndoe"
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

## Error Handling

The API uses standard HTTP status codes to indicate success or failure. In case of an error, the response will include an error message:

```json
{
  "error": "Detailed error message"
}
```

### Common Error Codes

- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Rate Limiting

Currently, there are no specific rate limits implemented, but it's recommended to implement reasonable request rates to ensure optimal performance.

## Data Types

### Card Object

```json
{
  "id": 1,
  "name": "string",
  "designation": "string", 
  "email": "string",
  "phone": "string",
  "image_url": "string",
  "location": "string",
  "list_id": 1,
  "card_order": 1.5,
  "company_name": "string",
  "company_role": "string", 
  "company_location": "string",
  "company_phone": "string",
  "company_email": "string"
}
```

### List Object

```json
{
  "id": 1,
  "name": "string",
  "color": "string",
  "user_id": 1,
  "list_order": 1.0
}
```

### Tag Object

```json
{
  "id": 1,
  "name": "string",
  "color": "string",
  "user_id": 1
}
```

### Activity Object

```json
{
  "id": 1,
  "content": "string",
  "card_id": 1,
  "created_at": "2024-01-15T10:30:00Z"
}
```

### User Object

```json
{
  "id": 1,
  "name": "string",
  "email": "string",
  "profile_picture": "string"
}
```