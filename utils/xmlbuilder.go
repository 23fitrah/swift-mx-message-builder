package utils

import "encoding/xml"

func MarshalMX(doc interface{}) ([]byte, error) {
	body, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	out := append([]byte(xml.Header), body...)
	return out, nil
}
