package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

func (c *CfAPI) GetAllSubdomains(domain, zoneID string) ([]string, error) {
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records?per_page=1000", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var parsedResp DNSListResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return nil, fmt.Errorf("failed to parse server response: %w", err)
	}

	zone := domain
	subdomainsSet := make(map[string]struct{})

	for _, record := range parsedResp.Result {
		if strings.HasSuffix(record.Name, "."+zone) {
			sub := strings.TrimSuffix(record.Name, "."+zone)
			if sub == "" {
				continue
			}
			parts := strings.Split(sub, ".")
			if len(parts) > 0 && parts[0] != "" {
				subdomainsSet[parts[0]] = struct{}{}
			}
		}
	}

	subdomains := make([]string, 0, len(subdomainsSet))
	for sub := range subdomainsSet {
		subdomains = append(subdomains, sub)
	}

	sort.Strings(subdomains)
	return subdomains, nil
}

func (c *CfAPI) CreateDNSRecord(zoneID string, data any) error {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var cfResp CloudflareResponse
	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if cfResp.Success {
		return nil
	}

	for _, errObj := range cfResp.Errors {
		if errObj.Code == 81058 {
			return ErrRecordExists
		}
	}

	return fmt.Errorf("failed to create DNS record: %v", cfResp.Errors)
}

func (c *CfAPI) DeleteDNSRecord(zoneID, name string) error {
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records?name="+name, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read GET response body: %w", err)
	}

	var listResp struct {
		Success bool        `json:"success"`
		Errors  []CFError   `json:"errors"`
		Result  []DNSRecord `json:"result"`
	}

	if err := json.Unmarshal(respBody, &listResp); err != nil {
		return fmt.Errorf("failed to unmarshal GET response: %w", err)
	}

	if !listResp.Success || len(listResp.Result) == 0 {
		return fmt.Errorf("record not found or API error: %v", listResp.Errors)
	}

	recordID := listResp.Result[0].ID

	delURL := "https://api.cloudflare.com/client/v4/zones/" + zoneID + "/dns_records/" + recordID

	delReq, err := http.NewRequest("DELETE", delURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	delResp, err := c.client.Do(delReq)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %w", err)
	}
	defer delResp.Body.Close()

	delRespBody, err := io.ReadAll(delResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read DELETE response body: %w", err)
	}

	var cfResp CloudflareResponse
	if err := json.Unmarshal(delRespBody, &cfResp); err != nil {
		return fmt.Errorf("failed to unmarshal DELETE response: %w", err)
	}

	if !cfResp.Success {
		return fmt.Errorf("failed to delete DNS record: %v", cfResp.Errors)
	}

	return nil
}
