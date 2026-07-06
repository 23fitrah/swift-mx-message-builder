package models

import "encoding/xml"

type Pacs008Document struct {
	XMLName           xml.Name          `xml:"urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08 Document"`
	FIToFICstmrCdtTrf FIToFICstmrCdtTrf `xml:"FIToFICstmrCdtTrf"`
}

type FIToFICstmrCdtTrf struct {
	GrpHdr      GroupHeader                            `xml:"GrpHdr"`
	CdtTrfTxInf []CreditTransferTransactionInformation `xml:"CdtTrfTxInf"`
}

type CreditTransferTransactionInformation struct {
	PmtId          PaymentIdentification                       `xml:"PmtId"`
	IntrBkSttlmAmt ActiveCurrencyAndAmount                     `xml:"IntrBkSttlmAmt"`
	Dbtr           PartyIdentification                         `xml:"Dbtr"`
	DbtrAgt        BranchAndFinancialInstitutionIdentification `xml:"DbtrAgt"`
	CdtrAgt        BranchAndFinancialInstitutionIdentification `xml:"CdtrAgt"`
	Cdtr           PartyIdentification                         `xml:"Cdtr"`
	RmtInf         *RemittanceInformation                      `xml:"RmtInf,omitempty"`
}
