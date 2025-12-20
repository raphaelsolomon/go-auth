# Testing Guide for Go Auth API

## Setup

1. **Create environment file:**
   ```bash
   cp .env.example .env
   ```

   Update `.env` with your configuration:
   - `PORT`: Server port (default: 8080)
   - `DB_NAME`: SQLite database file name (default: golang.db)
   - `SECRET_KEY`: JWT secret key (use a strong random string in production)

2. **Build the application:**
   ```bash
   go build -o bin/auth-server ./cmd/main.go
   ```

3. **Run the server:**
   ```bash
   ./bin/auth-server
   ```

## API Endpoints

### Public Endpoints

#### 1. Health Check
```bash
curl http://localhost:8080/health
```
Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-11-29T14:45:58Z"
}
```

#### 2. Welcome
```bash
curl http://localhost:8080/
```
Response:
```json
{
  "message": "Welcome to my public API"
}
```

#### 3. User Registration
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123",
    "firstname": "John",
    "lastname": "Doe",
    "phone_no": "+1234567890"
  }'
```

**Validation Rules:**
- Email: Valid email format
- Password: Minimum 8 characters
- First/Last Name: Minimum 2 characters
- Phone Number: Valid international format (E.164)

Response:
```json
{
  "message": "User created successfully"
}
```

#### 4. User Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": "2025-11-30T14:46:29Z",
  "token_type": "Bearer",
  "user": {
    "ID": 1,
    "Uuid": "ab8b1a24-9c4b-433d-98e0-dfbfeeb0a00e",
    "Email": "test@example.com",
    "Firstname": "John",
    "Lastname": "Doe",
    "Status": true,
    "PhoneNo": "+1234567890",
    "IsEmailverified": false
  }
}
```

### Protected Endpoints (Require JWT Token)

#### 5. Get User Profile
```bash
curl -X GET http://localhost:8080/api/v1/user \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "ID": 1,
  "Uuid": "ab8b1a24-9c4b-433d-98e0-dfbfeeb0a00e",
  "Email": "test@example.com",
  "Firstname": "John",
  "Lastname": "Doe",
  "Status": true,
  "PhoneNo": "+1234567890",
  "IsEmailverified": false
}
```

#### 6. Update User Profile
```bash
curl -X PUT http://localhost:8080/api/v1/user \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "Jane",
    "lastname": "Smith",
    "phone_no": "+9876543210"
  }'
```

Response:
```json
{
  "message": "User updated successfully",
  "user": {
    "ID": 1,
    "Uuid": "ab8b1a24-9c4b-433d-98e0-dfbfeeb0a00e",
    "Email": "test@example.com",
    "Firstname": "Jane",
    "Lastname": "Smith",
    "Status": true,
    "PhoneNo": "+9876543210",
    "IsEmailverified": false
  }
}
```

## Testing Validation

### Invalid Email
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "testpass123",
    "firstname": "John",
    "lastname": "Doe",
    "phone_no": "+1234567890"
  }'
```
Response: `{"error": "invalid email format"}`

### Short Password
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "short",
    "firstname": "John",
    "lastname": "Doe",
    "phone_no": "+1234567890"
  }'
```
Response: `{"error": "password must be at least 8 characters long"}`

### Invalid Phone Number
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123",
    "firstname": "John",
    "lastname": "Doe",
    "phone_no": "invalid"
  }'
```
Response: `{"error": "invalid phone number format"}`

## Features Added

### ✅ Fixed Issues
- Fixed typo: `prinmaryKey` → `primaryKey` in User model
- Fixed typo: `SetLasttname` → `SetLastname` in utils

### ✅ Input Validation
- Email format validation
- Password strength validation (min 8 characters)
- Name validation (min 2 characters)
- Phone number format validation (E.164 international format)

### ✅ CORS Configuration
- Enabled CORS for frontend integration
- Allows all origins (configure as needed for production)
- Supports all standard HTTP methods
- 12-hour preflight cache

### ✅ Improved Error Handling
- JWT generation errors are now properly handled
- Database operation errors are caught and reported
- Consistent error response format

### ✅ Database Connection Pooling
- Single database instance initialized at startup
- Connection pooling configured (10 idle, 100 max open)
- Auto-migration runs once at startup
- Improved performance and resource management

### ✅ Health Check Endpoint
- Added `/health` endpoint for monitoring
- Returns server status and timestamp

### ✅ Environment Configuration
- Created `.env.example` with all required variables
- Easy setup for new developers

## Testing with Frontend

The API now includes CORS support, making it easy to test with frontend applications:

```javascript
// Example frontend fetch
const response = await fetch('http://localhost:8080/api/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: 'test@example.com',
    password: 'testpass123'
  })
});

const data = await response.json();
console.log(data.access_token);
```

## Production Considerations

Before deploying to production:

1. **Change SECRET_KEY** in `.env` to a strong random string
2. **Update CORS origins** to allow only trusted domains
3. **Use PostgreSQL/MySQL** instead of SQLite for better concurrency
4. **Enable HTTPS** for secure communication
5. **Add rate limiting** to prevent abuse
6. **Set up proper logging** (file-based or centralized)
7. **Enable email verification** (implement the existing IsEmailverified field)
8. **Add token refresh** mechanism for better UX
9. **Set GIN_MODE=release** for production builds
