package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

func GetCharacters(c *gin.Context) {
	models.SendAPIResponse(c, http.StatusOK, models.APIResponse[string]{
		Status:  "success",
		Data:    models.Ptr("home data"),
		Message: "get-characters",
	})
}
