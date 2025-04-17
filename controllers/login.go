package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/internal/database"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

type LoginFormRequest struct {
	ID       string `json:"id" binding:"required" form:"user-id"`
	Password string `json:"password" binding:"required" form:"password"`
}

func LoginHandler(c *gin.Context) {
	req := LoginFormRequest{}
	err := c.Bind(&req)
	if err != nil {
		models.SendAPIResponse(c, http.StatusBadRequest, models.APIResponse[string]{Status: "error", Data: "", Message: "Form 형식이 올바르지 않습니다"})
		return
	}
	database.InitDB()
	db := database.GetDB()

	user := &models.User{
		ID:       req.ID,
		Password: req.Password,
	}

	statusCode, err := models.Login(db, user)

	status := "success"
	message := "Login success"

	if err != nil {
		status = "error"
		message = err.Error()
		models.SendAPIResponse(c, statusCode, models.APIResponse[string]{Status: status, Data: "", Message: message})
	}
	models.SendAPIResponse(c, statusCode, models.APIResponse[models.User]{Status: status, Data: *user, Message: message})
}
