package utils

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// "PACS008-20260706153012-A1B2C3".
func GenerateMessageID(prefix string) string {
	return fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("20060102150405"), randHex(6))
}

// GenerateUETR generates a UUIDv4-like unique end-to-end transaction
func GenerateUETR() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func randHex(n int) string {
	b := make([]byte, n/2+1)
	_, _ = rand.Read(b)
	s := fmt.Sprintf("%x", b)
	return strings.ToUpper(s[:n])
}

// DefaultOr returns fallback if value is empty.
func DefaultOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
