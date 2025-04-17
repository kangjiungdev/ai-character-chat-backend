package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

func HomeHandler(c *gin.Context) {
	models.SendAPIResponse(c, http.StatusBadRequest, models.APIResponse[string]{
		Status:  "success",
		Data:    "home data",
		Message: "get-characters",
	})
}
