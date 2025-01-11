package api

import (
	"net/http"

	"github.com/Vaibhavsahu2810/redis-go/internal/queue"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type SendEmailRequest struct {
	To           string                 `json:"to" binding:"required,email"`
	Subject      string                 `json:"subject" binding:"required"`
	TemplateName string                 `json:"templateName" binding:"required"`
	Data         map[string]interface{} `json:"data" binding:"required"`
}

func RegisterHandlers(router *gin.Engine, redisClient *redis.Client) {
	router.GET("/health", healthCheck)

	api := router.Group("/api")
	{
		api.POST("/send", sendEmailHandler(redisClient))
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func sendEmailHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SendEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		task := queue.EmailTask{
			To:           req.To,
			Subject:      req.Subject,
			TemplateName: req.TemplateName,
			Data:         req.Data,
		}

		if err := queue.EnqueueEmail(c.Request.Context(), redisClient, task); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error enqueueing email: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "email queued",
		})
	}
}
