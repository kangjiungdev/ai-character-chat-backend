package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/internal/database"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

func SignupHandler(c *gin.Context) {
	req := models.UserFormRequest{}
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": models.APIResponse[string]{
				Status:  "error",
				Data:    "",
				Message: "Form 형식이 올바르지 않습니다",
			},
		})
		return
	}

	// ID는 영어+숫자만 허용 (JS와 동일한 규칙 적용)
	idRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !idRegex.MatchString(req.ID) {
		fmt.Println("아이디는 영어와 숫자만 입력 가능합니다")
		c.JSON(http.StatusBadRequest, gin.H{
			"data": models.APIResponse[string]{
				Status:  "error",
				Data:    "",
				Message: "아이디는 영어와 숫자만 입력 가능합니다",
			},
		})
		return
	}

	if utf8.RuneCountInString(req.ID) < 6 || utf8.RuneCountInString(req.ID) > 15 {
		fmt.Println("아이디는 6자~15자여야 합니다")
		c.JSON(http.StatusBadRequest, gin.H{
			"data": models.APIResponse[string]{
				Status:  "error",
				Data:    "",
				Message: "아이디는 6자~15자여야 합니다",
			},
		})
		return
	}

	if utf8.RuneCountInString(req.Password) < 8 || utf8.RuneCountInString(req.Password) > 20 {
		fmt.Println("비밀번호는 8자~20자여야 합니다")
		c.JSON(http.StatusBadRequest, gin.H{
			"data": models.APIResponse[string]{
				Status:  "error",
				Data:    "",
				Message: "비밀번호는 8자~20자여야 합니다",
			},
		})
		return
	}

	phoneRegex := regexp.MustCompile(`^\d{3}-\d{4}-\d{4}$`)
	if !phoneRegex.MatchString(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": models.APIResponse[string]{
				Status:  "error",
				Data:    "",
				Message: "전화번호 형식이 올바르지 않습니다",
			},
		})
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": models.APIResponse[string]{
				Status:  "error",
				Data:    "",
				Message: "생년월일 형식이 올바르지 않습니다",
			},
		})
		return
	}
	user := &models.User{
		ID:          req.ID,
		Password:    req.Password,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		BirthDate:   birthDate,
	}
	fmt.Println("user", user)
	database.InitDB()
	db := database.GetDB()
	err = models.CreateUser(db, user)
	if err != nil {
		var statusCode int
		if err.Error() == "해당 아이디를 가진 유저가 이미 존재합니다" {
			statusCode = http.StatusConflict
		}
		if err.Error() == "요청을 처리하는 중 문제가 발생했습니다. 나중에 다시 시도해주세요" {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{
			"data": models.APIResponse[string]{
				Status:  "error",
				Data:    "",
				Message: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": models.APIResponse[string]{
			Status:  "success",
			Data:    "",
			Message: "signup success",
		},
	})

}
