# Grocery List API (MVP)

This document defines the minimal REST API for the shared grocery list app.

---

## Authentication
- **MVP Phase:** No user accounts. Lists are identified by a random `list_id` (invite link).
- **Future Phase:** Add simple user auth (email/password or Apple/Google login).

---

## Endpoints

### 1. Create a new list
**POST /lists**

Creates a new grocery list.

**Request**
```json
{
  "name": "Saturday Shopping"
}
```

**Response**
```json
{
  "list_id": "abcd1234",
  "name": "Saturday Shopping",
  "created_at": "2025-08-26T18:00:00Z"
}
```

### 2. Get a list (and its items)
**GET /lists/{list_id}**

Returns list metadata + items.

```json
{
  "list_id": "abcd1234",
  "name": "Saturday Shopping",
  "items": [
    {
      "item_id": "1",
      "name": "Tomatoes",
      "checked": false
    },
    {
      "item_id": "2",
      "name": "Milk",
      "checked": true
    }
  ]
}
```

### 3. Add an item to a list
**POST /lists/{list_id}/items**
```json
{
  "name": "Bread"
}
```

**Response**
```json
{
  "item_id": "3",
  "name": "Bread",
  "checked": false
}
```

### 4. Toggle item checked/unchecked
**PATCH /lists/{list_id}/items/{item_id}**

```json
{
  "checked": true
}
```

**Response**
```json
{
  "item_id": "3",
  "name": "Bread",
  "checked": true
}
```
### 5. Delete an item
**DELETE /lists/{list_id}/items/{item_id}**

**Response**
```json
{
  "deleted": true
}
```
### 6. Invite link (optional MVP)
**GET /lists/{list_id}/invite**

Returns a sharable invite URL.

**Response**
```json
{
  "invite_url": "https://yourapp.com/join/abcd1234"
}
```

### Obs
Real-time Sync
MVP can use polling (GET /lists/{list_id} every few seconds).

Later: add WebSocket /lists/{list_id}/ws → pushes item changes instantly.

### Summary of Routes
- POST	/lists	Create a new list
- GET	/lists/{list_id}	Fetch list + items
- POST	/lists/{list_id}/items	Add item to list
- PATCH	/lists/{list_id}/items/{item_id}	Toggle check/uncheck
- DELETE	/lists/{list_id}/items/{item_id}	Remove item from list
- GET	/lists/{list_id}/invite	Get invite link (optional)
- WS	/lists/{list_id}/ws (future)	Real-time updates
