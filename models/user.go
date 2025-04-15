package models

import (
	"database/sql"
	"fmt"
	"strings"
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
		// 중복 ID 에러인지 확인
		if isDuplicateEntryError(err) {
			return fmt.Errorf("해당 아이디를 가진 유저가 이미 존재합니다")
		}
		return fmt.Errorf("요청을 처리하는 중 문제가 발생했습니다. 나중에 다시 시도해주세요")
	}

	fmt.Println("유저 생성 완료")
	return nil
}

func isDuplicateEntryError(err error) bool {
	// MySQL의 duplicate entry 에러코드: 1062
	return err != nil &&
		// err 문자열로 체크
		(sqlErrorContains(err, "1062") || sqlErrorContains(err, "Duplicate entry"))
}

func sqlErrorContains(err error, substr string) bool {
	return err != nil &&
		strings.Contains(err.Error(), substr)
}
