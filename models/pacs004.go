package models

import "encoding/xml"

type Pacs004Document struct {
	XMLName xml.Name `xml:"urn:iso:std:iso:20022:tech:xsd:pacs.004.001.09 Document"`
	PmtRtr  PmtRtr   `xml:"PmtRtr"`
}

type PmtRtr struct {
	GrpHdr GroupHeader                     `xml:"GrpHdr"`
	TxInf  []PaymentTransactionInformation `xml:"TxInf"`
}

type PaymentTransactionInformation struct {
	RtrId              string                   `xml:"RtrId"`
	OrgnlInstrId       string                   `xml:"OrgnlInstrId,omitempty"`
	OrgnlEndToEndId    string                   `xml:"OrgnlEndToEndId,omitempty"`
	OrgnlTxId          string                   `xml:"OrgnlTxId,omitempty"`
	OrgnlUETR          string                   `xml:"OrgnlUETR,omitempty"`
	RtrdIntrBkSttlmAmt ActiveCurrencyAndAmount  `xml:"RtrdIntrBkSttlmAmt"`
	RtrRsnInf          *ReturnReasonInformation `xml:"RtrRsnInf,omitempty"`
}

type ReturnReasonInformation struct {
	Rsn      string `xml:"Rsn>Cd,omitempty"`
	AddtlInf string `xml:"AddtlInf,omitempty"`
}
