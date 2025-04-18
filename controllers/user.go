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

	models.SendAPIResponse(c, http.StatusOK, models.APIResponse[string]{
		Status:  "success",
		Data:    models.Ptr(userID.(string)),
		Message: "로그인 됨",
	})
}
