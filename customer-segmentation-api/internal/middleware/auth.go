package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Bearer token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        // Extract token from "Bearer <token>" format
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // Get JWT secret
        jwtSecret := os.Getenv("JWT_SECRET")
        if jwtSecret == "" {
            jwtSecret = "fallback-secret-key"
        }

        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Extract claims and set them in context
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", claims["user_id"])
            c.Set("username", claims["username"])
            c.Set("email", claims["email"])
            c.Set("role", claims["role"])
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        c.Next()
    }
}