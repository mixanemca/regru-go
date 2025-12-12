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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAddRecordPath(t *testing.T) {
	tests := []struct {
		name        string
		recordType  string
		wantPath    string
		wantErr     bool
		wantErrType error
	}{
		{
			name:       "A record",
			recordType: RecordTypeA,
			wantPath:   "zone/add_alias",
			wantErr:    false,
		},
		{
			name:       "AAAA record",
			recordType: RecordTypeAAAA,
			wantPath:   "zone/add_aaaa",
			wantErr:    false,
		},
		{
			name:       "CNAME record",
			recordType: RecordTypeCNAME,
			wantPath:   "zone/add_cname",
			wantErr:    false,
		},
		{
			name:       "MX record",
			recordType: RecordTypeMX,
			wantPath:   "zone/add_mx",
			wantErr:    false,
		},
		{
			name:       "NS record",
			recordType: RecordTypeNS,
			wantPath:   "zone/add_ns",
			wantErr:    false,
		},
		{
			name:       "TXT record",
			recordType: RecordTypeTXT,
			wantPath:   "zone/add_txt",
			wantErr:    false,
		},
		{
			name:        "unsupported record type",
			recordType:  "UNSUPPORTED",
			wantPath:    "",
			wantErr:     true,
			wantErrType: ErrUnsupportedRecordType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := getAddRecordPath(tt.recordType)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErrType))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPath, path)
			}
		})
	}
}

func TestGetRemoveRecordPath(t *testing.T) {
	tests := []struct {
		name        string
		recordType  string
		wantPath    string
		wantErr     bool
		wantErrType error
	}{
		{
			name:       "A record",
			recordType: RecordTypeA,
			wantPath:   "zone/remove_record",
			wantErr:    false,
		},
		{
			name:       "AAAA record",
			recordType: RecordTypeAAAA,
			wantPath:   "zone/remove_record",
			wantErr:    false,
		},
		{
			name:       "CNAME record",
			recordType: RecordTypeCNAME,
			wantPath:   "zone/remove_record",
			wantErr:    false,
		},
		{
			name:       "MX record",
			recordType: RecordTypeMX,
			wantPath:   "zone/remove_record",
			wantErr:    false,
		},
		{
			name:       "NS record",
			recordType: RecordTypeNS,
			wantPath:   "zone/remove_record",
			wantErr:    false,
		},
		{
			name:       "TXT record",
			recordType: RecordTypeTXT,
			wantPath:   "zone/remove_record",
			wantErr:    false,
		},
		{
			name:        "unsupported record type",
			recordType:  "UNSUPPORTED",
			wantPath:    "",
			wantErr:     true,
			wantErrType: ErrUnsupportedRecordType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := getRemoveRecordPath(tt.recordType)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErrType))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPath, path)
			}
		})
	}
}

func TestCreateAddRecordRequest(t *testing.T) {
	tests := []struct {
		name      string
		zone      string
		params    CreateDNSRecordParams
		wantErr   bool
		wantType  string
		checkType func(APIRequest) bool
	}{
		{
			name: "A record request",
			zone: "example.com",
			params: CreateDNSRecordParams{
				Name:    "www",
				Type:    RecordTypeA,
				Content: "192.0.2.1",
				TTL:     3600,
			},
			wantErr:  false,
			wantType: RecordTypeA,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*AddAliasRequest)
				return ok
			},
		},
		{
			name: "AAAA record request",
			zone: "example.com",
			params: CreateDNSRecordParams{
				Name:    "www",
				Type:    RecordTypeAAAA,
				Content: "2001:db8::1",
				TTL:     3600,
			},
			wantErr:  false,
			wantType: RecordTypeAAAA,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*AddAAAARequest)
				return ok
			},
		},
		{
			name: "CNAME record request",
			zone: "example.com",
			params: CreateDNSRecordParams{
				Name:    "blog",
				Type:    RecordTypeCNAME,
				Content: "example.github.io",
				TTL:     3600,
			},
			wantErr:  false,
			wantType: RecordTypeCNAME,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*AddCNAMERequest)
				return ok
			},
		},
		{
			name: "MX record request",
			zone: "example.com",
			params: CreateDNSRecordParams{
				Name:    "@",
				Type:    RecordTypeMX,
				Content: "10 mail.example.com",
				TTL:     3600,
			},
			wantErr:  false,
			wantType: RecordTypeMX,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*AddMXRequest)
				return ok
			},
		},
		{
			name: "NS record request",
			zone: "example.com",
			params: CreateDNSRecordParams{
				Name:    "@",
				Type:    RecordTypeNS,
				Content: "ns1.example.com",
				TTL:     3600,
			},
			wantErr:  false,
			wantType: RecordTypeNS,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*AddNSRequest)
				return ok
			},
		},
		{
			name: "TXT record request",
			zone: "example.com",
			params: CreateDNSRecordParams{
				Name:    "@",
				Type:    RecordTypeTXT,
				Content: "v=spf1 include:_spf.example.com ~all",
				TTL:     3600,
			},
			wantErr:  false,
			wantType: RecordTypeTXT,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*AddTXTRequest)
				return ok
			},
		},
		{
			name: "unsupported record type",
			zone: "example.com",
			params: CreateDNSRecordParams{
				Name:    "www",
				Type:    "UNSUPPORTED",
				Content: "value",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createAddRecordRequest(tt.zone, tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, req, "createAddRecordRequest() should not return nil request")
			assert.True(t, tt.checkType(req), "createAddRecordRequest() returned wrong request type for %s", tt.wantType)

			// Check that credentials can be set
			req.SetCredentials("test-user", "test-pass")
		})
	}
}

func TestCreateRemoveRecordRequest(t *testing.T) {
	tests := []struct {
		name      string
		zone      string
		record    DNSRecord
		wantErr   bool
		wantType  string
		checkType func(APIRequest) bool
	}{
		{
			name: "A record request",
			zone: "example.com",
			record: DNSRecord{
				Name:    "www",
				Type:    RecordTypeA,
				Content: "192.0.2.1",
			},
			wantErr:  false,
			wantType: RecordTypeA,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*RemoveAliasRequest)
				return ok
			},
		},
		{
			name: "AAAA record request",
			zone: "example.com",
			record: DNSRecord{
				Name:    "www",
				Type:    RecordTypeAAAA,
				Content: "2001:db8::1",
			},
			wantErr:  false,
			wantType: RecordTypeAAAA,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*RemoveAAAARequest)
				return ok
			},
		},
		{
			name: "CNAME record request",
			zone: "example.com",
			record: DNSRecord{
				Name:    "blog",
				Type:    RecordTypeCNAME,
				Content: "example.github.io",
			},
			wantErr:  false,
			wantType: RecordTypeCNAME,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*RemoveCNAMERequest)
				return ok
			},
		},
		{
			name: "MX record request",
			zone: "example.com",
			record: DNSRecord{
				Name:    "@",
				Type:    RecordTypeMX,
				Content: "10 mail.example.com",
			},
			wantErr:  false,
			wantType: RecordTypeMX,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*RemoveMXRequest)
				return ok
			},
		},
		{
			name: "NS record request",
			zone: "example.com",
			record: DNSRecord{
				Name:    "@",
				Type:    RecordTypeNS,
				Content: "ns1.example.com",
			},
			wantErr:  false,
			wantType: RecordTypeNS,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*RemoveNSRequest)
				return ok
			},
		},
		{
			name: "TXT record request",
			zone: "example.com",
			record: DNSRecord{
				Name:    "@",
				Type:    RecordTypeTXT,
				Content: "v=spf1 include:_spf.example.com ~all",
			},
			wantErr:  false,
			wantType: RecordTypeTXT,
			checkType: func(req APIRequest) bool {
				_, ok := req.(*RemoveTXTRequest)
				return ok
			},
		},
		{
			name: "unsupported record type",
			zone: "example.com",
			record: DNSRecord{
				Name:    "www",
				Type:    "UNSUPPORTED",
				Content: "value",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createRemoveRecordRequest(tt.zone, tt.record)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, req, "createRemoveRecordRequest() should not return nil request")
			assert.True(t, tt.checkType(req), "createRemoveRecordRequest() returned wrong request type for %s", tt.wantType)

			// Check that credentials can be set
			req.SetCredentials("test-user", "test-pass")
		})
	}
}
