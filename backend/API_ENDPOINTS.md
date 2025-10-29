# IoT Dashboard API Endpoints

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

### Web Application (JWT)
- Include `Authorization: Bearer <token>` header
- Login endpoint: `POST /api/v1/public/auth/login`

### IoT Devices (API Key)
- Include `Authorization: <api_key>` header
- Or `Authorization: Bearer <device_token>` header

## Public Endpoints (No Authentication)

### Authentication
- `POST /api/v1/public/auth/login` - User login
- `POST /api/v1/public/auth/register` - User registration

## Web Application Endpoints (JWT Required)

### Authentication
- `POST /api/v1/web/auth/logout` - User logout
- `GET /api/v1/web/auth/profile` - Get user profile
- `PUT /api/v1/web/auth/profile` - Update user profile

### Device Management
- `GET /api/v1/web/devices` - Get all devices
- `GET /api/v1/web/devices/:id` - Get device by ID
- `POST /api/v1/web/devices` - Create new device
- `PUT /api/v1/web/devices/:id` - Update device
- `DELETE /api/v1/web/devices/:id` - Delete device
- `GET /api/v1/web/devices/:id/data` - Get device sensor data

### Analytics
- `GET /api/v1/analytics/devices/:id/summary` - Get device analytics summary

## IoT Device Endpoints (API Key Required)

### Sensor Data
- `POST /api/v1/iot/data` - Send sensor data from device
- `GET /api/v1/iot/devices/:id/status` - Get device status
- `PUT /api/v1/iot/devices/:id/status` - Update device status

### Device Commands
- `GET /api/v1/iot/devices/:id/commands` - Get pending commands for device
- `POST /api/v1/iot/devices/:id/commands/:commandId/response` - Send command response

## Health Check
- `GET /health` - Server health check

## Request/Response Format

### Standard Response Format
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "error": "Error message if any"
}
```

### Error Response Format
```json
{
  "success": false,
  "error": "Error description"
}
```

## Example Requests

### Login
```bash
curl -X POST http://localhost:8080/api/v1/public/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "password"}'
```

### Send Sensor Data (IoT Device)
```bash
curl -X POST http://localhost:8080/api/v1/iot/data \
  -H "Content-Type: application/json" \
  -H "Authorization: your-api-key" \
  -d '{
    "device_id": "device123",
    "temperature": 22.5,
    "humidity": 45.2,
    "co2": 400,
    "pressure": 1013.25,
    "light": 850
  }'
```

### Get Devices (Web)
```bash
curl -X GET http://localhost:8080/api/v1/web/devices \
  -H "Authorization: Bearer your-jwt-token"
```
