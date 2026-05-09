package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"omnievent-backend/cmd"
	"omnievent-backend/internal/api"
	"omnievent-backend/internal/middleware"
	"omnievent-backend/internal/repository"
	"omnievent-backend/internal/service"
	apipkg "omnievent-backend/pkg/api"
	"omnievent-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

var (
	Version    = "1.0.0"
	BuildTime  = "unknown"
	GitCommit  = "unknown"
	ServerPort = "8080"
)

func init() {
	flag.StringVar(&ServerPort, "port", "8080", "server port")
}

func main() {
	flag.Parse()

	if err := cmd.InitConfig(); err != nil {
		log.Fatal("Failed to init config: ", err)
	}

	if err := cmd.InitDatabase(); err != nil {
		log.Fatal("Failed to init database: ", err)
	}

	if err := cmd.InitRedis(); err != nil {
		log.Fatal("Failed to init redis: ", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(requestLogger())
	r.Use(corsMiddleware())
	r.Use(requestIDMiddleware())

	// Health check endpoint
	r.GET("/health", healthCheck)
	r.GET("/version", versionInfo)

	// Initialize services
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService()
	tokenService := service.NewTokenService(userRepo)

	userHandler := api.NewUserHandler(userService, tokenService)
	avatarHandler := api.NewAvatarHandler(userService)
	authHandler := api.NewAuthHandler(userService, tokenService)
	tokenHandler := api.NewTokenHandler(tokenService)
	tokenHandler.SetUserService(userService)
	forgetPasswordHandler := api.NewForgetPasswordHandler(userService, tokenService)

	// Public routes
	r.POST("/api/register.json", apipkg.BindApi(authHandler.Register))
	r.POST("/api/authorize.json", apipkg.BindApiWithTokenUpdate(authHandler.Login))
	r.POST("/api/forget_password/request.json", apipkg.BindApi(forgetPasswordHandler.UserForgetPasswordRequest))
	r.POST("/api/forget_password/reset/by_token.json", apipkg.BindApi(forgetPasswordHandler.UserResetPassword))

	// Protected routes
	v1 := r.Group("/api/v1")
	v1.Use(middleware.JWTAuthorization())
	{
		v1.GET("/users/profile/get.json", apipkg.BindApi(userHandler.GetProfile))
		v1.POST("/users/profile/update.json", apipkg.BindApiWithTokenUpdate(userHandler.UpdateProfile))
		v1.POST("/users/avatar/update.json", apipkg.BindApiWithTokenUpdate(userHandler.UpdateAvatar))
		v1.POST("/users/avatar/remove.json", apipkg.BindApiWithTokenUpdate(userHandler.RemoveAvatar))
		v1.POST("/users/avatar/upload.json", apipkg.BindApiWithTokenUpdate(avatarHandler.UploadAvatar))
		v1.GET("/avatars/:fileName", apipkg.BindAvatarApi(avatarHandler.GetAvatar))
		v1.POST("/users/verify_email/resend.json", apipkg.BindApi(userHandler.SendVerifyEmail))

		v1.GET("/tokens/list.json", apipkg.BindApi(tokenHandler.ListTokens))
		v1.POST("/tokens/refresh.json", apipkg.BindApiWithTokenUpdate(tokenHandler.RefreshToken))
		v1.POST("/tokens/revoke.json", apipkg.BindApi(tokenHandler.RevokeToken))
		v1.POST("/tokens/revoke_all.json", apipkg.BindApi(tokenHandler.RevokeAllTokens))
		v1.POST("/tokens/revoke_current.json", apipkg.BindApi(tokenHandler.RevokeCurrentToken))
		v1.POST("/tokens/generate_api.json", apipkg.BindApi(tokenHandler.GenerateAPIToken))
		v1.POST("/tokens/generate_mcp.json", apipkg.BindApi(tokenHandler.GenerateMCPToken))
	}

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ServerPort
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("OmniEvent server starting on port %s", port)
		log.Printf("Version: %s, Build: %s, Commit: %s", Version, BuildTime, GitCommit)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Cleanup
	if err := cmd.CloseDatabase(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	if err := cmd.CloseRedis(); err != nil {
		log.Printf("Error closing redis: %v", err)
	}

	log.Println("Server exited gracefully")
}

// Recovery middleware with logging
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, string(debug.Stack()))
				response.Error(c, http.StatusInternalServerError, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// Request logger middleware
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID, _ := c.Get("RequestID")

		log.Printf("[%s] %s %s %s %d %s",
			requestID,
			c.Request.Method,
			path,
			query,
			status,
			latency,
		)
	}
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Request ID middleware
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID(c)
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func generateRequestID(c *gin.Context) string {
	ip := c.ClientIP()
	port := c.Request.RemoteAddr
	if idx := strings.LastIndex(port, ":"); idx > 0 {
		port = port[idx+1:]
	}
	return ip + ":" + port + ":" + strconv.FormatInt(time.Now().UnixNano(), 36)
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// Version info endpoint
func versionInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    Version,
		"build_time": BuildTime,
		"git_commit": GitCommit,
	})
}
