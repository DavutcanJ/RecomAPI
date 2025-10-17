# Customer Segmentation API

A comprehensive customer segmentation system built with Go (Gin), Python (FastAPI), MongoDB, and Redis. This API uses MiniBatchKMeans clustering and Ollama LLM to provide intelligent customer segmentation and personalized recommendations.

## 🏗️ **Architecture Overview**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go API        │    │  Python AI      │    │   Databases     │
│  (Port 8008)    │◄──►│  (Port 5005)    │    │                 │
│                 │    │                 │    │  MongoDB        │
│ • Authentication│    │ • ML Clustering │    │  Redis Cache    │
│ • JWT Tokens    │    │ • LLM Insights  │    │  Ollama LLM     │
│ • API Gateway   │    │ • Segmentation  │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 **Features**

### **Authentication & Security**
- ✅ JWT-based authentication
- ✅ User registration and login
- ✅ Password hashing with bcrypt
- ✅ Bearer token authorization
- ✅ MongoDB user storage with indexes

### **Customer Segmentation**
- ✅ MiniBatchKMeans clustering algorithm
- ✅ Ollama LLM integration for insights
- ✅ Multi-business type support (e-commerce, finance, retail)
- ✅ Confidence scoring
- ✅ Personalized recommendations

### **Performance & Scalability**
- ✅ Redis caching layer
- ✅ MongoDB connection pooling
- ✅ Async Python processing
- ✅ RESTful API design
- ✅ Comprehensive error handling

### **Developer Experience**
- ✅ Swagger/OpenAPI documentation
- ✅ Interactive API testing
- ✅ Comprehensive logging
- ✅ Environment-based configuration

## 📁 **Project Structure**

```
RecomAPI/
├── customer-segmentation-api/ (Go Backend)
│   ├── cmd/
│   │   └── main.go                    # Application entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── database.go            # MongoDB configuration
│   │   ├── handlers/
│   │   │   ├── auth.go               # Authentication endpoints
│   │   │   ├── user.go               # User management
│   │   │   └── segment.go            # Segmentation API
│   │   ├── middleware/
│   │   │   └── auth.go               # JWT middleware
│   │   ├── models/
│   │   │   ├── models.go             # API models
│   │   │   └── user.go               # User models
│   │   ├── repository/
│   │   │   └── user_repository.go    # Database operations
│   │   ├── service/
│   │   │   └── user_service.go       # Business logic
│   │   └── services/
│   │       └── cache.go              # Redis caching
│   ├── docs/                         # Swagger documentation
│   ├── go.mod
│   └── .env
│
└── rec-ai-service/ (Python AI Service)
    ├── main.py                       # FastAPI application
    ├── models/
    │   └── models.py                 # Pydantic models
    ├── service/
    │   └── segment.py                # ML clustering logic
    ├── requirements.txt
    └── .env
```

## 🛠️ **Technology Stack**

### **Backend (Go)**
- **Framework**: Gin Web Framework
- **Authentication**: JWT with golang-jwt/jwt/v5
- **Database**: MongoDB with official Go driver
- **Caching**: Redis with go-redis/v9
- **Documentation**: Swagger/OpenAPI with swaggo
- **Security**: bcrypt password hashing

### **AI Service (Python)**
- **Framework**: FastAPI
- **ML**: scikit-learn (MiniBatchKMeans)
- **LLM**: Ollama integration
- **Data Processing**: NumPy, Pandas
- **Validation**: Pydantic v2

### **Databases**
- **Primary**: MongoDB (User data, configurations)
- **Cache**: Redis (Segmentation results)
- **AI**: Ollama (LLM processing)

## ⚙️ **Installation & Setup**

### **Prerequisites**
```bash
# Required software
- Go 1.21+
- Python 3.9+
- MongoDB 6.0+
- Redis 7.0+
- Ollama
```

### **1. Clone the Repository**
```bash
git clone <repository-url>
cd RecomAPI
```

### **2. Setup Go Backend**
```bash
cd customer-segmentation-api

# Install dependencies
go mod tidy
go get go.mongodb.org/mongo-driver/mongo
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/redis/go-redis/v9

# Install Swagger tool
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swag init -g cmd/main.go -o docs
```

### **3. Setup Python AI Service**
```bash
cd ../rec-ai-service

# Create virtual environment
python -m venv venv
source venv/bin/activate  # Linux/Mac
# venv\Scripts\activate   # Windows

# Install dependencies
pip install -r requirements.txt
```

### **4. Setup Databases**

**MongoDB:**
```bash
# Install MongoDB (Ubuntu/Debian)
sudo apt update
sudo apt install mongodb

# Start MongoDB
sudo systemctl start mongodb
sudo systemctl enable mongodb
```

**Redis:**
```bash
# Install Redis (Ubuntu/Debian)
sudo apt update
sudo apt install redis-server

# Start Redis
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

**Ollama:**
```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model (e.g., llama3.1)
ollama pull llama3.1
```

### **5. Environment Configuration**

**Go Backend (.env):**
```bash
# customer-segmentation-api/.env
JWT_SECRET=your-super-secret-jwt-key-here-make-it-long-and-complex
MONGODB_URI=mongodb://localhost:27017
DB_NAME=recom_api
PORT=8008
PYTHON_SERVICE_URL=http://localhost:5005
REDIS_URI=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Python AI Service (.env):**
```bash
# rec-ai-service/.env
OLLAMA_MODEL=llama3.1
OLLAMA_HOST=http://localhost:11434
REDIS_URI=localhost:6379
LOG_LEVEL=INFO
```

## 🚀 **Running the Services**

### **1. Start AI Service**
```bash
cd rec-ai-service
source venv/bin/activate
python main.py
# Service starts on http://localhost:5005
```

### **2. Start Go Backend**
```bash
cd customer-segmentation-api
go run cmd/main.go
# Service starts on http://localhost:8008
```

### **3. Access Swagger Documentation**
Open your browser: `http://localhost:8008/swagger/index.html`

## 📊 **API Usage Examples**

### **1. User Registration**
```bash
curl -X POST http://localhost:8008/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword123",
    "full_name": "John Doe"
  }'
```

### **2. User Login**
```bash
curl -X POST http://localhost:8008/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "securepassword123"
  }'
```

### **3. Customer Segmentation**
```bash
curl -X POST http://localhost:8008/api/segment \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "customer": {
      "customer_id": "12345",
      "purchase_history": ["electronics", "books", "clothing"],
      "demographics": {
        "age": 30,
        "location": "Istanbul",
        "gender": "male"
      },
      "behavioral_data": {
        "website_visits": 25,
        "last_purchase_days": 15,
        "avg_order_value": 250.50
      }
    },
    "business_type": "e-commerce"
  }'
```

### **4. Expected Response**
```json
{
  "customer_id": "12345",
  "segment": "tech-savvy",
  "confidence": 0.85,
  "recommended_offer": "15% discount on electronics",
  "insights": [
    {
      "action": "Email campaign",
      "priority": "high"
    },
    {
      "action": "Personalized recommendations", 
      "priority": "medium"
    }
  ]
}
```

## 🧪 **Testing**

### **Health Checks**
```bash
# Go Backend Health
curl http://localhost:8008/health

# Python AI Service Health  
curl http://localhost:5005/health
```

### **Database Connections**
```bash
# Test MongoDB
mongosh recom_api

# Test Redis
redis-cli ping
```

### **LLM Service**
```bash
# Test Ollama
ollama list
```

## 🔧 **Development**

### **Hot Reload (Go)**
```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### **Python Development**
```bash
# Run with auto-reload
uvicorn main:app --reload --host 0.0.0.0 --port 5005
```

### **Regenerate Swagger Docs**
```bash
cd customer-segmentation-api
swag init -g cmd/main.go -o docs
```

## 📝 **API Documentation**

- **Swagger UI**: `http://localhost:8008/swagger/index.html`
- **OpenAPI JSON**: `http://localhost:8008/swagger/doc.json`

### **Available Endpoints**

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/register` | User registration | ❌ |
| POST | `/api/login` | User authentication | ❌ |
| GET | `/api/profile` | Get user profile | ✅ |
| POST | `/api/segment` | Customer segmentation | ✅ |
| GET | `/health` | Health check | ❌ |

## 🔒 **Security Features**

- **JWT Authentication** with configurable expiration
- **Password Hashing** using bcrypt with salt
- **Environment Variables** for sensitive data
- **Database Indexes** for performance and uniqueness
- **CORS Support** (configurable)
- **Rate Limiting** (can be implemented)

## 🚀 **Deployment**

### **Docker Support** (Future Enhancement)
```dockerfile
# Example Go Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### **Production Considerations**
- Use Docker containers for deployment
- Set up reverse proxy (nginx) for load balancing
- Configure database replicas for high availability
- Implement monitoring and logging (Prometheus, Grafana)
- Use environment-specific configurations
- Set up CI/CD pipelines

## 🤝 **Contributing**

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 **License**

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 **Support & Troubleshooting**

### **Common Issues**

**1. "Python service error"**
```bash
# Check if Python service is running
curl http://localhost:5005/health

# Check Python logs
cd rec-ai-service && python main.py
```

**2. "Database connection failed"**
```bash
# Check MongoDB status
sudo systemctl status mongodb

# Check Redis status  
redis-cli ping
```

**3. "Invalid JWT token"**
- Ensure Bearer token format: `Bearer <token>`
- Check JWT_SECRET environment variable
- Verify token hasn't expired (24h default)

**4. "Swagger not loading"**
```bash
# Regenerate docs
swag init -g cmd/main.go -o docs

# Restart Go service
go run cmd/main.go
```

### **Performance Tuning**
- Adjust MongoDB connection pool size
- Configure Redis memory settings
- Optimize ML model parameters
- Implement database indexing
- Use caching for frequent queries

---

**Built with ❤️ using Go, Python, MongoDB, Redis, and AI**
