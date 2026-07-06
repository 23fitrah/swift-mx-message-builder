package pacs009

import "swift-mx-message-builder/dto/general"

type Pacs009Request struct {
	general.GeneralRequest
	Payload FieldPacs009 `json:"payload" validate:"required"`
}

type FieldPacs009 struct {
	InstrId          string  `json:"instr_id"`
	EndToEndId       string  `json:"end_to_end_id" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	Currency         string  `json:"currency" validate:"required,len=3"`
	DebtorBIC        string  `json:"debtor_bic" validate:"required"`
	DebtorAgentBIC   string  `json:"debtor_agent_bic" validate:"required"`
	CreditorAgentBIC string  `json:"creditor_agent_bic" validate:"required"`
	CreditorBIC      string  `json:"creditor_bic" validate:"required"`
	SettlementMethod string  `json:"settlement_method"`
}
