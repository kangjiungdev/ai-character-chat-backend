package models

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	ID          string     `gorm:"primaryKey;size:20" json:"id"`
	Password    string     `gorm:"size:60" json:"-"`
	Name        string     `gorm:"size:15" json:"name"`
	PhoneNumber string     `gorm:"size:13" json:"phone_number"`
	BirthDate   time.Time  `json:"birth_date"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func CreateUser(db *gorm.DB, user *User) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("비밀번호 해시 처리 중 오류가 발생했습니다")
	}
	user.Password = string(hashedPassword)

	if err := db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || isDuplicateEntryError(err) {
			return http.StatusConflict, fmt.Errorf("해당 아이디를 가진 유저가 이미 존재합니다")
		}
		return http.StatusInternalServerError, fmt.Errorf("회원가입 처리 중 오류가 발생했습니다")
	}

	return http.StatusCreated, nil
}

func DeleteUser(db *gorm.DB, userID string) error {
	if err := db.Model(&User{}).Where("id = ?", userID).Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("계정 삭제 중 오류가 발생했습니다")
	}
	return nil
}

// gorm.ErrDuplicatedKey가 제대로 안 잡힐 경우 문자열로 체크
func isDuplicateEntryError(err error) bool {
	return err != nil && (contains(err.Error(), "Duplicate entry") || contains(err.Error(), "1062"))
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && (regexp.MustCompile(substr).MatchString(s))
}
