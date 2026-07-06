package models

// FinInstnId represents a Financial Institution Identification block (BIC or Name).
type FinInstnId struct {
	BICFI string `xml:"BICFI,omitempty"`
	Nm    string `xml:"Nm,omitempty"`
}

// BranchAndFinancialInstitutionIdentification wraps a FinInstnId (used for Agents).
type BranchAndFinancialInstitutionIdentification struct {
	FinInstnId FinInstnId `xml:"FinInstnId"`
}

// PartyIdentification represents a simple Debtor/Creditor party (Name only, simplified).
type PartyIdentification struct {
	Nm string `xml:"Nm,omitempty"`
}

// ActiveCurrencyAndAmount represents an amount together with its ISO currency code.
type ActiveCurrencyAndAmount struct {
	Ccy   string  `xml:"Ccy,attr"`
	Value float64 `xml:",chardata"`
}

// SettlementInformation carries the interbank settlement method (e.g. CLRG, INDA, INGA, COVE).
type SettlementInformation struct {
	SttlmMtd string `xml:"SttlmMtd"`
}

// GroupHeader is the common Group Header (GrpHdr) block shared across pacs messages.
type GroupHeader struct {
	MsgId    string                 `xml:"MsgId"`
	CreDtTm  string                 `xml:"CreDtTm"`
	NbOfTxs  string                 `xml:"NbOfTxs"`
	SttlmInf *SettlementInformation `xml:"SttlmInf,omitempty"`
}

// PaymentIdentification holds the identifiers correlating a payment transaction.
type PaymentIdentification struct {
	InstrId    string `xml:"InstrId,omitempty"`
	EndToEndId string `xml:"EndToEndId"`
	TxId       string `xml:"TxId"`
	UETR       string `xml:"UETR,omitempty"`
}

// RemittanceInformation carries unstructured remittance info.
type RemittanceInformation struct {
	Ustrd string `xml:"Ustrd,omitempty"`
}
