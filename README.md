# Distribyte

A modern Distributed Object Storage System built with Go, PostgreSQL, Redis, React, and TailwindCSS.

Distribyte allows users to securely upload, store, download, restore, and manage files while providing authentication, caching, duplicate detection, and a modern dashboard interface.

---

## Features

### Authentication & Security

* User Registration
* User Login
* JWT Authentication
* Password Hashing using bcrypt
* Protected API Routes
* User-specific File Access

### File Management

* Upload Files
* Download Files
* Soft Delete Files
* Restore Deleted Files
* SHA-256 Duplicate File Detection
* Automatic Restore of Previously Deleted Duplicate Files

### Performance

* Redis Caching
* PostgreSQL Database Storage
* Cache Invalidation on Upload/Delete/Restore
* Optimized File Metadata Retrieval

### User Interface

* React Frontend
* TailwindCSS UI
* Dark / Light Theme
* Dashboard Analytics
* Search Files
* Storage Usage Visualization
* Responsive Design

---

## Tech Stack

### Backend

* Go
* Gin Framework
* PostgreSQL
* Redis
* JWT
* bcrypt

### Frontend

* React
* Axios
* TailwindCSS
* Lucide React Icons
* React Hot Toast

---

## Project Structure

```text
Distribyte/
│
├── backend/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   ├── services/
│   ├── utils/
│   ├── storage/
│   ├── .env
│   └── main.go
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   └── App.jsx
│   │
│   └── package.json
│
└── README.md
```

---

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    password_hash TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Files Table

```sql
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    original_name TEXT NOT NULL,
    stored_name TEXT NOT NULL,
    filepath TEXT NOT NULL,
    size BIGINT NOT NULL,
    file_hash TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    uploaded_at TIMESTAMP DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL
);
```

---

## Environment Variables

Create a `.env` file inside the backend directory.

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=distribyte

JWT_SECRET=your_secret_key

MAX_FILE_SIZE=10485760
```

---

## Running PostgreSQL

Ensure PostgreSQL is running.

Create database:

```sql
CREATE DATABASE distribyte;
```

---

## Running Redis

Install Redis and start server:

```bash
redis-server
```

Verify:

```bash
redis-cli ping
```

Expected:

```text
PONG
```

---

## Backend Setup

Navigate to backend:

```bash
cd backend
```

Install dependencies:

```bash
go mod tidy
```

Run server:

```bash
go run main.go
```

Backend runs on:

```text
http://localhost:8080
```

---

## Frontend Setup

Navigate to frontend:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Start development server:

```bash
npm run dev
```

Frontend runs on:

```text
http://localhost:5173
```

---

## API Endpoints

### Authentication

#### Register

```http
POST /register
```

#### Login

```http
POST /login
```

---

### Files

#### Upload File

```http
POST /upload
```

#### Get Files

```http
GET /files
```

#### Download File

```http
GET /download/:id
```

#### Delete File

```http
DELETE /files/:id
```

#### Restore File

```http
POST /restore/:id
```

#### Get Deleted Files

```http
GET /deleted-files
```

---

## Redis Cache

Distribyte uses Redis for:

* File List Caching
* Deleted File List Caching
* Faster Dashboard Loading
* Reduced Database Queries

Cache is automatically cleared when:

* File Uploaded
* File Deleted
* File Restored

---

## Duplicate Detection

Each uploaded file generates a SHA-256 hash.

Workflow:

```text
Upload File
      ↓
Generate SHA-256
      ↓
Check Existing Hash
      ↓
Duplicate?
 ┌───────────────┐
 │ Yes           │
 └──────┬────────┘
        ↓
Return Duplicate Error

OR

Restore Previously Deleted File
```

---

## Future Roadmap

### Phase 7

* Distributed Storage Nodes
* Node Registration
* File Replication
* Health Checks
* Storage Balancing

### Phase 8

* Docker Deployment
* Kubernetes Support
* Object Versioning
* File Sharing
* Role Based Access Control

### Phase 9

* Monitoring Dashboard
* Metrics & Analytics
* Audit Logs
* Prometheus Integration

---

## Author

Shubham Sharma

Built as a distributed systems and cloud storage project using Go, PostgreSQL, Redis, React, and TailwindCSS.

Currently I am working on Phase 7 of the project.

---

## License

MIT License
