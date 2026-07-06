package pacs009

import "swift-mx-message-builder/dto/general"

// Pacs009Request is the JSON payload accepted by POST /pacs009/generate.
type Pacs009Request struct {
	general.GeneralRequest
	Payload FieldPacs009 `json:"payload" binding:"required"`
}

type FieldPacs009 struct {
	InstrId          string  `json:"instr_id"`
	EndToEndId       string  `json:"end_to_end_id" binding:"required"`
	Amount           float64 `json:"amount" binding:"required"`
	Currency         string  `json:"currency" binding:"required,len=3"`
	DebtorBIC        string  `json:"debtor_bic" binding:"required"`
	DebtorAgentBIC   string  `json:"debtor_agent_bic" binding:"required"`
	CreditorAgentBIC string  `json:"creditor_agent_bic" binding:"required"`
	CreditorBIC      string  `json:"creditor_bic" binding:"required"`
	SettlementMethod string  `json:"settlement_method"`
}
