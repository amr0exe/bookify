package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/amr0exe/bookify/internal/tokener"
	"github.com/gin-gonic/gin"
)

// Middleware to to authorize with tokens

func Authorize(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortUnauthorized(c, "no token found")
			return
		}

		scheme, token, found := strings.Cut(authHeader, " ")
		if !found || strings.ToLower(scheme) != "bearer" || token == "" {
			abortUnauthorized(c, "no token found")
			return
		}

		claims, err := tokener.ValidateToken(token, jwtSecret)
		if err != nil {
			if errors.Is(err, tokener.ErrExpiredToken) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success":    false,
					"error_code": "TOKEN_EXPIRED",
					"data": gin.H{
						"message": "token has expired",
					},
				})
				c.Abort()
				return
			}
			abortUnauthorized(c, "Invalid token token")
			return
		}

		if claims.TokenType != "access" {
			abortUnauthorized(c, "Invalid Token type")
			return
		}

		// set context values for downstream handlers
		c.Set("token", token)
		c.Set("account_id", claims.AccountId)
		c.Set("account_role", claims.Role)

		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"data": gin.H{
			"message": msg,
		},
	})
	c.Abort()
}
