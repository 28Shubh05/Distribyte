# Distribyte API Documentation

## Base URL

```http
http://localhost:8080
```

---

# Upload File

### Endpoint

```http
POST /upload
```

### Request

Multipart form-data:

| Key  | Type |
| ---- | ---- |
| file | File |

### Success Response

```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "id": 1,
    "original_name": "resume.pdf",
    "stored_name": "uuid.pdf",
    "filepath": "../storage/uuid.pdf",
    "size": 12345,
    "file_hash": "abc123",
    "uploaded_at": "2025-08-01T10:00:00Z"
  }
}
```

---

# List Files

### Endpoint

```http
GET /files
```

### Success Response

```json
{
  "success": true,
  "files": []
}
```

---

# Download File

### Endpoint

```http
GET /download/:id
```

### Example

```http
GET /download/5
```

Downloads the file.

---

# Soft Delete File

### Endpoint

```http
DELETE /files/:id
```

### Example

```http
DELETE /files/5
```

### Response

```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

---

# Restore File

### Endpoint

```http
POST /restore/:id
```

### Example

```http
POST /restore/5
```

### Response

```json
{
  "success": true,
  "message": "File restored successfully"
}
```

---

# Error Response Format

```json
{
  "success": false,
  "error": "Error message"
}
```
