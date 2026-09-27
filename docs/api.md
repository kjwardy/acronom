# API Endpoints

## Create Acronym

Create a new acronym entry.

**Endpoint:** `POST /api/acronyms`

**Request Body:**
```json
{
  "acronym": "API",
  "definition": "Application Programming Interface",
  "link": "https://en.wikipedia.org/wiki/API"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "acronym": "API",
  "definition": "Application Programming Interface",
  "link": "https://en.wikipedia.org/wiki/API",
  "created_at": "2026-09-27T14:32:13.578996-04:00",
  "updated_at": "2026-09-27T14:32:13.578996-04:00"
}
```

**Error Responses:**
- `400 Bad Request`: Empty acronym or definition
- `409 Conflict`: Acronym already exists
- `500 Internal Server Error`: Database error

## Update Acronym

Update an existing acronym by ID.

**Endpoint:** `PUT /api/acronyms/:id`

**Request Body:**
```json
{
  "definition": "Application Programming Interface - A set of rules for building software",
  "link": "https://developer.mozilla.org/en-US/docs/Glossary/API"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "acronym": "API",
  "definition": "Application Programming Interface - A set of rules for building software",
  "link": "https://developer.mozilla.org/en-US/docs/Glossary/API",
  "created_at": "2026-09-27T14:32:13.578996-04:00",
  "updated_at": "2026-09-27T14:35:20.123456-04:00"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format or empty definition
- `404 Not Found`: Acronym not found
- `500 Internal Server Error`: Database error

## Delete Acronym

Delete an acronym by ID.

**Endpoint:** `DELETE /api/acronyms/:id`

**Response (204 No Content):** No body

**Error Responses:**
- `400 Bad Request`: Invalid ID format
- `404 Not Found`: Acronym not found
- `500 Internal Server Error`: Database error
