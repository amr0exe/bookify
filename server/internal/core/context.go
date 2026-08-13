package core

// Provides value from request context to downstream handlers

import (
	"errors"

	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrAccountIDNotFound = errors.New("account_id not found in context")
	ErrInvalidType       = errors.New("invaid account_id type in context")

	ErrRoleNotFound    = errors.New("account_role not found in context")
	ErrInvalidRoleType = errors.New("invalid account_role type in context")
)

func GetAccountId(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("account_id")
	if !exists {
		return uuid.Nil, ErrAccountIDNotFound
	}

	accID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidType
	}

	return accID, nil
}

func GetAccountRole(c *gin.Context) (models.Role, error) {
	val, exists := c.Get("account_role")
	if !exists {
		return "", ErrRoleNotFound
	}

	role, ok := val.(models.Role)
	if !ok {
		return "", ErrInvalidRoleType
	}

	return role, nil
}
