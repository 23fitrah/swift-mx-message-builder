package pacs008

import "swift-mx-message-builder/dto/general"

// Pacs008Request is the JSON payload accepted by POST /pacs008/generate.
type Pacs008Request struct {
	general.GeneralRequest
	Payload FieldPacs008 `json:"payload" binding:"required"`
}

type FieldPacs008 struct {
	InstrId          string  `json:"instr_id"`
	EndToEndId       string  `json:"end_to_end_id" binding:"required"`
	Amount           float64 `json:"amount" binding:"required"`
	Currency         string  `json:"currency" binding:"required,len=3"`
	DebtorName       string  `json:"debtor_name" binding:"required"`
	DebtorAgentBIC   string  `json:"debtor_agent_bic" binding:"required"`
	CreditorName     string  `json:"creditor_name" binding:"required"`
	CreditorAgentBIC string  `json:"creditor_agent_bic" binding:"required"`
	RemittanceInfo   string  `json:"remittance_info"`
	SettlementMethod string  `json:"settlement_method"`
}
