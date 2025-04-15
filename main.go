package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kangjiungdev/ai-character-chat/backend/controllers"
)

func main() {
	r := gin.Default()

	// CORS 설정
	r.Use(cors.Default())

	// .env 로드
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env 파일을 로드하지 못했습니다. (무시하고 계속 진행)")
	}

	// 나중에 캐릭터 갖고 올 때

	// DB 초기화
	// database.InitDB()

	// 예시용 출력
	// db := database.GetDB()
	// fmt.Println("DB 인스턴스:", db)

	r.POST("/api/get-characters", controllers.HomeHandler)

	r.POST("/api/login", controllers.LoginHandler)
	r.POST("/api/signup", controllers.SignupHandler)

	r.Run(":4000") // 기본 포트 4000
}
