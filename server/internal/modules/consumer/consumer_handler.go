package consumer

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/amr0exe/bookify/internal/core"
	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/amr0exe/bookify/internal/middleware"
	"github.com/gin-gonic/gin"
)

type ConsumerHandler struct {
	service   ConsumerService
	res       core.Responder
	jwtSecret string
}

func NewConsumerHandler(service ConsumerService, res core.Responder, jwtSecret string) *ConsumerHandler {
	return &ConsumerHandler{service: service, res: res, jwtSecret: jwtSecret}
}

func (h *ConsumerHandler) RegisterRoute(r *gin.RouterGroup) {
	consumer := r.Group("/consumer")
	{
		consumer.POST("", middleware.Authorize(h.jwtSecret), h.CreateProfile)
	}
}

func (h *ConsumerHandler) CreateProfile(c *gin.Context) {
	var req CreateConsumer

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind CreateProfile reqest", "warning", err)
		h.res.ValidationError(c, err)
		return
	}

	role, err := core.GetAccountRole(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "INSUFICIENT_CONTEXT", "request doesn't contain required context", nil)
		return
	}

	if role == models.RoleBusiness {
		h.res.Error(c, http.StatusForbidden, "FORBIDDEN", "only consumer can create consumer profile", nil)
		return
	}

	accountID, err := core.GetAccountId(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "INSUFICIENT_CONTEXT", "request doesn't contain required context", nil)
		return
	}

	if err := h.service.CreateConsumer(c, accountID, &req); err != nil {
		if errors.Is(err, ErrConsumerAlreadyExists) {
			slog.Error("failed creating consumer profile", "warning", err)
			h.res.Error(c, http.StatusConflict, "CONSUMER_ALREADY_EXISTS", "A consumer profile already exists for this account", nil)
			return
		}

		slog.Error("failed creating consumer profile", "error", err)
		h.res.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed creating consumer profile", nil)
		return
	}

	h.res.SuccessWithMeta(c, http.StatusCreated, req, "ConsumerProfile Creation successfull")
}
