# 📬 Daily Fortune API Manual

This manual explains how to use the API endpoints provided by the Daily Fortune service.

---

## 🃏 Draw a Card

### `GET /draw?username=<your-username>`

- Draws a fortune card for the given username.
- A user can only draw once per day.
- The result is cached using Redis.

**Example:**

```http
GET /draw?username=testuser
```

**Success Response:**
```json
{
  "card": {
    "name": "Lucky Star",
    "message": "You will shine today!",
    "type": "fortune",
    "rarity": "rare",
    "imagePath": "/img/luckystar.png",
    "status": "Y"
  },
  "message": "You already drew a card today"
}
```

---

## 🧾 Card Management (Admin Only)

All endpoints require a valid `username` of an admin user.

### 🔍 Get All Cards

**GET** `/cards`

**Query Params:**
- `status` (optional): Y / N
- `name` (optional): partial match
- `username` (required): must be an admin

**Example:**
```http
GET /cards?status=Y&name=luck&username=admin
```

---

### ➕ Create a Card

**POST** `/card`

```json
{
  "username": "admin",
  "name": "Lucky Star",
  "message": "You will shine today!",
  "type": "fortune",
  "rarity": "rare",
  "imagePath": "/img/luckystar.png",
  "status": "Y"
}
```

---

### ✏️ Update a Card

**PATCH** `/card`

```json
{
  "username": "admin",
  "name": "Lucky Star",
  "message": "Updated message",
  "type": "fortune",
  "rarity": "legendary",
  "imagePath": "/img/luckystar.png",
  "status": "Y"
}
```

---

### ❌ Delete a Card

**DELETE** `/card`

```json
{
  "username": "admin",
  "name": "Lucky Star"
}
```

Soft-deletes the card by setting `"status": "N"`.

---

## 🔐 Authentication

No token-based auth yet. Admin permissions are verified using a hardcoded admin username (via `utils.IsAdmin()` logic).

---

## 💬 Notes

- Redis must be running for `/draw` caching to work properly.
- All input must be valid JSON where applicable.
- Consider adding Swagger or Postman collection for easier testing.

---

Happy drawing! 🍀