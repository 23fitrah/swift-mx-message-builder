package models

import "encoding/xml"

// Pacs002Document is the root Document element for pacs.002.001.10
// (FIToFIPaymentStatusReport) - used to report the status (accepted,
// rejected, pending, etc.) of a previously sent payment instruction,
// and also used as the response body for status inquiries.
type Pacs002Document struct {
	XMLName         xml.Name        `xml:"urn:iso:std:iso:20022:tech:xsd:pacs.002.001.10 Document"`
	FIToFIPmtStsRpt FIToFIPmtStsRpt `xml:"FIToFIPmtStsRpt"`
}

type FIToFIPmtStsRpt struct {
	GrpHdr            GroupHeader                              `xml:"GrpHdr"`
	OrgnlGrpInfAndSts OriginalGroupInformationAndStatus        `xml:"OrgnlGrpInfAndSts"`
	TxInfAndSts       []PaymentTransactionInformationAndStatus `xml:"TxInfAndSts,omitempty"`
}

type OriginalGroupInformationAndStatus struct {
	OrgnlMsgId   string `xml:"OrgnlMsgId"`
	OrgnlMsgNmId string `xml:"OrgnlMsgNmId"`
	OrgnlCreDtTm string `xml:"OrgnlCreDtTm,omitempty"`
	GrpSts       string `xml:"GrpSts,omitempty"`
}

type PaymentTransactionInformationAndStatus struct {
	OrgnlInstrId    string                   `xml:"OrgnlInstrId,omitempty"`
	OrgnlEndToEndId string                   `xml:"OrgnlEndToEndId,omitempty"`
	OrgnlTxId       string                   `xml:"OrgnlTxId,omitempty"`
	OrgnlUETR       string                   `xml:"OrgnlUETR,omitempty"`
	TxSts           string                   `xml:"TxSts,omitempty"`
	StsRsnInf       *StatusReasonInformation `xml:"StsRsnInf,omitempty"`
}

type StatusReasonInformation struct {
	Rsn      string `xml:"Rsn>Cd,omitempty"`
	AddtlInf string `xml:"AddtlInf,omitempty"`
}
