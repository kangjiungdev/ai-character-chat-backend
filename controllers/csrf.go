package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kangjiungdev/ai-character-chat/backend/models"
	csrf "github.com/utrack/gin-csrf"
)

func CsrfTokenHandler(c *gin.Context) {
	fmt.Println("CSRF token handler called")
	if !strings.Contains(c.GetHeader("Accept"), "application/json") {
		models.SendAPIResponse(c, http.StatusNotFound, models.APIResponse[string]{
			Status:  "error",
			Data:    nil,
			Message: "Not Found",
		})
		return
	}

	token := csrf.GetToken(c)

	models.SendAPIResponse(c, http.StatusOK, models.APIResponse[gin.H]{
		Status:  "success",
		Data:    models.Ptr(gin.H{"csrf_token": token}),
		Message: "CSRF token generated successfully",
	})
}
