package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Username  string            `bson:"username" json:"username" binding:"required,min=3,max=50"`
    Email     string            `bson:"email" json:"email" binding:"required,email"`
    Password  string            `bson:"password" json:"-"` // Never return password in JSON
    FullName  string            `bson:"full_name" json:"full_name" binding:"required,min=2,max=100"`
    Role      string            `bson:"role" json:"role"` // admin, user
    IsActive  bool              `bson:"is_active" json:"is_active"`
    CreatedAt time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=100"`
    FullName string `json:"full_name" binding:"required,min=2,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" example:"test-user"`
	Password string `json:"password" example:"test-password"`
}

type LoginResponse struct {
    Token     string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    User      UserProfile `json:"user"`
    ExpiresAt int64       `json:"expires_at"`
}

type SegmentRequest struct {
	Customer struct {
		CustomerID      string                 `json:"customer_id" example:"12345"`
		PurchaseHistory []string               `json:"purchase_history" example:"laptop,mouse"`
		Demographics    map[string]interface{} `json:"demographics" example:"{\"age\": 28, \"location\": \"Istanbul\", \"income_level\": \"medium\"}"`
		LastInteraction string                 `json:"last_interaction" example:"2025-10-04"`
	} `json:"customer"`
	BusinessType string `json:"business_type" example:"ecommerce"`
}

type SegmentResponse struct {
	CustomerID       string    `json:"customer_id" example:"12345"`
	Segment          string    `json:"segment" example:"tech-savvy"`
	RecommendedOffer string    `json:"recommended_offer" example:"15% discount"`
	Confidence       float64   `json:"confidence" example:"0.85"`
	Insights         []Insight `json:"insights"`
}

type Insight struct {
	Priority string `json:"priority" example:"high"`
	Action   string `json:"action" example:"Email campaign"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserProfile struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    FullName string `json:"full_name"`
    Role     string `json:"role"`
    IsActive bool   `json:"is_active"`
}