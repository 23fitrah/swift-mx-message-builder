package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"swift-mx-message-builder/worker"
)

func InquiryHandler(pool *worker.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageId := c.Param("messageId")
		if messageId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "messageId is required"})
			return
		}

		result, ok := pool.GetStatus(messageId)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error":      "message id not found",
				"message_id": messageId,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message_id":   result.MessageID,
			"message_type": result.MsgType,
			"status":       result.Status,
			"file_path":    result.FilePath,
			"error":        result.Error,
			"submitted_at": result.SubmittedAt,
			"updated_at":   result.UpdatedAt,
		})
	}
}
