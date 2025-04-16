package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserFormRequest struct {
	ID          string `json:"id" binding:"required" form:"user-id"`
	Password    string `json:"password" binding:"required" form:"password"`
	Name        string `json:"name" binding:"required" form:"name"`
	PhoneNumber string `json:"phone_number" binding:"required" form:"phone-number"`
	BirthDate   string `json:"birth_date" binding:"required" form:"birth-date"`
}

func (u *UserFormRequest) UserFormValidation() error {
	inputRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !inputRegex.MatchString(u.ID) {
		return fmt.Errorf("아이디는 영어와 숫자만 입력 가능합니다")
	}
	if utf8.RuneCountInString(u.ID) < 6 || utf8.RuneCountInString(u.ID) > 15 {
		return fmt.Errorf("아이디는 6자~15자여야 합니다")
	}

	if !inputRegex.MatchString(u.Password) {
		return fmt.Errorf("비밀번호는 영어와 숫자만 입력 가능합니다")
	}
	if utf8.RuneCountInString(u.Password) < 8 || utf8.RuneCountInString(u.Password) > 20 {
		return fmt.Errorf("비밀번호는 8자~20자여야 합니다")
	}

	phoneRegex := regexp.MustCompile(`^\d{3}-\d{4}-\d{4}$`)
	if !phoneRegex.MatchString(u.PhoneNumber) {
		return fmt.Errorf("전화번호 형식이 올바르지 않습니다")
	}

	_, err := time.Parse("2006-01-02", u.BirthDate)
	if err != nil {
		return fmt.Errorf("생년월일 형식이 올바르지 않습니다")
	}
	return nil
}

type User struct {
	ID          string
	Password    string
	Name        string
	PhoneNumber string
	BirthDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func CreateUser(db *sql.DB, user *User) (int, error) {
	// 비밀번호 해시
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// 비밀번호 해시 에러
		fmt.Println("비밀번호 해시 에러", err)
		return http.StatusInternalServerError, fmt.Errorf("요청을 처리하는 중 문제가 발생했습니다. 나중에 다시 시도해주세요")
	}

	query := `
		INSERT INTO users (id, password, name, phone_number, birth_date)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, user.ID, string(hashedPassword), user.Name, user.PhoneNumber, user.BirthDate)
	if err != nil {
		// 중복 ID 에러인지 확인
		sqlErrCode := checkSqlError(err)
		if sqlErrCode == 1062 {
			return http.StatusConflict, fmt.Errorf("해당 아이디를 가진 유저가 이미 존재합니다")
		}
		return http.StatusInternalServerError, fmt.Errorf("요청을 처리하는 중 문제가 발생했습니다. 나중에 다시 시도해주세요")
	}

	return http.StatusCreated, nil
}

func DeleteUser(db *sql.DB, userID string) error {
	now := time.Now()
	query := `UPDATE users SET deleted_at = ? WHERE id = ?`
	_, err := db.Exec(query, now, userID)
	if err != nil {
		return fmt.Errorf("계정 삭제 중 오류가 발생했습니다. 잠시 후 다시 시도해 주세요")
	}
	return nil
}

func checkSqlError(err error) uint16 {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number
	}
	return 0
}
