package models

type FinInstnId struct {
	BICFI string `xml:"BICFI,omitempty"`
	Nm    string `xml:"Nm,omitempty"`
}

type BranchAndFinancialInstitutionIdentification struct {
	FinInstnId FinInstnId `xml:"FinInstnId"`
}

type PartyIdentification struct {
	Nm string `xml:"Nm,omitempty"`
}

type ActiveCurrencyAndAmount struct {
	Ccy   string  `xml:"Ccy,attr"`
	Value float64 `xml:",chardata"`
}

type SettlementInformation struct {
	SttlmMtd string `xml:"SttlmMtd"`
}

type GroupHeader struct {
	MsgId    string                 `xml:"MsgId"`
	CreDtTm  string                 `xml:"CreDtTm"`
	NbOfTxs  string                 `xml:"NbOfTxs"`
	SttlmInf *SettlementInformation `xml:"SttlmInf,omitempty"`
}

type PaymentIdentification struct {
	InstrId    string `xml:"InstrId,omitempty"`
	EndToEndId string `xml:"EndToEndId"`
	TxId       string `xml:"TxId"`
	UETR       string `xml:"UETR,omitempty"`
}

type RemittanceInformation struct {
	Ustrd string `xml:"Ustrd,omitempty"`
}
