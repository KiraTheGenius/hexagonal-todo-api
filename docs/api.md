# Todo Service API Documentation

## Overview

The Todo Service provides a RESTful API for managing todo items with file upload capabilities, database persistence, and event streaming.

## Base URL

- Development: `http://localhost:8080`
- Production: `https://api.yourdomain.com`

## Authentication

Currently, the API does not require authentication.

## Endpoints

### Health Check

#### GET /health

Check if the service is running.

**Response:**
```json
{
  "status": "ok"
}
```

### File Upload

#### POST /upload

Upload a file to S3 storage.

**Request:**
- Content-Type: `multipart/form-data`
- Body: Form data with `file` field

**Supported file types:**
- Images: `.jpg`, `.jpeg`, `.png`, `.gif`
- Documents: `.txt`, `.pdf`, `.doc`, `.docx`

**File size limit:** 10MB

**Response:**
```json
{
  "fileId": "uuid-of-uploaded-file"
}
```

**Error responses:**
- `400 Bad Request`: Invalid file type or size
- `500 Internal Server Error`: Upload failed

### Todo Management

#### POST /todo

Create a new todo item.

**Request:**
```json
{
  "description": "Complete the assignment",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-id-from-upload"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "description": "Complete the assignment",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-id",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

**Error responses:**
- `400 Bad Request`: Invalid input (missing description, past due date)
- `500 Internal Server Error`: Creation failed

#### GET /todo/:id

Get a specific todo item by ID.

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "description": "Complete the assignment",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-id",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

**Error responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Todo not found

#### GET /todo

List todos with pagination.

**Query parameters:**
- `limit` (optional): Number of items per page (default: 10, max: 100)
- `offset` (optional): Number of items to skip (default: 0)

**Example:** `GET /todo?limit=5&offset=10`

**Response:** `200 OK`
```json
{
  "todos": [
    {
      "id": "uuid",
      "description": "Complete the assignment",
      "dueDate": "2024-12-31T23:59:59Z",
      "fileId": "optional-file-id",
      "createdAt": "2024-01-01T10:00:00Z",
      "updatedAt": "2024-01-01T10:00:00Z"
    }
  ],
  "pagination": {
    "limit": 5,
    "offset": 10,
    "count": 1
  }
}
```

#### PUT /todo/:id

Update an existing todo item.

**Request:**
```json
{
  "description": "Updated description",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "new-file-id"
}
```

All fields are optional. Only provided fields will be updated.

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "description": "Updated description",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "new-file-id",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T11:00:00Z"
}
```

**Error responses:**
- `400 Bad Request`: Invalid UUID format or request body
- `404 Not Found`: Todo not found
- `500 Internal Server Error`: Update failed

#### DELETE /todo/:id

Delete a todo item.

**Response:** `204 No Content`

**Error responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Todo not found
- `500 Internal Server Error`: Deletion failed

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": "Human readable error message",
  "details": "Additional error details (optional)"
}
```

## Rate Limiting

Currently, no rate limiting is implemented. Consider adding rate limiting for production use.

## CORS

The API supports CORS with the following configuration:
- All origins allowed (`*`)
- Methods: GET, POST, PUT, DELETE, OPTIONS
- Headers: Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization

## Request/Response Headers

### Request Headers
- `Content-Type`: application/json (for JSON endpoints)
- `X-Request-ID`: Optional request ID for tracing

### Response Headers
- `X-Request-ID`: Request ID for tracing
- `Content-Type`: application/json

## Examples

### Complete workflow example:

1. **Upload a file:**
```bash
curl -X POST -F "file=@document.pdf" http://localhost:8080/upload
```

2. **Create a todo with the file:**
```bash
curl -X POST http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Review uploaded document",
    "dueDate": "2024-12-31T23:59:59Z",
    "fileId": "file-id-from-step-1"
  }'
```

3. **List todos:**
```bash
curl http://localhost:8080/todo?limit=5
```

4. **Update a todo:**
```bash
curl -X PUT http://localhost:8080/todo/todo-id \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description"
  }'
```

5. **Delete a todo:**
```bash
curl -X DELETE http://localhost:8080/todo/todo-id
``` 