package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

func MeHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	fmt.Println("userID", userID.(string))

	if userID == nil {
		models.SendAPIResponse(c, http.StatusUnauthorized, models.APIResponse[string]{
			Status:  "error",
			Data:    nil,
			Message: "로그인 필요",
		})
		return
	}

	type myData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	data := myData{
		ID:   userID.(string),
		Name: "형님",
	}
	models.SendAPIResponse(c, http.StatusOK, models.APIResponse[myData]{
		Status:  "success",
		Data:    models.Ptr(data),
		Message: "로그인 됨",
	})
}
