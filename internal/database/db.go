package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("환경변수 DB_DSN이 설정되지 않았습니다.")
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("DB 핑 실패: %v", err)
	}
}

// GetDB returns the db connection
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("DB가 초기화되지 않았습니다. InitDB()를 먼저 호출하세요.")
	}
	return db
}
