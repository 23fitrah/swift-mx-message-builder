package handlers

import (
	"net/http"
	"time"

	"swift-mx-message-builder/dto/pacs009"
	"swift-mx-message-builder/models"
	"swift-mx-message-builder/utils"
	"swift-mx-message-builder/worker"

	"github.com/gin-gonic/gin"
)

func Pacs009Handler(pool *worker.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pacs009.Pacs009Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var reqPayload = req.Payload
		msgId := utils.GenerateMessageID("PACS009")
		txId := utils.GenerateMessageID("TX")

		doc := models.Pacs009Document{
			FICdtTrf: models.FICdtTrf{
				GrpHdr: models.GroupHeader{
					MsgId:   msgId,
					CreDtTm: time.Now().Format(time.RFC3339),
					NbOfTxs: "1",
					SttlmInf: &models.SettlementInformation{
						SttlmMtd: utils.DefaultOr(reqPayload.SettlementMethod, "COVE"),
					},
				},
				CdtTrfTxInf: []models.FinancialInstitutionCreditTransferTransaction{
					{
						PmtId: models.PaymentIdentification{
							InstrId:    reqPayload.InstrId,
							EndToEndId: reqPayload.EndToEndId,
							TxId:       txId,
							UETR:       utils.GenerateUETR(),
						},
						IntrBkSttlmAmt: models.ActiveCurrencyAndAmount{Ccy: reqPayload.Currency, Value: reqPayload.Amount},
						Dbtr:           models.BranchAndFinancialInstitutionIdentification{FinInstnId: models.FinInstnId{BICFI: reqPayload.DebtorBIC}},
						DbtrAgt:        models.BranchAndFinancialInstitutionIdentification{FinInstnId: models.FinInstnId{BICFI: reqPayload.DebtorAgentBIC}},
						CdtrAgt:        models.BranchAndFinancialInstitutionIdentification{FinInstnId: models.FinInstnId{BICFI: reqPayload.CreditorAgentBIC}},
						Cdtr:           models.BranchAndFinancialInstitutionIdentification{FinInstnId: models.FinInstnId{BICFI: reqPayload.CreditorBIC}},
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
			MsgType:   "pacs009",
		})

		c.JSON(http.StatusAccepted, gin.H{
			"message_id":     msgId,
			"transaction_id": txId,
			"message_type":   "pacs.009.001.08",
			"status":         worker.StatusPending,
			"file_name":      fileName,
			"inquiry_url":    "/api/v1/pacs009/inquiry/" + msgId,
		})
	}
}
