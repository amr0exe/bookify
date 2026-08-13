package business

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/amr0exe/bookify/internal/core"
	"github.com/amr0exe/bookify/internal/db/models"
	"github.com/amr0exe/bookify/internal/middleware"
	"github.com/gin-gonic/gin"
)

type BusinessHandler struct {
	service   BusinessService
	res       core.Responder
	jwtSecret string
}

func NewBusinessHandler(service BusinessService, res core.Responder, jwtSecret string) *BusinessHandler {
	return &BusinessHandler{service: service, res: res, jwtSecret: jwtSecret}
}

func (h *BusinessHandler) RegisterRoute(r *gin.RouterGroup) {
	business := r.Group("/business")
	{
		business.POST("", middleware.Authorize(h.jwtSecret), h.CreateNewBusinessProfile)
	}
}

func (h *BusinessHandler) CreateNewBusinessProfile(c *gin.Context) {
	var req CreateBusiness

	if err := c.ShouldBindJSON(&req); err != nil {
		h.res.ValidationError(c, err)
		return
	}

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

	if err := h.service.CreateProfile(c, accID, &req); err != nil {
		if errors.Is(err, ErrBusinessAlreadyExists) {
			h.res.Error(c, http.StatusConflict, "BAD_REQUEST", "Creation Failed!! Profile already exists.", nil)
			return
		}

		slog.Warn("Failed creating business profile", "error", err)
		h.res.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "something went wrong on our side, please try again", nil)
		return
	}

	h.res.SuccessWithMeta(c, http.StatusCreated, req, "Business Profile created successfully")
}
