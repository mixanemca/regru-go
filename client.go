/*
Copyright © 2025 Michael Bruskov <mixanemca@yandex.ru>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package regru

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for reg.ru API.
	DefaultBaseURL = "https://api.reg.ru/api/regru2"
	// DefaultTimeout is the default timeout for HTTP requests.
	DefaultTimeout = 30 * time.Second
)

// Client represents a client for working with reg.ru API.
type Client struct {
	username   string
	password   string
	baseURL    string
	httpClient *http.Client
}

// ClientOption represents an option for configuring the client.
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the timeout for HTTP requests.
// If a custom HTTP client is set via WithHTTPClient, this option will update its timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new instance of reg.ru client.
func NewClient(username, password string, opts ...ClientOption) *Client {
	client := &Client{
		username: username,
		password: password,
		baseURL:  DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// apiRequest performs a request to reg.ru API.
func (c *Client) apiRequest(ctx context.Context, path string, apiReq APIRequest) ([]byte, error) {
	// Set credentials in the request
	apiReq.SetCredentials(c.username, c.password)

	// Build URL
	apiURL := fmt.Sprintf("%s/%s", c.baseURL, path)

	// Serialize parameters to JSON for sending
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request params: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: string(bodyBytes)}
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors in response
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil {
		if apiResp.ErrorText != "" {
			return nil, &APIError{Message: apiResp.ErrorText}
		}
	}

	return body, nil
}

// getAddRecordPath returns the API path for adding a record of the specified type.
func getAddRecordPath(recordType string) (string, error) {
	switch recordType {
	case "A":
		return "zone/add_alias", nil
	case "AAAA":
		return "zone/add_aaaa", nil
	case "CNAME":
		return "zone/add_cname", nil
	case "MX":
		return "zone/add_mx", nil
	case "NS":
		return "zone/add_ns", nil
	case "TXT":
		return "zone/add_txt", nil
	default:
		return "", &UnsupportedRecordTypeError{RecordType: recordType}
	}
}

// getRemoveRecordPath returns the API path for removing a record of the specified type.
func getRemoveRecordPath(recordType string) (string, error) {
	switch recordType {
	case "A":
		return "zone/remove_alias", nil
	case "AAAA":
		return "zone/remove_aaaa", nil
	case "CNAME":
		return "zone/remove_cname", nil
	case "MX":
		return "zone/remove_mx", nil
	case "NS":
		return "zone/remove_ns", nil
	case "TXT":
		return "zone/remove_txt", nil
	default:
		return "", &UnsupportedRecordTypeError{RecordType: recordType}
	}
}

// createAddRecordRequest creates an appropriate request structure based on record type.
func createAddRecordRequest(zone string, params CreateDNSRecordParams) (APIRequest, error) {
	domain := AddRecordDomain{
		DName:     zone,
		Subdomain: params.Name,
		Content:   params.Content,
	}

	baseReq := AddRecordRequest{
		BaseRequest: BaseRequest{},
		Domains:     []AddRecordDomain{domain},
	}

	if params.TTL > 0 {
		baseReq.TTL = params.TTL
	}

	switch params.Type {
	case "A":
		return &AddAliasRequest{AddRecordRequest: baseReq}, nil
	case "AAAA":
		return &AddAAAARequest{AddRecordRequest: baseReq}, nil
	case "CNAME":
		return &AddCNAMERequest{AddRecordRequest: baseReq}, nil
	case "MX":
		// MX records may have priority in content (e.g., "10 mail.example.com")
		mxReq := &AddMXRequest{AddRecordRequest: baseReq}
		// Try to extract priority from content if format is "priority hostname"
		// This is a simple implementation - may need adjustment based on actual API
		return mxReq, nil
	case "NS":
		return &AddNSRequest{AddRecordRequest: baseReq}, nil
	case "TXT":
		return &AddTXTRequest{AddRecordRequest: baseReq}, nil
	default:
		return nil, &UnsupportedRecordTypeError{RecordType: params.Type}
	}
}

// createRemoveRecordRequest creates an appropriate request structure based on record type.
func createRemoveRecordRequest(zone string, rr DNSRecord) (APIRequest, error) {
	domain := RemoveRecordDomain{
		DName:     zone,
		Subdomain: rr.Name,
		Content:   rr.Content,
	}

	baseReq := RemoveRecordRequest{
		BaseRequest: BaseRequest{},
		Domains:     []RemoveRecordDomain{domain},
	}

	switch rr.Type {
	case "A":
		return &RemoveAliasRequest{RemoveRecordRequest: baseReq}, nil
	case "AAAA":
		return &RemoveAAAARequest{RemoveRecordRequest: baseReq}, nil
	case "CNAME":
		return &RemoveCNAMERequest{RemoveRecordRequest: baseReq}, nil
	case "MX":
		return &RemoveMXRequest{RemoveRecordRequest: baseReq}, nil
	case "NS":
		return &RemoveNSRequest{RemoveRecordRequest: baseReq}, nil
	case "TXT":
		return &RemoveTXTRequest{RemoveRecordRequest: baseReq}, nil
	default:
		return nil, &UnsupportedRecordTypeError{RecordType: rr.Type}
	}
}

// AddRR creates a new DNS record for the specified zone.
func (c *Client) AddRR(ctx context.Context, zone string, params CreateDNSRecordParams) (DNSRecord, error) {
	// Get the appropriate API path for this record type
	path, err := getAddRecordPath(params.Type)
	if err != nil {
		return DNSRecord{}, err
	}

	// Create the appropriate request structure
	apiReq, err := createAddRecordRequest(zone, params)
	if err != nil {
		return DNSRecord{}, err
	}

	// Execute API request
	body, err := c.apiRequest(ctx, path, apiReq)
	if err != nil {
		return DNSRecord{}, err
	}

	// Parse response
	var resp AddNSResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return DNSRecord{}, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert response to DNSRecord
	record := DNSRecord{
		Name:    params.Name,
		Type:    params.Type,
		Content: params.Content,
		TTL:     params.TTL,
	}

	// Extract record ID from response if available
	if len(resp.Answer.Domains) > 0 {
		domain := resp.Answer.Domains[0]
		if domain.Result == "success" {
			record.ID = domain.DNSID
		}
	}

	return record, nil
}

// DeleteRR deletes a DNS record from the specified zone.
func (c *Client) DeleteRR(ctx context.Context, zone string, rr DNSRecord) error {
	// Get the appropriate API path for this record type
	path, err := getRemoveRecordPath(rr.Type)
	if err != nil {
		return err
	}

	// Create the appropriate request structure
	apiReq, err := createRemoveRecordRequest(zone, rr)
	if err != nil {
		return err
	}

	// Execute API request
	_, err = c.apiRequest(ctx, path, apiReq)
	if err != nil {
		return err
	}

	return nil
}

// GetRRByName returns a DNS record by name in the specified zone.
func (c *Client) GetRRByName(ctx context.Context, zone, name string) (DNSRecord, error) {
	// Get all zone records
	params := ListDNSRecordsParams{
		ZoneName: zone,
	}

	records, err := c.ListRecords(ctx, params)
	if err != nil {
		return DNSRecord{}, err
	}

	// Search for record by name
	for _, record := range records {
		if record.Name == name {
			return record, nil
		}
	}

	return DNSRecord{}, &RecordNotFoundError{RecordName: name}
}

// ListZones returns a list of all zones in the account.
func (c *Client) ListZones(ctx context.Context) ([]Zone, error) {
	// Prepare API request
	apiReq := ServiceListRequest{
		BaseRequest: BaseRequest{},
		PageSize:    1000, // Maximum number of zones per request
	}

	body, err := c.apiRequest(ctx, "service/get_list", &apiReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ServiceListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var zones []Zone
	for _, service := range resp.Answer.Services {
		if service.ServiceType == "domain" {
			zones = append(zones, Zone{
				Name: service.Domain,
				ID:   service.ServiceID,
			})
		}
	}

	return zones, nil
}

// ListZonesByName returns a list of zones by name.
func (c *Client) ListZonesByName(ctx context.Context, name string) ([]Zone, error) {
	zones, err := c.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Zone
	for _, zone := range zones {
		if zone.Name == name {
			filtered = append(filtered, zone)
		}
	}

	return filtered, nil
}

// ListRecords returns a list of DNS records for the specified zone.
func (c *Client) ListRecords(ctx context.Context, params ListDNSRecordsParams) ([]DNSRecord, error) {
	zoneName := params.ZoneName
	if zoneName == "" {
		zoneName = params.ZoneID // Fallback to ZoneID if ZoneName is not set
	}

	// Prepare API request
	apiReq := ZoneGetNSRequest{
		BaseRequest: BaseRequest{},
		Domains:     []string{zoneName},
	}

	body, err := c.apiRequest(ctx, "zone/get_ns", &apiReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ZoneListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var records []DNSRecord
	for _, domain := range resp.Answer.Domains {
		if domain.DName == zoneName {
			for _, nsRecord := range domain.NSList {
				record := DNSRecord{
					Name:    nsRecord.Subdomain,
					Type:    nsRecord.Type,
					Content: nsRecord.Content,
					TTL:     nsRecord.TTL,
					ID:      nsRecord.DNSID,
				}

				// Apply filters if specified
				if params.Name != "" && record.Name != params.Name {
					continue
				}
				if params.Type != "" && record.Type != params.Type {
					continue
				}

				records = append(records, record)
			}
		}
	}

	return records, nil
}

// ListRecordsByZoneID returns a list of DNS records by zone identifier.
func (c *Client) ListRecordsByZoneID(ctx context.Context, id string, params ListDNSRecordsParams) ([]DNSRecord, error) {
	// In reg.ru API, zone identifier usually matches zone name
	// Get zone by ID and use its name
	zones, err := c.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	var zoneName string
	for _, zone := range zones {
		if zone.ID == id {
			zoneName = zone.Name
			break
		}
	}

	if zoneName == "" {
		return nil, &ZoneNotFoundError{ZoneID: id}
	}

	params.ZoneName = zoneName
	return c.ListRecords(ctx, params)
}

// UpdateRR updates an existing DNS record in the specified zone.
func (c *Client) UpdateRR(ctx context.Context, zone string, rr DNSRecord) (DNSRecord, error) {
	// In reg.ru API, record update is usually performed through delete and create
	// First, delete the old record
	if err := c.DeleteRR(ctx, zone, rr); err != nil {
		return DNSRecord{}, err
	}

	// Create a new record with updated data
	createParams := CreateDNSRecordParams{
		Name:    rr.Name,
		Type:    rr.Type,
		Content: rr.Content,
		TTL:     rr.TTL,
	}

	return c.AddRR(ctx, zone, createParams)
}
