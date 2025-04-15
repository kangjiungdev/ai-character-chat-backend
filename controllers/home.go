package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": models.APIResponse[string]{
			Status:  "success",
			Data:    "home data",
			Message: "get-characters",
		},
	})
}
