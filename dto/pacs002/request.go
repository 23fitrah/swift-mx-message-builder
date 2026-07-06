package pacs002

import "swift-mx-message-builder/dto/general"

type Pacs002Request struct {
	general.GeneralRequest
	Payload FieldPacs002 `json:"payload" binding:"required"`
}

type FieldPacs002 struct {
	OriginalMsgId      string `json:"original_msg_id" binding:"required"`
	OriginalMsgNameId  string `json:"original_msg_name_id" binding:"required"`
	OriginalCreDtTm    string `json:"original_cre_dt_tm"`
	OriginalEndToEndId string `json:"original_end_to_end_id"`
	OriginalInstrId    string `json:"original_instr_id"`
	OriginalTxId       string `json:"original_tx_id"`
	OriginalUETR       string `json:"original_uetr"`
	GroupStatus        string `json:"group_status"`
	TxStatus           string `json:"tx_status" binding:"required"`
	ReasonCode         string `json:"reason_code"`
	AdditionalInfo     string `json:"additional_info"`
}
