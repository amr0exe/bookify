package router

import (
	"github.com/amr0exe/bookify/internal/config"
	"github.com/amr0exe/bookify/internal/core"
	"github.com/amr0exe/bookify/internal/modules/auth"
	"github.com/amr0exe/bookify/internal/modules/business"
	"github.com/amr0exe/bookify/internal/modules/consumer"
	"github.com/amr0exe/bookify/internal/modules/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUp(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	responder := core.NewResponder()

	// auth
	accountRepo := auth.NewAccountRepo(db)
	refreshTokenRepo := auth.NewRefreshTknRepository(db)
	authService := auth.NewAuthService(accountRepo, refreshTokenRepo, cfg.JWT_SECRET)
	authHandler := auth.NewAuthHandler(authService, responder, cfg.JWT_SECRET)

	authHandler.RegisterRoute(api)

	// consumer
	cRepo := consumer.NewConsumerRepository(db)
	cService := consumer.NewConsumerService(cRepo)
	cHandler := consumer.NewConsumerHandler(cService, responder, cfg.JWT_SECRET)

	cHandler.RegisterRoute(api)

	// business
	bRepo := business.NewBusinessRepository(db)
	bService := business.NewBusinessService(bRepo)
	bHandler := business.NewBusinessHandler(bService, responder, cfg.JWT_SECRET)

	bHandler.RegisterRoute(api)

	// business_service
	sRepo := service.NewServiceRepository(db)
	sService := service.NewServiceService(sRepo)
	sHandler := service.NewServiceHandler(sService, responder, cfg.JWT_SECRET)

	sHandler.RegisterRoute(api)

	return r
}
