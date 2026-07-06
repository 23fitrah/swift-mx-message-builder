package pacs004

import "swift-mx-message-builder/dto/general"

type Pacs004Request struct {
	general.GeneralRequest
	Payload FieldPacs004 `json:"payload" binding:"required"`
}

type FieldPacs004 struct {
	OriginalInstrId    string  `json:"original_instr_id"`
	OriginalEndToEndId string  `json:"original_end_to_end_id"`
	OriginalTxId       string  `json:"original_tx_id"`
	OriginalUETR       string  `json:"original_uetr"`
	ReturnedAmount     float64 `json:"returned_amount" binding:"required"`
	Currency           string  `json:"currency" binding:"required,len=3"`
	ReasonCode         string  `json:"reason_code" binding:"required"`
	AdditionalInfo     string  `json:"additional_info"`
}
