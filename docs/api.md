# API Endpoints

## Total Acronyms Metric

Return the total number of glossary entries stored. Each meaning is counted as an individual entry, even when several entries share the same acronym.

**Endpoint:** `GET /api/metrics/total-acronyms`

**Response (200 OK):**
```json
{
  "total_acronyms": 42
}
```

An empty database returns `200 OK` with `total_acronyms` set to `0`.

**Error Responses:**
- `500 Internal Server Error`: Database error

## List Categories

Return all reusable categories in alphabetical order.

**Endpoint:** `GET /api/categories`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Platform",
    "colour": "#1976D2",
    "created_at": "2026-09-27T14:32:13.578996-04:00",
    "updated_at": "2026-09-27T14:32:13.578996-04:00"
  }
]
```

An empty category collection returns `200 OK` with an empty array.

## Create Category

Create a reusable category. Names are trimmed, repeated whitespace is collapsed, names are converted to Title Case, and uniqueness is case-insensitive. Names may contain at most 30 characters.

**Endpoint:** `POST /api/categories`

**Request Body:**
```json
{
  "name": "Platform",
  "colour": "#1976D2"
}
```

Approved colour values are `#1976D2`, `#00796B`, `#388E3C`, `#FFA000`, `#F57C00`, `#D32F2F`, `#7B1FA2`, and `#455A64`.

**Response (201 Created):** The persisted category object.

**Error Responses:**
- `400 Bad Request`: Missing or invalid name, or a colour outside the approved palette
- `409 Conflict`: A category with the same case-insensitive name already exists
- `500 Internal Server Error`: Database error

Category update and deletion endpoints are not currently available.

## Search Acronyms

Search acronym names and definitions. Matching is case-insensitive, with exact acronym matches followed by acronym prefixes and then other matches. Responses are limited to 50 entries.

**Endpoint:** `GET /api/search?q=api`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "acronym": "API",
    "definition": "Application Programming Interface",
    "link": "https://en.wikipedia.org/wiki/API",
    "created_at": "2026-09-27T14:32:13.578996-04:00",
    "updated_at": "2026-09-27T14:32:13.578996-04:00",
    "categories": [
      {
        "id": 1,
        "name": "Platform",
        "colour": "#1976D2"
      }
    ]
  }
]
```

When no entries match, the endpoint returns `200 OK` with an empty array.

**Error Responses:**
- `400 Bad Request`: Missing or empty `q` query parameter
- `500 Internal Server Error`: Database error

## Create Acronym

Create a new acronym entry. Multiple entries may use the same acronym when it has more than one meaning.

**Endpoint:** `POST /api/acronyms`

**Request Body:**
```json
{
  "acronym": "API",
  "definition": "Application Programming Interface",
  "link": "https://en.wikipedia.org/wiki/API",
  "category_ids": [1, 3]
}
```

`category_ids` is optional. At most five unique category IDs are accepted; duplicate IDs are normalized and unknown IDs return `400 Bad Request`. Omitting it or sending an empty array creates an uncategorized entry. Acronym responses always include `categories`, using an empty array for uncategorized entries.

**Response (201 Created):**
```json
{
  "id": 1,
  "acronym": "API",
  "definition": "Application Programming Interface",
  "link": "https://en.wikipedia.org/wiki/API",
  "created_at": "2026-09-27T14:32:13.578996-04:00",
  "updated_at": "2026-09-27T14:32:13.578996-04:00",
  "categories": [
    {
      "id": 1,
      "name": "Platform",
      "colour": "#1976D2"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Empty acronym or definition
- `409 Conflict`: The same acronym and definition already exist
- `500 Internal Server Error`: Database error

## Update Acronym

Update an existing acronym by ID.

**Endpoint:** `PUT /api/acronyms/:id`

**Request Body:**
```json
{
  "definition": "Application Programming Interface - A set of rules for building software",
  "link": "https://developer.mozilla.org/en-US/docs/Glossary/API",
  "category_ids": [1, 3]
}
```

Omitting `category_ids` preserves existing assignments. Sending an empty array removes all assignments. A supplied array replaces assignments atomically and accepts at most five unique, existing IDs.

**Response (200 OK):**
```json
{
  "id": 1,
  "acronym": "API",
  "definition": "Application Programming Interface - A set of rules for building software",
  "link": "https://developer.mozilla.org/en-US/docs/Glossary/API",
  "created_at": "2026-09-27T14:32:13.578996-04:00",
  "updated_at": "2026-09-27T14:35:20.123456-04:00",
  "categories": [
    {
      "id": 1,
      "name": "Platform",
      "colour": "#1976D2"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID format or empty definition
- `409 Conflict`: The update would duplicate an existing acronym and definition
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
