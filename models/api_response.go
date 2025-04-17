package models

import (
	"github.com/gin-gonic/gin"
)

type APIResponse[T any] struct {
	Status  string `json:"status"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

func SendAPIResponse[T any](c *gin.Context, statusCode int, apiResponse APIResponse[T]) {
	c.JSON(statusCode, gin.H{
		"data": apiResponse,
	})
}
