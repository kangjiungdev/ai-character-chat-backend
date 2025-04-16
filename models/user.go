package models

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string `gorm:"primaryKey"`
	Password    string
	Name        string
	PhoneNumber string
	BirthDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

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
		return errors.New("아이디는 영어와 숫자만 입력 가능합니다")
	}
	if utf8.RuneCountInString(u.ID) < 6 || utf8.RuneCountInString(u.ID) > 15 {
		return errors.New("아이디는 6자~15자여야 합니다")
	}

	if !inputRegex.MatchString(u.Password) {
		return errors.New("비밀번호는 영어와 숫자만 입력 가능합니다")
	}
	if utf8.RuneCountInString(u.Password) < 8 || utf8.RuneCountInString(u.Password) > 20 {
		return errors.New("비밀번호는 8자~20자여야 합니다")
	}

	phoneRegex := regexp.MustCompile(`^\d{3}-\d{4}-\d{4}$`)
	if !phoneRegex.MatchString(u.PhoneNumber) {
		return errors.New("전화번호 형식이 올바르지 않습니다")
	}

	_, err := time.Parse("2006-01-02", u.BirthDate)
	if err != nil {
		return errors.New("생년월일 형식이 올바르지 않습니다")
	}
	return nil
}

func CreateUser(db *gorm.DB, user *User) (int, error) {
	// 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("비밀번호 해시 에러:", err)
		return http.StatusInternalServerError, errors.New("요청을 처리하는 중 문제가 발생했습니다. 나중에 다시 시도해주세요")
	}
	user.Password = string(hashedPassword)

	// 유저 생성
	err = db.Create(user).Error
	if err != nil {
		// 중복 ID 에러 확인
		if isDuplicateKeyError(err) {
			return http.StatusConflict, errors.New("해당 아이디를 가진 유저가 이미 존재합니다")
		}
		return http.StatusInternalServerError, errors.New("요청을 처리하는 중 문제가 발생했습니다. 나중에 다시 시도해주세요")
	}

	return http.StatusCreated, nil
}

func DeleteUser(db *gorm.DB, userID string) error {
	err := db.Where("id = ?", userID).Delete(&User{}).Error
	if err != nil {
		return errors.New("계정 삭제 중 오류가 발생했습니다. 잠시 후 다시 시도해 주세요")
	}
	return nil
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
