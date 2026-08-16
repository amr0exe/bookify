package service

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/amr0exe/bookify/internal/core"
	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/amr0exe/bookify/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	service   ServiceService
	res       core.Responder
	jwtSecret string
}

func NewServiceHandler(service ServiceService, res core.Responder, jwtSecret string) *ServiceHandler {
	return &ServiceHandler{service: service, res: res, jwtSecret: jwtSecret}
}

func (h *ServiceHandler) RegisterRoute(r *gin.RouterGroup) {
	sr := r.Group("/business")
	{
		sr.GET("/service/all", middleware.Authorize(h.jwtSecret), h.GetAllServices)
		sr.POST("/service", middleware.Authorize(h.jwtSecret), h.CreateNewService)
		sr.PUT("/service/:id", middleware.Authorize(h.jwtSecret), h.UpdateService)
		sr.DELETE("/service/:id", middleware.Authorize(h.jwtSecret), h.DeleteService)
	}
}

func (h *ServiceHandler) GetAllServices(c *gin.Context) {
	accID, err := core.GetAccountId(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}
	role, err := core.GetAccountRole(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}

	if role != models.RoleBusiness {
		h.res.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "you are not permitted this action.", nil)
		return
	}

	services, err := h.service.GetAllServices(c, accID)
	if err != nil {
		if errors.Is(err, ErrNoBusinessFound) {
			h.res.Error(c, http.StatusNotFound, "NOT_FOUND", "No business account found for this accountID", nil)
			return
		}

		slog.Error("failed fetching services", "error", err)
		h.res.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Something went wrong here, Please try again", nil)
		return
	}

	h.res.SuccessWithMeta(c, http.StatusOK, services, "all services fetched successfully")
}

func (h *ServiceHandler) CreateNewService(c *gin.Context) {
	var req CreateService

	if err := c.ShouldBindJSON(&req); err != nil {
		h.res.ValidationError(c, err)
		return
	}

	// take Id, role from token, for context to downstream handlers
	accID, err := core.GetAccountId(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}
	role, err := core.GetAccountRole(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}

	if role != models.RoleBusiness {
		h.res.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "you are not permitted this action.", nil)
		return
	}

	info, err := h.service.CreateService(c, accID, &req)
	if err != nil {
		if errors.Is(err, ErrNoBusinessFound) {
			h.res.Error(c, http.StatusNotFound, "NOT_FOUND", "No business account found for this accountID", nil)
			return
		}

		slog.Error("failed creating service", "error", nil)
		h.res.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Something went wrong here, Please try again", nil)
		return
	}

	h.res.SuccessWithMeta(c, http.StatusCreated, info, "Service created successfully")
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	serviceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}

	var req UpdateService

	if err := c.ShouldBindJSON(&req); err != nil {
		h.res.ValidationError(c, err)
		return
	}

	// take Id, role from token, for context to downstream handlers
	accID, err := core.GetAccountId(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}
	role, err := core.GetAccountRole(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}

	if role != models.RoleBusiness {
		h.res.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "you are not permitted this action.", nil)
		return
	}

	info, err := h.service.UpdateService(c, accID, serviceId, &req)
	if err != nil {
		if errors.Is(err, ErrNoBusinessFound) {
			h.res.Error(c, http.StatusNotFound, "NOT_FOUND", "No business account found for this accountID", nil)
			return
		}
		if errors.Is(err, ErrNoServiceFound) {
			h.res.Error(c, http.StatusNotFound, "NOT_FOUND", "No service found for this id", nil)
			return
		}

		slog.Error("failed updating service", "error", err)
		h.res.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Something went wrong here, Please try again", nil)
		return
	}

	h.res.SuccessWithMeta(c, http.StatusOK, info, "Service updated successfully")
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	serviceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}

	// take Id, role from token, for context to downstream handlers
	accID, err := core.GetAccountId(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}
	role, err := core.GetAccountRole(c)
	if err != nil {
		h.res.Error(c, http.StatusBadRequest, "BAD_REQUEST", "request doesn't contain required context", nil)
		return
	}

	if role != models.RoleBusiness {
		h.res.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "you are not permitted this action.", nil)
		return
	}

	if err := h.service.DeleteService(c, accID, serviceId); err != nil {
		if errors.Is(err, ErrNoBusinessFound) {
			h.res.Error(c, http.StatusNotFound, "NOT_FOUND", "No business account found for this accountID", nil)
			return
		}

		slog.Error("failed deleting service", "error", nil)
		h.res.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Something went wrong here, Please try again", nil)
		return
	}

	h.res.Success(c, http.StatusOK, "Service deleted successfully")
}
