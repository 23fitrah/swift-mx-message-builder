package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"swift-mx-message-builder/dto/pacs004"
	"swift-mx-message-builder/models"
	"swift-mx-message-builder/utils"
	"swift-mx-message-builder/worker"
)

func Pacs004Handler(pool *worker.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pacs004.Pacs004Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var reqPayload = req.Payload

		msgId := utils.GenerateMessageID("PACS004")
		rtrId := utils.GenerateMessageID("RTR")

		doc := models.Pacs004Document{
			PmtRtr: models.PmtRtr{
				GrpHdr: models.GroupHeader{
					MsgId:   msgId,
					CreDtTm: time.Now().Format(time.RFC3339),
					NbOfTxs: "1",
				},
				TxInf: []models.PaymentTransactionInformation{
					{
						RtrId:              rtrId,
						OrgnlInstrId:       reqPayload.OriginalInstrId,
						OrgnlEndToEndId:    reqPayload.OriginalInstrId,
						OrgnlTxId:          reqPayload.OriginalTxId,
						OrgnlUETR:          reqPayload.OriginalUETR,
						RtrdIntrBkSttlmAmt: models.ActiveCurrencyAndAmount{Ccy: reqPayload.Currency, Value: reqPayload.ReturnedAmount},
						RtrRsnInf: &models.ReturnReasonInformation{
							Rsn:      reqPayload.ReasonCode,
							AddtlInf: reqPayload.AdditionalInfo,
						},
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
			MsgType:   "pacs004",
		})

		c.JSON(http.StatusAccepted, gin.H{
			"message_id":   msgId,
			"return_id":    rtrId,
			"message_type": "pacs.004.001.09",
			"status":       worker.StatusPending,
			"file_name":    fileName,
			"inquiry_url":  "/api/v1/pacs004/inquiry/" + msgId,
		})
	}
}
