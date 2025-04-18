package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/controllers"
	csrf "github.com/utrack/gin-csrf"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	frontendUrl := os.Getenv("FRONTEND_URL")

	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET_KEY")))
	r.Use(sessions.Sessions("aichat_sess", store))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// 라우트 등록
	api := r.Group("/api")
	api.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CSRF_SECRET_KEY"),
		ErrorFunc: func(c *gin.Context) {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
		},
	}))
	api.Use(EnforceJSONOnly())
	{
		api.GET("/csrf-token", controllers.CsrfTokenHandler)
		api.GET("/me", controllers.MeHandler)
		api.POST("/signup", controllers.SignupHandler)
		api.POST("/login", controllers.LoginHandler)
	}

	public := r.Group("/api/public")
	public.Use(EnforceJSONOnly())
	{
		public.GET("/get-characters", controllers.GetCharacters)
	}

	return r
}

func EnforceJSONOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Accept 헤더가 application/json을 포함하지 않으면
		if !strings.Contains(c.GetHeader("Accept"), "application/json") {
			// 강제 404 응답
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Not Found",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
