// Rental Management System API
// @title 租房管理系统 API
// @version 1.0
// @description 面向个人房东的租房管理后端系统，支持房源管理、租客管理、水电气收费、租金账单等功能。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "rental-management/api/swagger" // swagger docs
	"rental-management/internal/config"
	"rental-management/internal/handler"
	"rental-management/internal/middleware"
	"rental-management/internal/model"
	"rental-management/internal/pkg/db"
	"rental-management/internal/pkg/logger"
	"rental-management/internal/repository"
	"rental-management/internal/service"
)

func main() {
	// 加载配置
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	if err := config.Init(configPath); err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	cfg := config.Get()

	// 初始化日志
	log, err := logger.InitLogger(&cfg.Log)
	if err != nil {
		panic(fmt.Sprintf("Failed to init logger: %v", err))
	}
	logger.Logger = log

	// 初始化数据库
	database, err := db.InitMySQL(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database", zap.Error(err))
	}
	log.Info("Database connected successfully")

	// 自动迁移
	if err := database.AutoMigrate(
		&model.User{},
		&model.FeeRate{},
		&model.Property{},
		&model.Room{},
		&model.Tenant{},
		&model.Contract{},
		&model.MeterReading{},
		&model.Bill{},
		&model.Payment{},
	); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}
	log.Info("Database migrated successfully")

	// 初始化Redis
	_, err = db.InitRedis(&cfg.Redis)
	if err != nil {
		log.Warn("Redis connection failed, running without cache", zap.Error(err))
	} else {
		log.Info("Redis connected successfully")
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化Repository
	userRepo := repository.NewUserRepository(database)
	propertyRepo := repository.NewPropertyRepository(database)
	roomRepo := repository.NewRoomRepository(database)
	tenantRepo := repository.NewTenantRepository(database)
	contractRepo := repository.NewContractRepository(database)
	billRepo := repository.NewBillRepository(database)

	// 初始化Service
	authService := service.NewAuthService(userRepo, propertyRepo, roomRepo, tenantRepo, billRepo)
	propertyService := service.NewPropertyService(propertyRepo, roomRepo)
	tenantService := service.NewTenantService(tenantRepo, contractRepo, roomRepo, userRepo)
	billService := service.NewBillService(billRepo, roomRepo, userRepo, contractRepo)

	// 初始化Handler
	authHandler := handler.NewAuthHandler(authService)
	propertyHandler := handler.NewPropertyHandler(propertyService)
	tenantHandler := handler.NewTenantHandler(tenantService)
	billHandler := handler.NewBillHandler(billService)

	// 创建路由
	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证路由(无需登录)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 需要登录的路由
		protected := api.Use(middleware.JWTAuth())
		{
			// 用户信息
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.PUT("/auth/password", authHandler.ChangePassword)

			// 统计
			protected.GET("/statistics", authHandler.GetStatistics)

			// 房源管理
			protected.POST("/properties", propertyHandler.CreateProperty)
			protected.GET("/properties", propertyHandler.ListProperties)
			protected.GET("/properties/:id", propertyHandler.GetProperty)
			protected.PUT("/properties/:id", propertyHandler.UpdateProperty)
			protected.DELETE("/properties/:id", propertyHandler.DeleteProperty)

			// 房间管理
			protected.POST("/rooms", propertyHandler.CreateRoom)
			protected.GET("/rooms", propertyHandler.ListRooms)
			protected.GET("/rooms/:id", propertyHandler.GetRoom)
			protected.PUT("/rooms/:id", propertyHandler.UpdateRoom)
			protected.DELETE("/rooms/:id", propertyHandler.DeleteRoom)
			protected.PUT("/rooms/:id/status", propertyHandler.UpdateRoomStatus)

			// 租客管理
			protected.POST("/tenants", tenantHandler.CreateTenant)
			protected.GET("/tenants", tenantHandler.ListTenants)
			protected.GET("/tenants/:id", tenantHandler.GetTenant)
			protected.PUT("/tenants/:id", tenantHandler.UpdateTenant)
			protected.DELETE("/tenants/:id", tenantHandler.DeleteTenant)

			// 合同管理
			protected.POST("/contracts", tenantHandler.CreateContract)
			protected.GET("/contracts", tenantHandler.ListContracts)
			protected.GET("/contracts/:id", tenantHandler.GetContract)
			protected.PUT("/contracts/:id", tenantHandler.UpdateContract)
				protected.PUT("/contracts/:id/terminate", tenantHandler.TerminateContract)
				protected.DELETE("/contracts/:id", tenantHandler.DeleteContract)

			// 抄表管理
			protected.POST("/meter-readings", billHandler.CreateMeterReading)
			protected.GET("/meter-readings", billHandler.ListMeterReadings)
			protected.DELETE("/meter-readings/:id", billHandler.DeleteMeterReading)

			// 账单管理
			protected.POST("/bills", billHandler.CreateBill)
			protected.GET("/bills", billHandler.ListBills)
				protected.GET("/bills/preview", billHandler.PreviewBill)
			protected.GET("/bills/statistics", billHandler.GetBillStatistics)
			protected.GET("/bills/monthly-statistics", billHandler.GetMonthlyStatistics)
			protected.GET("/bills/:id", billHandler.GetBill)
			protected.PUT("/bills/:id/pay", billHandler.PayBill)
			protected.DELETE("/bills/:id", billHandler.DeleteBill)

			// 系统配置
			protected.GET("/fee-rates", billHandler.GetFeeRate)
			protected.PUT("/fee-rates", billHandler.UpdateFeeRate)
		}
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info("Server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
