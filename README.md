# 🔥 NoBurn - AI-Powered HR Analytics Platform

> Reduce employee churn with AI-driven sentiment analysis and predictive analytics for Indian startups

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-19.2+-61DAFB?style=flat&logo=react)](https://reactjs.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Features

- 🎯 **Churn Prediction** - ML-powered attrition risk scoring (0-100%)
- 💬 **Sentiment Analysis** - IndicBERT for Indian language support (Hindi, Tamil, etc.)
- 📊 **Real-time Dashboard** - Employee analytics and risk insights
- 📧 **Smart Notifications** - Automated email alerts for high-risk employees
- 🌐 **Multi-language** - Regional Indian language support
- ⚡ **Background Jobs** - Async processing with Asynq
- 🏢 **Multi-tenancy** - Company-specific SMTP and settings
- 📱 **Modern UI** - React + TailwindCSS responsive interface

## 🏗️ Architecture

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   Frontend  │─────▶│   API (Go)   │─────▶│ PostgreSQL  │
│ React+Nginx │      │  Chi Router  │      │   Database  │
└─────────────┘      └──────┬───────┘      └─────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │    Redis     │
                     │ Queue/Cache  │
                     └──────┬───────┘
                            │
                            ▼
                     ┌──────────────┐
                     │    Worker    │
                     │ Asynq Jobs   │
                     └──────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)
- Node.js 20+ (for frontend development)

### 🐳 Docker Setup (Recommended)

```bash
# Clone repository
git clone https://github.com/yourusername/NoBurn.git
cd NoBurn

# Start all services
docker compose up -d

# Check status
docker ps
```

**Services:**
- Frontend: http://localhost:3002 (Vite + React + TypeScript + Shadcn UI)
- API: http://localhost:3000
- PostgreSQL: localhost:5433
- Redis: localhost:6379

### 🛠️ Local Development

```bash
# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Setup environment
cp .env.example .env
# Edit .env with your credentials

# Run migrations
make migrate-up

# Terminal 1: API Server
make dev-api

# Terminal 2: Worker
make dev-worker

# Terminal 3: Frontend
cd frontend && npm run dev
```

## 📚 API Documentation

### Authentication

```bash
# Register company
POST /auth/register
{
  "email": "admin@company.com",
  "password": "secure123",
  "name": "Admin Name",
  "company_name": "Your Company"
}

# Login
POST /auth/login
{
  "email": "admin@company.com",
  "password": "secure123"
}
```

### Employee Management

```bash
# Add single employee
POST /api/employees
Authorization: Bearer <token>
{
  "email": "employee@company.com",
  "name": "John Doe"
}

# Bulk upload (CSV)
POST /api/employees/bulk
Authorization: Bearer <token>
Content-Type: multipart/form-data
file: employees.csv

# List employees
GET /api/employees
Authorization: Bearer <token>
```

### Survey Management

```bash
# Create survey (auto-sends emails to employees)
POST /api/surveys
Authorization: Bearer <token>
{
  "title": "Q4 Employee Satisfaction",
  "questions": [
    "How satisfied are you with your role?",
    "Rate work-life balance (1-5)",
    "Any suggestions?"
  ]
}

# List surveys
GET /api/surveys
Authorization: Bearer <token>

# Submit response (public endpoint via email link)
POST /api/surveys/responses/public
{
  "survey_id": 1,
  "user_token": 123,
  "responses": ["Very satisfied", "5", "Great team!"]
}
```

### Analytics

```bash
# Dashboard
GET /api/dashboard
Authorization: Bearer <token>

# Attrition risks
GET /api/attrition-risks
Authorization: Bearer <token>
```

### Settings

```bash
# Get SMTP settings
GET /api/settings/smtp
Authorization: Bearer <token>

# Update SMTP (company-specific email)
PUT /api/settings/smtp
Authorization: Bearer <token>
{
  "smtp_host": "smtp.gmail.com",
  "smtp_port": 587,
  "smtp_user": "company@gmail.com",
  "smtp_password": "app_password"
}
```

## 🗂️ Project Structure

```
NoBurn/
├── cmd/
│   ├── api/main.go              # HTTP API server
│   └── worker/main.go           # Background job processor
├── internal/
│   ├── config/                  # Configuration management
│   ├── database/                # GORM database connection
│   ├── handlers/                # HTTP request handlers
│   │   ├── auth/                # Authentication
│   │   ├── employee/            # Employee management
│   │   ├── survey/              # Survey CRUD
│   │   ├── analytics/           # Dashboard & analytics
│   │   └── settings/            # Company settings
│   ├── middleware/              # Auth, logging middleware
│   ├── models/                  # GORM models
│   ├── repository/              # Database queries
│   ├── services/                # Business logic
│   │   ├── sentiment/           # ML sentiment analysis
│   │   └── notification/        # Email/Slack alerts
│   ├── worker/                  # Asynq background jobs
│   └── utils/                   # Helper functions
├── frontend/                    # React application
│   ├── src/
│   │   ├── pages/               # Dashboard, Surveys, Employees, Settings
│   │   ├── components/          # Reusable components
│   │   └── context/             # Auth context
│   ├── Dockerfile               # Frontend container
│   └── nginx.conf               # Nginx config
├── migrations/                  # SQL migrations
├── docker-compose.yml           # Multi-container setup
└── Makefile                     # Development commands
```

## 🔧 Configuration

### Environment Variables

```bash
# Server
PORT=3000
ENV=production

# Database
DB_URL=postgres://user:pass@localhost:5432/noburn_db?sslmode=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# Redis
REDIS_URL=redis://localhost:6379

# Hugging Face (IndicBERT)
HUGGING_FACE_TOKEN=hf_xxxxxxxxxxxxx

# Email (Optional - can be configured per company in UI)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=default@gmail.com
SMTP_PASSWORD=app_password

# Slack (Optional)
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/xxx
```

## 🤖 ML & AI

### IndicBERT Integration

NoBurn uses **IndicBERT** (ai4bharat/indic-bert) for sentiment analysis:

- **12 Indian Languages**: Hindi, Bengali, Gujarati, Kannada, Malayalam, Marathi, Oriya, Punjabi, Tamil, Telugu, Urdu
- **State-of-the-art**: Developed by AI4Bharat (IIT Madras)
- **Fallback**: Rule-based sentiment analysis when ML unavailable

### Churn Prediction Algorithm

```
Risk Score = f(
  - Average Sentiment Score
  - Response Frequency
  - Sentiment Trend
  - Negative Keywords
)

Thresholds:
- High Risk: > 70%
- Medium Risk: 40-70%
- Low Risk: < 40%
```

## 📦 Deployment

### Docker Production

```bash
# Build and deploy
docker compose up -d

# View logs
docker logs -f noburn-api
docker logs -f noburn-worker

# Scale workers
docker compose up -d --scale worker=3
```

### Manual Deployment

```bash
# Build binaries
make build-all

# Run migrations
DB_URL=<prod-url> make migrate-up

# Start services
./bin/api &
./bin/worker &
```

### Environment Setup

1. **Database**: PostgreSQL 15+
2. **Cache**: Redis 7+
3. **SMTP**: Gmail App Password or SendGrid
4. **ML**: Hugging Face API token (free tier available)

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Test specific package
go test ./internal/services/sentiment/...
```

## 📊 Monitoring

### Health Check

```bash
curl http://localhost:3000/health
```

### Metrics

- API response times
- Worker job processing
- Database connections
- Redis queue length

## 🔐 Security

- ✅ JWT-based authentication
- ✅ Password hashing (bcrypt)
- ✅ SQL injection prevention (GORM)
- ✅ CORS configuration
- ✅ Rate limiting
- ✅ Environment variable secrets

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) file

## 🙏 Acknowledgments

- [IndicBERT](https://huggingface.co/ai4bharat/indic-bert) by AI4Bharat
- [Asynq](https://github.com/hibiken/asynq) for background jobs
- [Chi](https://github.com/go-chi/chi) for routing
- [GORM](https://gorm.io) for ORM

## 📧 Support

For issues and questions:
- GitHub Issues: [Create Issue](https://github.com/yourusername/NoBurn/issues)
- Email: support@noburn.ai

---

**Built with ❤️ for Indian startups**
