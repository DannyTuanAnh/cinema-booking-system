package middleware

import (
	"net/http"
	"strings"

	jwt "cinema.com/demo/pkg/jwt_service"
	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	validator jwt.Validator
}

func NewJWTMiddleware(validator jwt.Validator) *JWTMiddleware {
	return &JWTMiddleware{
		validator: validator,
	}
}

func (m *JWTMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		// Expect: Authorization: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			return
		}

		tokenString := parts[1]

		claims, err := m.validator.Validate(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Inject claims vào context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
