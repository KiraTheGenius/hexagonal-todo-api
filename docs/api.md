# API Documentation

## Base URL
```
http://localhost:8080
```

## Health Check

### GET /health
Check if the service is running.

**Response:**
```json
{
  "status": "ok"
}
```

## Todo Management

### Create Todo
**POST** `/todo`

**Request Body:**
```json
{
  "description": "Learn hexagonal architecture",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-uuid"
}
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "description": "Learn hexagonal architecture",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-uuid",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

### Get Todo
**GET** `/todo/{id}`

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "description": "Learn hexagonal architecture",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-uuid",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

### List Todos
**GET** `/todo?limit=10&offset=0`

**Query Parameters:**
- `limit` (optional): Number of todos to return (default: 10)
- `offset` (optional): Number of todos to skip (default: 0)

**Response:**
```json
{
  "todos": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "description": "Learn hexagonal architecture",
      "dueDate": "2024-12-31T23:59:59Z",
      "fileId": "optional-file-uuid",
      "createdAt": "2024-01-01T10:00:00Z",
      "updatedAt": "2024-01-01T10:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "count": 1
  }
}
```

### Update Todo
**PUT** `/todo/{id}`

**Request Body:**
```json
{
  "description": "Updated description",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-uuid"
}
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "description": "Updated description",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-file-uuid",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T11:00:00Z"
}
```

### Delete Todo
**DELETE** `/todo/{id}`

**Response:**
```
204 No Content
```

## File Management

### Upload File
**POST** `/upload`

**Request:** `multipart/form-data`
- `file`: The file to upload

**Supported File Types:**
- Images: `.jpg`, `.jpeg`, `.png`, `.gif`
- Documents: `.txt`, `.pdf`, `.doc`, `.docx`

**File Size Limit:** 10MB

**Response:**
```json
{
  "fileId": "123e4567-e89b-12d3-a456-426614174000",
  "url": "https://s3.amazonaws.com/bucket/file-key"
}
```

### Get File Metadata
**GET** `/file/{id}`

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "filename": "document.pdf",
  "contentType": "application/pdf",
  "size": 1024000,
  "url": "https://s3.amazonaws.com/bucket/file-key",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

### Download File
**GET** `/file/{id}/download`

**Response:**
- Content-Type: Based on file type
- Content-Disposition: `attachment; filename="original-filename"`
- Body: File content

### List Files
**GET** `/files?limit=10&offset=0`

**Query Parameters:**
- `limit` (optional): Number of files to return (default: 10)
- `offset` (optional): Number of files to skip (default: 0)

**Response:**
```json
{
  "files": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "filename": "document.pdf",
      "contentType": "application/pdf",
      "size": 1024000,
      "url": "https://s3.amazonaws.com/bucket/file-key",
      "createdAt": "2024-01-01T10:00:00Z",
      "updatedAt": "2024-01-01T10:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "count": 1
  }
}
```

### Update File Metadata
**PUT** `/file/{id}`

**Request Body:**
```json
{
  "filename": "new-filename.pdf",
  "contentType": "application/pdf"
}
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "filename": "new-filename.pdf",
  "contentType": "application/pdf",
  "size": 1024000,
  "url": "https://s3.amazonaws.com/bucket/file-key",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T11:00:00Z"
}
```

### Delete File
**DELETE** `/file/{id}`

**Response:**
```
204 No Content
```

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "Error message",
  "details": "Additional error details (optional)"
}
```

### Common HTTP Status Codes

- `200 OK` - Success
- `201 Created` - Resource created successfully
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Example Usage

### Using curl

```bash
# Create a todo
curl -X POST http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{"description": "Learn Go", "dueDate": "2024-12-31T23:59:59Z"}'

# Upload a file
curl -X POST http://localhost:8080/upload \
  -F "file=@document.pdf"

# Get all todos
curl http://localhost:8080/todo
```

### Using JavaScript

```javascript
// Create a todo
const response = await fetch('http://localhost:8080/todo', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    description: 'Learn Go',
    dueDate: '2024-12-31T23:59:59Z'
  })
});

const todo = await response.json();
console.log(todo);
```