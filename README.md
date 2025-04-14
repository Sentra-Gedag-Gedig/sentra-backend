# SentraPay 💰💳

SentraPay is a comprehensive financial management platform built with Go that provides user authentication, budget management, identity verification, and digital wallet capabilities.

## Features ✨

### User Authentication 🔐
- Registration with email/phone verification
- Multi-factor authentication
- Social login (Google OAuth)
- Biometric authentication (Touch ID)
- User profile management

### Budget Management 📊
- Income and expense tracking
- Transaction categorization
- Period-based financial reports
- Audio notes for transactions
- Customizable categories

### Identity Verification 🔐
- KTP (ID card) detection and data extraction
- Face recognition for authentication
- QRIS code scanning
- Currency recognition

### Digital Wallet 👛
- Balance management
- Virtual account creation via DOKU
- Payment processing with callbacks
- Transaction history
- Secure fund transfers

## Architecture 🏗️

SentraPay follows a clean, modular architecture:

- **API Layer**: RESTful API built with Fiber framework
- **Service Layer**: Core business logic implementation
- **Repository Layer**: Data access patterns for persistence
- **Infrastructure Layer**: External service integrations

## Tech Stack 🛠️

- **Backend**: Go (Golang)
- **Web Framework**: Fiber
- **Database**: PostgreSQL
- **Caching**: Redis
- **Storage**: AWS S3
- **Messaging**: WhatsApp API
- **AI Services**: Google Gemini for image analysis
- **Payment Gateway**: DOKU API

## Prerequisites ✅

- Go 1.18+
- PostgreSQL 13+
- Redis 6+
- AWS S3 credentials
- DOKU payment gateway account
- Google Cloud project for Gemini AI
- WhatsApp integration for notifications

## Environment Variables 🔧

Create a `.env` file in the root directory with the following variables:

```bash
# Application
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=sentrapay
DB_SSLMODE=disable

# Redis
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_ACCESS_TOKEN_SECRET=your-jwt-secret

# AWS S3
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=ap-southeast-1
AWS_BUCKET_NAME=sentrapay-files

# DOKU Payment Gateway
DOKU_CLIENT_ID=your-client-id
DOKU_SECRET_KEY=your-secret-key
DOKU_PUBLIC_KEY=your-public-key
DOKU_IS_PRODUCTION=false

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_STATE=random-state-string

# Gemini AI
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL_NAME=gemini-pro-vision

# AI Services
AI_FACE_DETECTION_URL=ws://face-detection-service:8000/api/v1/face/ws
AI_KTP_DETECTION_URL=ws://ktp-service:8000/api/v1/ktp/ws
AI_QRIS_DETECTION_URL=ws://qris-service:8001/api/v1/qris/ws

# SMTP for email
SMTP_MAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## Installation 📥

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/sentrapay.git
   cd sentrapay
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the database:
   ```bash
   # Run PostgreSQL migrations
   migrate -database "postgres://postgres:yourpassword@localhost:5432/sentrapay?sslmode=disable" -path database/migrations up
   ```

4. Build the application:
   ```bash
   go build -o sentrapay ./cmd/app
   ```

5. Run the application:
   ```bash
   ./sentrapay
   ```

## Docker Deployment 🐳

You can also use Docker Compose to run the entire application stack:

```bash
docker-compose up -d
```

This will start:
- The main Go application
- PostgreSQL database
- Redis cache
- Face detection service
- KTP detection service
- QRIS detection service
- NGINX as a reverse proxy

## API Documentation 📘

### Postman Collection

Access our complete API documentation and test endpoints using our Postman collection:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://braciate-backend.postman.co/workspace/My-Workspace~3c0895d0-8f47-45ff-8232-9471b36c8289/collection/32354585-ae5b5ec5-ccbf-46a0-b4a5-1375abc5d2e4?action=share&creator=32354585&active-environment=32354585-f992d894-dc2a-4b75-8494-aefe3fa343d9)

## Project Structure 📂

```
ProjectGolang/
├── cmd/app/                  # Application entry point
├── database/                 # Database migrations and config
│   ├── migrations/           # SQL migration files
│   └── postgres/             # PostgreSQL connection
├── internal/                 # Internal application code
│   ├── api/                  # API handlers and services
│   │   ├── auth/             # Authentication module
│   │   ├── budget_manager/   # Budget management module
│   │   ├── detection/        # Detection services
│   │   └── sentra_pay/       # Wallet and payments
│   ├── config/               # Application configuration
│   ├── entity/               # Domain entities
│   └── middleware/           # HTTP middleware
├── nginx/                    # NGINX configuration
├── pkg/                      # Shared packages
│   ├── bcrypt/               # Password hashing
│   ├── context/              # Context utilities
│   ├── doku/                 # DOKU payment gateway
│   ├── gemini/               # Google Gemini AI
│   ├── google/               # Google OAuth
│   ├── handlerUtil/          # Handler utilities
│   ├── jwt/                  # JWT authentication
│   ├── log/                  # Logging
│   ├── redis/                # Redis client
│   ├── response/             # HTTP response utilities
│   ├── s3/                   # AWS S3 client
│   ├── smtp/                 # Email sending
│   ├── utils/                # General utilities
│   ├── websocket/            # WebSocket utilities
│   └── whatsapp/             # WhatsApp messaging
└── .env                      # Environment variables
```

## Contributing 🤝

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -am 'Add my feature'`
4. Push to the branch: `git push origin feature/my-feature`
5. Submit a pull request

## License 📝

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments 🙏

- The Go Fiber team for their excellent web framework
- All contributors to the open-source libraries used in this project