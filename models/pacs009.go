package models

import "encoding/xml"

type Pacs009Document struct {
	XMLName  xml.Name `xml:"urn:iso:std:iso:20022:tech:xsd:pacs.009.001.08 Document"`
	FICdtTrf FICdtTrf `xml:"FICdtTrf"`
}

type FICdtTrf struct {
	GrpHdr      GroupHeader                                     `xml:"GrpHdr"`
	CdtTrfTxInf []FinancialInstitutionCreditTransferTransaction `xml:"CdtTrfTxInf"`
}

type FinancialInstitutionCreditTransferTransaction struct {
	PmtId          PaymentIdentification                       `xml:"PmtId"`
	IntrBkSttlmAmt ActiveCurrencyAndAmount                     `xml:"IntrBkSttlmAmt"`
	Dbtr           BranchAndFinancialInstitutionIdentification `xml:"Dbtr"`
	DbtrAgt        BranchAndFinancialInstitutionIdentification `xml:"DbtrAgt"`
	CdtrAgt        BranchAndFinancialInstitutionIdentification `xml:"CdtrAgt"`
	Cdtr           BranchAndFinancialInstitutionIdentification `xml:"Cdtr"`
}
