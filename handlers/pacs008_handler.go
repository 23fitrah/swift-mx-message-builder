package handlers

import (
	"net/http"
	"time"

	"swift-mx-message-builder/dto/pacs008"
	"swift-mx-message-builder/models"
	"swift-mx-message-builder/utils"
	"swift-mx-message-builder/worker"

	"github.com/gin-gonic/gin"
)

// Pacs008Handler builds a pacs.008.001.08 MX document from the request,
// hands the marshalled XML off to the worker pool for asynchronous file
// writing, and immediately returns the generated message ID + PENDING
// status so the caller can poll the inquiry endpoint.
func Pacs008Handler(pool *worker.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pacs008.Pacs008Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var reqPayload = req.Payload
		msgId := utils.GenerateMessageID("PACS008")
		txId := utils.GenerateMessageID("TX")

		doc := models.Pacs008Document{
			FIToFICstmrCdtTrf: models.FIToFICstmrCdtTrf{
				GrpHdr: models.GroupHeader{
					MsgId:   msgId,
					CreDtTm: time.Now().Format(time.RFC3339),
					NbOfTxs: "1",
					SttlmInf: &models.SettlementInformation{
						SttlmMtd: utils.DefaultOr(reqPayload.SettlementMethod, "CLRG"),
					},
				},
				CdtTrfTxInf: []models.CreditTransferTransactionInformation{
					{
						PmtId: models.PaymentIdentification{
							InstrId:    reqPayload.InstrId,
							EndToEndId: reqPayload.EndToEndId,
							TxId:       txId,
							UETR:       utils.GenerateUETR(),
						},
						IntrBkSttlmAmt: models.ActiveCurrencyAndAmount{Ccy: reqPayload.Currency, Value: reqPayload.Amount},
						Dbtr:           models.PartyIdentification{Nm: reqPayload.DebtorName},
						DbtrAgt:        models.BranchAndFinancialInstitutionIdentification{FinInstnId: models.FinInstnId{BICFI: reqPayload.DebtorAgentBIC}},
						CdtrAgt:        models.BranchAndFinancialInstitutionIdentification{FinInstnId: models.FinInstnId{BICFI: reqPayload.CreditorAgentBIC}},
						Cdtr:           models.PartyIdentification{Nm: reqPayload.CreditorName},
						RmtInf:         &models.RemittanceInformation{Ustrd: reqPayload.RemittanceInfo},
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
			MsgType:   "pacs008",
		})

		c.JSON(http.StatusAccepted, gin.H{
			"message_id":     msgId,
			"transaction_id": txId,
			"message_type":   "pacs.008.001.08",
			"status":         worker.StatusPending,
			"file_name":      fileName,
			"inquiry_url":    "/api/v1/pacs008/inquiry/" + msgId,
		})
	}
}
