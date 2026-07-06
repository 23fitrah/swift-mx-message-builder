package pacs008

import "swift-mx-message-builder/dto/general"

type Pacs008Request struct {
	general.GeneralRequest
	Payload FieldPacs008 `json:"payload" validate:"required"`
}

type FieldPacs008 struct {
	InstrId          string  `json:"instr_id"`
	EndToEndId       string  `json:"end_to_end_id" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	Currency         string  `json:"currency" validate:"required,len=3"`
	DebtorName       string  `json:"debtor_name" validate:"required"`
	DebtorAgentBIC   string  `json:"debtor_agent_bic" validate:"required"`
	CreditorName     string  `json:"creditor_name" validate:"required"`
	CreditorAgentBIC string  `json:"creditor_agent_bic" validate:"required"`
	RemittanceInfo   string  `json:"remittance_info"`
	SettlementMethod string  `json:"settlement_method"`
}
