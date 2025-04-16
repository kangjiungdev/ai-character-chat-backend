package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("환경변수 DB_DSN이 설정되지 않았습니다.")
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("DB 객체 가져오기 실패: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("DB 연결 성공")
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("DB가 초기화되지 않았습니다. InitDB()를 먼저 호출하세요.")
	}
	return db
}
