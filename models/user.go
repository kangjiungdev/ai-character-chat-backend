package models

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserFormRequest struct {
	ID          string `json:"id" binding:"required" form:"user-id"`
	Password    string `json:"password" binding:"required" form:"password"`
	Name        string `json:"name" binding:"required" form:"name"`
	PhoneNumber string `json:"phone_number" binding:"required" form:"phone-number"`
	BirthDate   string `json:"birth_date" binding:"required" form:"birth-date"`
}

type User struct {
	ID          string
	Password    string
	Name        string
	PhoneNumber string
	BirthDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func CreateUser(db *sql.DB, user *User) error {
	// 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("비밀번호 해싱 실패: %w", err)
	}

	query := `
		INSERT INTO users (id, password, name, phone_number, birth_date)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, user.ID, string(hashedPassword), user.Name, user.PhoneNumber, user.BirthDate)
	if err != nil {
		return fmt.Errorf("유저 생성 실패: %w", err)
	}

	fmt.Println("✅ 유저 생성 완료")
	return nil
}
