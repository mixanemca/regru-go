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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates a test HTTP server that returns the provided response.
func setupTestServer(t *testing.T, responseBody interface{}, statusCode int) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if responseBody != nil {
			require.NoError(t, json.NewEncoder(w).Encode(responseBody))
		}
	})

	return httptest.NewServer(handler)
}

// setupTestClient creates a test client with a test server.
func setupTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()

	client := NewClient("test-username", "test-password",
		WithBaseURL(server.URL),
		WithTimeout(5*time.Second),
	)

	return client
}

func TestNewClient(t *testing.T) {
	client := NewClient("username", "password")
	require.NotNil(t, client)
	assert.Equal(t, DefaultBaseURL, client.baseURL)
	assert.Equal(t, DefaultTimeout, client.httpClient.Timeout)
}

func TestNewClient_WithOptions(t *testing.T) {
	customURL := "https://custom.api.url"
	customTimeout := 60 * time.Second

	client := NewClient("username", "password",
		WithBaseURL(customURL),
		WithTimeout(customTimeout),
	)

	assert.Equal(t, customURL, client.baseURL)
	assert.Equal(t, customTimeout, client.httpClient.Timeout)
}

func TestClient_AddRR(t *testing.T) {
	tests := []struct {
		name           string
		recordType     string
		content        string
		response       AddNSResponse
		expectedMethod string
		wantErr        bool
	}{
		{
			name:       "add A record",
			recordType: "A",
			content:    "192.0.2.1",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
							DNSID:  "12345",
						},
					},
				},
			},
			expectedMethod: "zone/add_alias",
			wantErr:        false,
		},
		{
			name:       "add AAAA record",
			recordType: "AAAA",
			content:    "2001:db8::1",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
							DNSID:  "12346",
						},
					},
				},
			},
			expectedMethod: "zone/add_aaaa",
			wantErr:        false,
		},
		{
			name:       "add CNAME record",
			recordType: "CNAME",
			content:    "example.github.io",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
							DNSID:  "67890",
						},
					},
				},
			},
			expectedMethod: "zone/add_cname",
			wantErr:        false,
		},
		{
			name:       "add MX record",
			recordType: "MX",
			content:    "10 mail.example.com",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
							DNSID:  "12347",
						},
					},
				},
			},
			expectedMethod: "zone/add_mx",
			wantErr:        false,
		},
		{
			name:       "add NS record",
			recordType: "NS",
			content:    "ns1.example.com",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
							DNSID:  "12348",
						},
					},
				},
			},
			expectedMethod: "zone/add_ns",
			wantErr:        false,
		},
		{
			name:       "add TXT record",
			recordType: "TXT",
			content:    "v=spf1 include:_spf.example.com ~all",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
							DNSID:  "12349",
						},
					},
				},
			},
			expectedMethod: "zone/add_txt",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(t, tt.response, http.StatusOK)
			defer server.Close()

			client := setupTestClient(t, server)

			params := CreateDNSRecordParams{
				Name:    "www",
				Type:    tt.recordType,
				Content: tt.content,
				TTL:     3600,
			}

			record, err := client.AddRR(context.Background(), "example.com", params)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, params.Name, record.Name)
			assert.Equal(t, params.Type, record.Type)
			assert.Equal(t, params.Content, record.Content)
		})
	}
}

func TestClient_AddRR_UnsupportedType(t *testing.T) {
	client := NewClient("username", "password")

	params := CreateDNSRecordParams{
		Name:    "www",
		Type:    "UNSUPPORTED",
		Content: "192.0.2.1",
	}

	_, err := client.AddRR(context.Background(), "example.com", params)
	require.Error(t, err)

	var unsupportedErr *UnsupportedRecordTypeError
	assert.True(t, errors.As(err, &unsupportedErr), "error should be UnsupportedRecordTypeError")
}

func TestClient_ListZones(t *testing.T) {
	response := ServiceListResponse{
		Answer: ServiceListAnswer{
			Services: []Service{
				{
					ServiceType: "domain",
					Domain:      "example.com",
					ServiceID:   "12345",
				},
				{
					ServiceType: "domain",
					Domain:      "test.com",
					ServiceID:   "67890",
				},
				{
					ServiceType: "hosting",
					Domain:      "other.com",
					ServiceID:   "11111",
				},
			},
		},
	}

	server := setupTestServer(t, response, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	zones, err := client.ListZones(context.Background())
	require.NoError(t, err)
	require.Len(t, zones, 2)
	assert.Equal(t, "example.com", zones[0].Name)
	assert.Equal(t, "12345", zones[0].ID)
}

func TestClient_ListZonesByName(t *testing.T) {
	response := ServiceListResponse{
		Answer: ServiceListAnswer{
			Services: []Service{
				{
					ServiceType: "domain",
					Domain:      "example.com",
					ServiceID:   "12345",
				},
				{
					ServiceType: "domain",
					Domain:      "test.com",
					ServiceID:   "67890",
				},
			},
		},
	}

	server := setupTestServer(t, response, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	zones, err := client.ListZonesByName(context.Background(), "example.com")
	require.NoError(t, err)
	require.Len(t, zones, 1)
	assert.Equal(t, "example.com", zones[0].Name)
}

func TestClient_ListRecords(t *testing.T) {
	response := ZoneListResponse{
		Answer: ZoneListAnswer{
			Domains: []DomainWithRecords{
				{
					DName: "example.com",
					NSList: []NSRecord{
						{
							Subdomain: "www",
							Type:      "A",
							Content:   "192.0.2.1",
							TTL:       3600,
							DNSID:     "11111",
						},
						{
							Subdomain: "@",
							Type:      "A",
							Content:   "192.0.2.2",
							TTL:       3600,
							DNSID:     "22222",
						},
						{
							Subdomain: "mail",
							Type:      "MX",
							Content:   "10 mail.example.com",
							TTL:       3600,
							DNSID:     "33333",
						},
					},
				},
			},
		},
	}

	server := setupTestServer(t, response, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	params := ListDNSRecordsParams{
		ZoneName: "example.com",
	}

	records, err := client.ListRecords(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, records, 3)
	assert.Equal(t, "www", records[0].Name)
	assert.Equal(t, "A", records[0].Type)
}

func TestClient_ListRecords_WithFilters(t *testing.T) {
	response := ZoneListResponse{
		Answer: ZoneListAnswer{
			Domains: []DomainWithRecords{
				{
					DName: "example.com",
					NSList: []NSRecord{
						{
							Subdomain: "www",
							Type:      "A",
							Content:   "192.0.2.1",
							TTL:       3600,
							DNSID:     "11111",
						},
						{
							Subdomain: "@",
							Type:      "A",
							Content:   "192.0.2.2",
							TTL:       3600,
							DNSID:     "22222",
						},
						{
							Subdomain: "mail",
							Type:      "MX",
							Content:   "10 mail.example.com",
							TTL:       3600,
							DNSID:     "33333",
						},
					},
				},
			},
		},
	}

	server := setupTestServer(t, response, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	params := ListDNSRecordsParams{
		ZoneName: "example.com",
		Type:     "A", // Filter by type
	}

	records, err := client.ListRecords(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, records, 2)
	for _, record := range records {
		assert.Equal(t, "A", record.Type)
	}
}

func TestClient_GetRRByName(t *testing.T) {
	response := ZoneListResponse{
		Answer: ZoneListAnswer{
			Domains: []DomainWithRecords{
				{
					DName: "example.com",
					NSList: []NSRecord{
						{
							Subdomain: "www",
							Type:      "A",
							Content:   "192.0.2.1",
							TTL:       3600,
							DNSID:     "11111",
						},
						{
							Subdomain: "@",
							Type:      "A",
							Content:   "192.0.2.2",
							TTL:       3600,
							DNSID:     "22222",
						},
					},
				},
			},
		},
	}

	server := setupTestServer(t, response, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	record, err := client.GetRRByName(context.Background(), "example.com", "www")
	require.NoError(t, err)
	assert.Equal(t, "www", record.Name)
	assert.Equal(t, "192.0.2.1", record.Content)
}

func TestClient_GetRRByName_NotFound(t *testing.T) {
	response := ZoneListResponse{
		Answer: ZoneListAnswer{
			Domains: []DomainWithRecords{
				{
					DName: "example.com",
					NSList: []NSRecord{
						{
							Subdomain: "www",
							Type:      "A",
							Content:   "192.0.2.1",
							TTL:       3600,
							DNSID:     "11111",
						},
					},
				},
			},
		},
	}

	server := setupTestServer(t, response, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	_, err := client.GetRRByName(context.Background(), "example.com", "nonexistent")
	require.Error(t, err)

	var notFoundErr *RecordNotFoundError
	assert.True(t, errors.As(err, &notFoundErr), "error should be RecordNotFoundError")
}

func TestClient_DeleteRR(t *testing.T) {
	tests := []struct {
		name       string
		recordType string
		content    string
		response   AddNSResponse
		wantErr    bool
	}{
		{
			name:       "delete A record",
			recordType: "A",
			content:    "192.0.2.1",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "delete AAAA record",
			recordType: "AAAA",
			content:    "2001:db8::1",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "delete CNAME record",
			recordType: "CNAME",
			content:    "example.github.io",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "delete MX record",
			recordType: "MX",
			content:    "10 mail.example.com",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "delete NS record",
			recordType: "NS",
			content:    "ns1.example.com",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "delete TXT record",
			recordType: "TXT",
			content:    "v=spf1 include:_spf.example.com ~all",
			response: AddNSResponse{
				Answer: AddNSAnswer{
					Domains: []DomainResult{
						{
							DName:  "example.com",
							Result: "success",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(t, tt.response, http.StatusOK)
			defer server.Close()

			client := setupTestClient(t, server)

			record := DNSRecord{
				Name:    "www",
				Type:    tt.recordType,
				Content: tt.content,
			}

			err := client.DeleteRR(context.Background(), "example.com", record)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_DeleteRR_UnsupportedType(t *testing.T) {
	client := NewClient("username", "password")

	record := DNSRecord{
		Name:    "www",
		Type:    "UNSUPPORTED",
		Content: "192.0.2.1",
	}

	err := client.DeleteRR(context.Background(), "example.com", record)
	require.Error(t, err)

	var unsupportedErr *UnsupportedRecordTypeError
	assert.True(t, errors.As(err, &unsupportedErr), "error should be UnsupportedRecordTypeError")
}

func TestClient_UpdateRR(t *testing.T) {
	// First call - delete
	deleteResponse := AddNSResponse{
		Answer: AddNSAnswer{
			Domains: []DomainResult{
				{
					DName:  "example.com",
					Result: "success",
				},
			},
		},
	}

	// Second call - add
	addResponse := AddNSResponse{
		Answer: AddNSAnswer{
			Domains: []DomainResult{
				{
					DName:  "example.com",
					Result: "success",
					DNSID:  "12345",
				},
			},
		},
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		callCount++
		if callCount == 1 {
			// Delete response
			require.NoError(t, json.NewEncoder(w).Encode(deleteResponse))
		} else {
			// Add response
			require.NoError(t, json.NewEncoder(w).Encode(addResponse))
		}
	}))
	defer server.Close()

	client := setupTestClient(t, server)

	record := DNSRecord{
		Name:    "www",
		Type:    "A",
		Content: "192.0.2.3",
		TTL:     3600,
	}

	updated, err := client.UpdateRR(context.Background(), "example.com", record)
	require.NoError(t, err)
	assert.Equal(t, "192.0.2.3", updated.Content)
	assert.Equal(t, 2, callCount, "UpdateRR should make 2 API calls")
}

func TestClient_apiRequest_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := setupTestClient(t, server)

	req := &ServiceListRequest{
		BaseRequest: BaseRequest{},
		PageSize:    1000,
	}

	_, err := client.apiRequest(context.Background(), "service/get_list", req)
	require.Error(t, err)

	var httpErr *HTTPError
	assert.True(t, errors.As(err, &httpErr), "error should be HTTPError")
	assert.Equal(t, http.StatusInternalServerError, httpErr.StatusCode)
}

func TestClient_apiRequest_APIError(t *testing.T) {
	response := APIResponse{
		ErrorText: "Invalid credentials",
	}

	server := setupTestServer(t, response, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	req := &ServiceListRequest{
		BaseRequest: BaseRequest{},
		PageSize:    1000,
	}

	_, err := client.apiRequest(context.Background(), "service/get_list", req)
	require.Error(t, err)

	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr), "error should be APIError")
	assert.Equal(t, "Invalid credentials", apiErr.Message)
}

func TestClient_ListRecordsByZoneID(t *testing.T) {
	// First call - list zones
	zonesResponse := ServiceListResponse{
		Answer: ServiceListAnswer{
			Services: []Service{
				{
					ServiceType: "domain",
					Domain:      "example.com",
					ServiceID:   "12345",
				},
			},
		},
	}

	// Second call - list records
	recordsResponse := ZoneListResponse{
		Answer: ZoneListAnswer{
			Domains: []DomainWithRecords{
				{
					DName: "example.com",
					NSList: []NSRecord{
						{
							Subdomain: "www",
							Type:      "A",
							Content:   "192.0.2.1",
							TTL:       3600,
							DNSID:     "11111",
						},
					},
				},
			},
		},
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		callCount++
		if callCount == 1 {
			require.NoError(t, json.NewEncoder(w).Encode(zonesResponse))
		} else {
			require.NoError(t, json.NewEncoder(w).Encode(recordsResponse))
		}
	}))
	defer server.Close()

	client := setupTestClient(t, server)

	params := ListDNSRecordsParams{}

	records, err := client.ListRecordsByZoneID(context.Background(), "12345", params)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "www", records[0].Name)
}

func TestClient_ListRecordsByZoneID_ZoneNotFound(t *testing.T) {
	zonesResponse := ServiceListResponse{
		Answer: ServiceListAnswer{
			Services: []Service{
				{
					ServiceType: "domain",
					Domain:      "example.com",
					ServiceID:   "12345",
				},
			},
		},
	}

	server := setupTestServer(t, zonesResponse, http.StatusOK)
	defer server.Close()

	client := setupTestClient(t, server)

	params := ListDNSRecordsParams{}

	_, err := client.ListRecordsByZoneID(context.Background(), "nonexistent", params)
	require.Error(t, err)

	var notFoundErr *ZoneNotFoundError
	assert.True(t, errors.As(err, &notFoundErr), "error should be ZoneNotFoundError")
}
