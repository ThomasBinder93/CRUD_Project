# CRUD API Documentation

This is a comprehensive REST API for managing items with a SQLite backend.

## Base URL

```
http://localhost:8080/api
```

## API Endpoints

### Health Check

#### GET /health

Returns the health status of the server.

**Response:** `200 OK`
```json
{
  "status": "healthy"
}
```

---

### Items

#### GET /api/items

Retrieve all items.

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": 1,
      "name": "Item 1"
    },
    {
      "id": 2,
      "name": "Item 2"
    }
  ]
}
```

---

#### POST /api/items

Create a new item.

**Request Body:**
```json
{
  "name": "New Item"
}
```

**Response:** `201 Created`
```json
{
  "data": {
    "id": 3,
    "name": "New Item"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid JSON
- `422 Unprocessable Entity` - Validation failed (empty name, name too long)

---

#### GET /api/items/{id}

Retrieve a specific item by ID.

**Parameters:**
- `id` (path, required): Item ID (integer)

**Response:** `200 OK`
```json
{
  "data": {
    "id": 1,
    "name": "Item 1"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid ID format
- `404 Not Found` - Item not found

---

#### PUT /api/items/{id}

Update an existing item.

**Parameters:**
- `id` (path, required): Item ID (integer)

**Request Body:**
```json
{
  "name": "Updated Item Name"
}
```

**Response:** `200 OK`
```json
{
  "data": {
    "id": 1,
    "name": "Updated Item Name"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid ID format or invalid JSON
- `404 Not Found` - Item not found
- `422 Unprocessable Entity` - Validation failed

---

#### DELETE /api/items/{id}

Delete an item.

**Parameters:**
- `id` (path, required): Item ID (integer)

**Response:** `204 No Content`

**Error Responses:**
- `400 Bad Request` - Invalid ID format

---

## Error Response Format

All errors follow a consistent format:

```json
{
  "error": {
    "code": 400,
    "message": "bad request",
    "details": "detailed error message"
  }
}
```

## Validation Rules

### Item Name
- **Required**: Yes
- **Type**: String
- **Min Length**: 1 character (cannot be empty or whitespace-only)
- **Max Length**: 255 characters

## Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 204 | No Content - Successful deletion |
| 400 | Bad Request - Invalid input or format |
| 404 | Not Found - Resource not found |
| 422 | Unprocessable Entity - Validation error |
| 500 | Internal Server Error - Server error |

## Content-Type

All requests and responses use `application/json` content type.

## Examples

### Create an Item
```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name": "My Item"}'
```

### Get All Items
```bash
curl http://localhost:8080/api/items
```

### Get Item by ID
```bash
curl http://localhost:8080/api/items/1
```

### Update an Item
```bash
curl -X PUT http://localhost:8080/api/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Item"}'
```

### Delete an Item
```bash
curl -X DELETE http://localhost:8080/api/items/1
```

## Rate Limiting

Currently no rate limiting is implemented.

## Authentication

Currently no authentication is required.

## CORS Support

The API supports CORS for cross-origin requests. All origins are allowed.

Allowed Methods: GET, POST, PUT, DELETE, OPTIONS
Allowed Headers: Content-Type, Authorization
