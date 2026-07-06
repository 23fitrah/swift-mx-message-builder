package handlers

import (
	"net/http"
	"time"

	"swift-mx-message-builder/dto/pacs002"
	"swift-mx-message-builder/models"
	"swift-mx-message-builder/utils"
	"swift-mx-message-builder/worker"

	"github.com/gin-gonic/gin"
)

// Pacs002Handler builds a pacs.002.001.10 Payment Status Report MX
// document - this is the standard way to answer/confirm the status of a
// previously submitted pacs.008/pacs.009 message - and queues it for
// asynchronous file generation.
func Pacs002Handler(pool *worker.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pacs002.Pacs002Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var reqPayload = req.Payload

		msgId := utils.GenerateMessageID("PACS002")

		var stsRsnInf *models.StatusReasonInformation
		if reqPayload.ReasonCode != "" || reqPayload.AdditionalInfo != "" {
			stsRsnInf = &models.StatusReasonInformation{
				Rsn:      reqPayload.ReasonCode,
				AddtlInf: reqPayload.AdditionalInfo,
			}
		}

		doc := models.Pacs002Document{
			FIToFIPmtStsRpt: models.FIToFIPmtStsRpt{
				GrpHdr: models.GroupHeader{
					MsgId:   msgId,
					CreDtTm: time.Now().Format(time.RFC3339),
					NbOfTxs: "1",
				},
				OrgnlGrpInfAndSts: models.OriginalGroupInformationAndStatus{
					OrgnlMsgId:   reqPayload.OriginalMsgId,
					OrgnlMsgNmId: reqPayload.OriginalMsgNameId,
					OrgnlCreDtTm: reqPayload.OriginalCreDtTm,
					GrpSts:       reqPayload.GroupStatus,
				},
				TxInfAndSts: []models.PaymentTransactionInformationAndStatus{
					{
						OrgnlInstrId:    reqPayload.OriginalInstrId,
						OrgnlEndToEndId: reqPayload.OriginalEndToEndId,
						OrgnlTxId:       reqPayload.OriginalTxId,
						OrgnlUETR:       reqPayload.OriginalUETR,
						TxSts:           reqPayload.TxStatus,
						StsRsnInf:       stsRsnInf,
					},
				},
			},
		}

		xmlBytes, err := utils.MarshalMX(doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build MX document: " + err.Error()})
			return
		}

		fileName := msgId + ".xml"
		pool.Submit(worker.Job{
			MessageID: msgId,
			FileName:  fileName,
			Content:   xmlBytes,
			MsgType:   "pacs002",
		})

		c.JSON(http.StatusAccepted, gin.H{
			"message_id":   msgId,
			"message_type": "pacs.002.001.10",
			"status":       worker.StatusPending,
			"file_name":    fileName,
			"inquiry_url":  "/api/v1/pacs002/inquiry/" + msgId,
		})
	}
}
