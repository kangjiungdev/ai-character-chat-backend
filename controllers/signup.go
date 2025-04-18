package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/internal/database"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

func SignupHandler(c *gin.Context) {
	req := models.UserFormRequest{}
	err := c.Bind(&req)
	if err != nil {
		models.SendAPIResponse(c, http.StatusBadRequest, models.APIResponse[string]{Status: "error", Data: nil, Message: "입력하신 정보가 올바르지 않습니다"})
		return
	}

	err = req.UserFormValidation()
	if err != nil {
		models.SendAPIResponse(c, http.StatusBadRequest, models.APIResponse[string]{Status: "error", Data: nil, Message: err.Error()})
		return
	}

	birthDate, _ := time.Parse("2006-01-02", req.BirthDate)
	user := &models.User{
		ID:          req.ID,
		Password:    req.Password,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		BirthDate:   birthDate,
	}
	database.InitDB()
	db := database.GetDB()
	statusCode, err := models.CreateUser(db, user)

	// 문제 해결: 이제 유저 삭제하면 삭제된 유저 ID를 다시 생성할 수 없음. 유저 삭제는 DB에서 삭제하는 것이 아니라
	// deleted_at에 현재 시간을 넣는 것으로 처리함. (soft delete)
	// 그래서 ID가 중복되는 경우는 없을 것임.

	status := "success"
	message := "Signup Success"

	if err != nil {
		status = "error"
		message = err.Error()
	}
	models.SendAPIResponse(c, statusCode, models.APIResponse[string]{Status: status, Data: nil, Message: message})

}
