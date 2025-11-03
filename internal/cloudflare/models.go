package cloudflare

import (
	"errors"
	"time"
)

var ErrRecordExists = errors.New("dns_record_exists")

type DNSListResponse struct {
	Result []DNSRecord `json:"result"`
}

type DNSRecord struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	TTL        int       `json:"ttl"`
	Settings   struct{}  `json:"settings"`
	Meta       struct{}  `json:"meta"`
	Comment    *string   `json:"comment"`
	Tags       []string  `json:"tags"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
}

type CloudflareResponse struct {
	Success bool      `json:"success"`
	Errors  []CFError `json:"errors"`
	Result  DNSRecord `json:"result"`
}

type CFError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
