package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"swift-mx-message-builder/worker"
)

// InquiryHandler is a generic "inquiry status message" endpoint used by
// all four services (pacs008/pacs009/pacs002/pacs004). It reports the
// current lifecycle status of a previously generated MX message:
// PENDING -> PROCESSING -> COMPLETED (or FAILED), as tracked by the
// worker pool that actually writes the file to disk in the background.
//
// GET /api/v1/{service}/inquiry/:messageId
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
