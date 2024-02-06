package router

import (
	"chat/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func authMiddleware(checkIsAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(authHeader, &user.MyJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*user.MyJWTClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if checkIsAdmin && !claims.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You must have the admin role."})
			c.Abort()
			return
		}

		c.Set("id", claims.ID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}
