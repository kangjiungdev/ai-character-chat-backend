package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
)

func LoginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": models.APIResponse[string]{
			Status:  "success",
			Data:    "login data",
			Message: "login success",
		},
	})
}
