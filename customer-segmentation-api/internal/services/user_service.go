package services

import (
    "errors"
    "os"
    "time"

    "github.com/DavutcanJ/customer-segmentation-api/internal/models"
    "github.com/DavutcanJ/customer-segmentation-api/internal/repository"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    userRepo *repository.UserRepository
}

func NewUserService() *UserService {
    return &UserService{
        userRepo: repository.NewUserRepository(),
    }
}

func (s *UserService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
    // Check if username already exists
    existingUser, _ := s.userRepo.GetUserByUsername(req.Username)
    if existingUser != nil {
        return nil, errors.New("username already exists")
    }

    // Check if email already exists
    existingUser, _ = s.userRepo.GetUserByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("email already exists")
    }

    // Hash password with high cost for security
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }

    user := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
        FullName: req.FullName,
    }

    err = s.userRepo.CreateUser(user)
    if err != nil {
        return nil, errors.New("failed to create user")
    }

    return user, nil
}

func (s *UserService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
    // Get user by username
    user, err := s.userRepo.GetUserByUsername(req.Username)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if !user.IsActive {
        return nil, errors.New("user account is deactivated")
    }

    // Verify password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Generate JWT token
    token, expiresAt, err := s.generateJWTToken(user)
    if err != nil {
        return nil, errors.New("failed to generate token")
    }

    userProfile := models.UserProfile{
        ID:       user.ID.Hex(),
        Username: user.Username,
        Email:    user.Email,
        FullName: user.FullName,
        Role:     user.Role,
        IsActive: user.IsActive,
    }

    return &models.LoginResponse{
        Token:     token,
        User:      userProfile,
        ExpiresAt: expiresAt,
    }, nil
}

func (s *UserService) generateJWTToken(user *models.User) (string, int64, error) {
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "fallback-secret-key"
    }

    expiresAt := time.Now().Add(time.Hour * 24).Unix() // 24 hours

    claims := jwt.MapClaims{
        "user_id":  user.ID.Hex(),
        "username": user.Username,
        "email":    user.Email,
        "role":     user.Role,
        "exp":      expiresAt,
        "iat":      time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        return "", 0, err
    }

    return tokenString, expiresAt, nil
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
    return s.userRepo.GetUserByID(userID)
}

func (s *UserService) InitializeIndexes() error {
    return s.userRepo.CreateIndexes()
}