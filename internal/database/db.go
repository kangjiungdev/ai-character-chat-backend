package database

import (
	"log"
	"os"

	"github.com/kangjiungdev/ai-character-chat/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("환경변수 DB_DSN이 설정되지 않았습니다.")
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	// 마이그레이션 (선택)
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("AutoMigrate 실패: %v", err)
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("DB가 초기화되지 않았습니다. InitDB()를 먼저 호출하세요.")
	}
	return db
}
