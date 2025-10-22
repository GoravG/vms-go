# Gateway Service

A Go-based gateway service that provides authentication, user management, and WebSocket capabilities. This service is part of the VMS (Visitor Management System) project.

## Features

- User Authentication (Login/Registration)
- WebSocket Support for Real-time Communications
- Health Check Endpoint
- MySQL Database Integration
- Redis Integration for Real-time Updates
- HMAC-based Token Authentication
- Structured Logging

## Prerequisites

- Go 1.20 or higher
- Docker and Docker Compose
- MySQL 8.0
- Redis

## Project Structure

```
├── cmd/
│   └── server/           # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── db/             # Database connection and operations
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── migrations/     # Database migrations
│   ├── models/         # Data models
│   ├── redis/          # Redis integration
│   ├── security/       # Security utilities
│   ├── token/          # Token management
│   ├── utils/          # Utility functions
│   └── ws/             # WebSocket implementation
└── bruno-api/          # API Collection for testing
```

## Getting Started

### Running with Docker Compose

1. Clone the repository
2. Navigate to the project directory
3. Start the services:
   ```bash
   docker-compose up -d
   ```

This will start:
- MySQL database on port 3306
- The backend service with all necessary configurations

### Environment Variables

The service uses the following environment variables:
- `DB_USER` - Database username
- `DB_PASS` - Database password
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `HMAC_SECRET_KEY` - Secret key for HMAC token generation
- `REDIS_HOST` - Redis host
- `REDIS_PORT` - Redis port
- `LOG_FILE` - Log file location
- `LOG_FORMAT` - Logging format ('ecs' or 'plain')

## API Endpoints

- `GET /health` - Health check endpoint
- `POST /login` - User login
- `POST /users` - Create new user
- `POST /checkin` - User check-in
- `GET /ws` - WebSocket endpoint

## Testing

The project includes a Bruno API collection for testing the endpoints. You can find the collection in the `bruno-api` directory.

## Logging

The service implements structured logging with the following features:
- Log rotation
- Error tracking
- Application events logging

## Development

To run the service locally for development:

1. Set up your environment variables
2. Start MySQL and Redis (using Docker or locally)
3. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

## Docker Support

The project includes:
- `Dockerfile` for building the service

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
