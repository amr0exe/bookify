package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/amr0exe/bookify/internal/core"
	"github.com/amr0exe/bookify/internal/middleware"
	"github.com/amr0exe/bookify/internal/tokener"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service   AuthService
	responder core.Responder
	jwtSecret string
}

func NewAuthHandler(service AuthService, responder core.Responder, jwtSecret string) *AuthHandler {
	return &AuthHandler{service: service, responder: responder, jwtSecret: jwtSecret}
}

func (h *AuthHandler) RegisterRoute(r *gin.RouterGroup) {
	accounts := r.Group("/auth")
	{
		accounts.POST("/register", h.Register)
		accounts.POST("/login", h.Login)
		accounts.POST("/refresh", h.RefreshIt)
		accounts.GET("/test", middleware.Authorize(h.jwtSecret), h.TestAuth)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req CreateAccount

	// parse JSON req.body
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind registration request", "warning", err)
		//h.responder.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Please check your input details and try again.", nil)
		h.responder.ValidationError(c, err)
		return
	}

	// create user account and store refresh-token
	accModel, err := h.service.CreateAccount(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			h.responder.Error(c, http.StatusConflict, "EMAIL_EXISTS", "Email already in use.", nil)
			return
		}

		slog.Error("account creation failed", "error", err)
		h.responder.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong. Please try again later.", nil)
		return
	}

	// craft json response, with user-info, tokens
	accResp := ToAccountResponse(accModel.Account, accModel.Tokens)
	h.responder.SuccessWithMeta(c, http.StatusCreated, accResp, "Registration Successfull!!")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginAccount

	// parse JSON req.body
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind login request", "warning", err)
		//h.responder.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Please check your input details and try again.", nil)
		h.responder.ValidationError(c, err)
		return
	}

	usr, tokens, err := h.service.LoginToAccount(c, req.Email)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			h.responder.Error(c, http.StatusConflict, "ACC_NOT_FOUND", "Account not found.", nil)
			return
		}

		slog.Error("Login failed.", "error", err)
		h.responder.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong. Please try again later.", nil)
		return
	}

	accResp := ToAccountResponse(usr, tokens)
	h.responder.SuccessWithMeta(c, http.StatusOK, accResp, "Login Successfull!!")
}

func (h *AuthHandler) RefreshIt(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind refresh_token request", "warning", err)
		h.responder.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Please check your input details and try again.", nil)
		return
	}

	claims, err := tokener.ValidateToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		slog.Warn("failed in validating token", "error", err)
		h.responder.Error(c, http.StatusBadRequest, "INVALID_INPUT", "failed validating token", nil)
		return
	}

	tknPair, err := h.service.RefreshTheToken(c, req.RefreshToken, claims)
	if err != nil {
		slog.Warn("refreshing the token failed", "error", err)
		h.responder.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong. Please try again later.", nil)
		return
	}

	h.responder.Success(c, http.StatusCreated, tknPair)
}

func (h *AuthHandler) TestAuth(c *gin.Context) {
	h.responder.Success(c, http.StatusFound, gin.H{"greet": "endpoint is hit"})
}
