package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "yourmodule/internal/utils"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            // Check cookie
            cookie, err := c.Cookie("token")
            if err == nil {
                tokenString = "Bearer " + cookie
            }
        }
        
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token"})
            c.Abort()
            return
        }
        
        // Remove Bearer prefix
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("is_admin", claims.IsAdmin)
        c.Next()
    }
}

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        isAdmin := c.GetBool("is_admin")
        if !isAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
            c.Abort()
            return
        }
        c.Next()
    }
}
