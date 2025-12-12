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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

	// Serialize request structure to JSON for input_data parameter
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request params: %w", err)
	}

	// Create form data with required parameters
	formData := url.Values{}
	formData.Set("input_format", "json")
	formData.Set("input_data", string(jsonData))
	formData.Set("username", c.username)
	formData.Set("password", c.password)

	// Create HTTP request with form data
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	case RecordTypeA:
		return "zone/add_alias", nil
	case RecordTypeAAAA:
		return "zone/add_aaaa", nil
	case RecordTypeCNAME:
		return "zone/add_cname", nil
	case RecordTypeMX:
		return "zone/add_mx", nil
	case RecordTypeNS:
		return "zone/add_ns", nil
	case RecordTypeTXT:
		return "zone/add_txt", nil
	default:
		return "", &UnsupportedRecordTypeError{RecordType: recordType}
	}
}

// getRemoveRecordPath returns the API path for removing a record of the specified type.
// According to reg.ru API documentation, all record types use the same endpoint: zone/remove_record
func getRemoveRecordPath(recordType string) (string, error) {
	switch recordType {
	case RecordTypeA, RecordTypeAAAA, RecordTypeCNAME, RecordTypeMX, RecordTypeNS, RecordTypeTXT:
		return "zone/remove_record", nil
	default:
		return "", &UnsupportedRecordTypeError{RecordType: recordType}
	}
}

// createAddRecordRequest creates an appropriate request structure based on record type.
func createAddRecordRequest(zone string, params CreateDNSRecordParams) (APIRequest, error) {
	switch params.Type {
	case RecordTypeA:
		// For A records (add_alias), ipaddr and subdomain are at request level
		aliasReq := &AddAliasRequest{
			BaseRequest: BaseRequest{},
			Domains: []AddAliasDomain{
				{DName: zone},
			},
			Subdomain: params.Name,
			IPAddr:    params.Content,
		}
		if params.TTL > 0 {
			aliasReq.TTL = params.TTL
		}
		return aliasReq, nil
	case RecordTypeAAAA:
		// For AAAA records (add_aaaa), ipaddr and subdomain are at request level
		aaaaReq := &AddAAAARequest{
			BaseRequest: BaseRequest{},
			Domains: []AddAliasDomain{
				{DName: zone},
			},
			Subdomain: params.Name,
			IPAddr:    params.Content,
		}
		if params.TTL > 0 {
			aaaaReq.TTL = params.TTL
		}
		return aaaaReq, nil
	case RecordTypeCNAME:
		// For CNAME records (add_cname), canonical_name and subdomain are at request level
		cnameReq := &AddCNAMERequest{
			BaseRequest: BaseRequest{},
			Domains: []AddAliasDomain{
				{DName: zone},
			},
			Subdomain:     params.Name,
			CanonicalName: params.Content,
		}
		if params.TTL > 0 {
			cnameReq.TTL = params.TTL
		}
		return cnameReq, nil
	case RecordTypeMX:
		// For MX records (add_mx), mail_server and subdomain are at request level
		mxReq := &AddMXRequest{
			BaseRequest: BaseRequest{},
			Domains: []AddAliasDomain{
				{DName: zone},
			},
			Subdomain:  params.Name,
			MailServer: params.Content,
		}
		if params.TTL > 0 {
			mxReq.TTL = params.TTL
		}
		return mxReq, nil
	case RecordTypeNS:
		// For NS records (add_ns), dns_server and subdomain are at request level
		nsReq := &AddNSRequest{
			BaseRequest: BaseRequest{},
			Domains: []AddAliasDomain{
				{DName: zone},
			},
			Subdomain: params.Name,
			DNSServer: params.Content,
		}
		if params.TTL > 0 {
			nsReq.TTL = params.TTL
		}
		return nsReq, nil
	case RecordTypeTXT:
		// For TXT records (add_txt), text and subdomain are at request level
		txtReq := &AddTXTRequest{
			BaseRequest: BaseRequest{},
			Domains: []AddAliasDomain{
				{DName: zone},
			},
			Subdomain: params.Name,
			Text:      params.Content,
		}
		if params.TTL > 0 {
			txtReq.TTL = params.TTL
		}
		return txtReq, nil
	default:
		return nil, &UnsupportedRecordTypeError{RecordType: params.Type}
	}
}

// createRemoveRecordRequest creates an appropriate request structure based on record type.
// According to reg.ru API documentation, remove_record uses subdomain, content, and record_type at request level.
func createRemoveRecordRequest(zone string, rr DNSRecord) (APIRequest, error) {
	// All record types use the same structure for removal
	req := &RemoveRecordRequest{
		BaseRequest: BaseRequest{},
		Domains: []RemoveRecordDomain{
			{DName: zone},
		},
		Subdomain:  rr.Name,
		Content:    rr.Content,
		RecordType: rr.Type,
	}

	// All remove requests use the same structure, but we return typed requests for consistency
	switch rr.Type {
	case RecordTypeA:
		return &RemoveAliasRequest{RemoveRecordRequest: *req}, nil
	case RecordTypeAAAA:
		return &RemoveAAAARequest{RemoveRecordRequest: *req}, nil
	case RecordTypeCNAME:
		return &RemoveCNAMERequest{RemoveRecordRequest: *req}, nil
	case RecordTypeMX:
		return &RemoveMXRequest{RemoveRecordRequest: *req}, nil
	case RecordTypeNS:
		return &RemoveNSRequest{RemoveRecordRequest: *req}, nil
	case RecordTypeTXT:
		return &RemoveTXTRequest{RemoveRecordRequest: *req}, nil
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
		serviceType := service.GetServiceType()
		if serviceType == "domain" {
			zones = append(zones, Zone{
				Name: service.GetDomain(),
				ID:   service.GetServiceID(),
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
	apiReq := ZoneGetResourceRecordsRequest{
		BaseRequest: BaseRequest{},
		Domains: []ZoneGetResourceRecordsDomain{
			{
				DName: zoneName,
			},
		},
	}

	body, err := c.apiRequest(ctx, "zone/get_resource_records", &apiReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ZoneGetResourceRecordsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var records []DNSRecord
	for _, domain := range resp.Answer.Domains {
		if domain.DName == zoneName {
			for _, rr := range domain.RRList {
				record := DNSRecord{
					Name:    rr.Subname,
					Type:    rr.Rectype,
					Content: rr.Content,
					// TTL and ID are not available in get_resource_records response
					// TTL:     rr.TTL,
					// ID:      rr.DNSID,
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
