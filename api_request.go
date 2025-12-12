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

// APIRequest represents an API request that can set credentials.
type APIRequest interface {
	SetCredentials(username, password string)
}

// BaseRequest contains common fields for all API requests.
type BaseRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SetCredentials sets username and password in the request.
func (b *BaseRequest) SetCredentials(username, password string) {
	b.Username = username
	b.Password = password
}

// AddRecordDomain represents a domain in add record requests.
type AddRecordDomain struct {
	DName     string `json:"dname"`
	Subdomain string `json:"subdomain"`
	Content   string `json:"content"`
}

// AddAliasDomain represents a domain in add_alias requests (only dname).
type AddAliasDomain struct {
	DName string `json:"dname"`
}

// AddRecordRequest represents base structure for adding DNS records.
type AddRecordRequest struct {
	BaseRequest
	Domains []AddRecordDomain `json:"domains"`
	TTL     int               `json:"ttl,omitempty"`
}

// AddNSRequest represents parameters for zone/add_ns API method.
// For add_ns, dns_server and subdomain are at the request level, not in domains.
type AddNSRequest struct {
	BaseRequest
	Domains      []AddAliasDomain `json:"domains"`
	Subdomain    string           `json:"subdomain"`
	DNSServer    string           `json:"dns_server"`
	RecordNumber string           `json:"record_number,omitempty"`
	TTL          int              `json:"ttl,omitempty"`
}

// AddAliasRequest represents parameters for zone/add_alias API method (A records).
// For add_alias, ipaddr and subdomain are at the request level, not in domains.
type AddAliasRequest struct {
	BaseRequest
	Domains   []AddAliasDomain `json:"domains"`
	Subdomain string           `json:"subdomain"`
	IPAddr    string           `json:"ipaddr"`
	TTL       int              `json:"ttl,omitempty"`
}

// AddAAAARequest represents parameters for zone/add_aaaa API method.
// For add_aaaa, ipaddr and subdomain are at the request level, not in domains.
type AddAAAARequest struct {
	BaseRequest
	Domains   []AddAliasDomain `json:"domains"`
	Subdomain string           `json:"subdomain"`
	IPAddr    string           `json:"ipaddr"`
	TTL       int              `json:"ttl,omitempty"`
}

// AddCNAMERequest represents parameters for zone/add_cname API method.
// For add_cname, canonical_name and subdomain are at the request level, not in domains.
type AddCNAMERequest struct {
	BaseRequest
	Domains       []AddAliasDomain `json:"domains"`
	Subdomain     string           `json:"subdomain"`
	CanonicalName string           `json:"canonical_name"`
	TTL           int              `json:"ttl,omitempty"`
}

// AddMXRequest represents parameters for zone/add_mx API method.
// For add_mx, mail_server and subdomain are at the request level, not in domains.
type AddMXRequest struct {
	BaseRequest
	Domains    []AddAliasDomain `json:"domains"`
	Subdomain  string           `json:"subdomain"`
	MailServer string           `json:"mail_server"`
	TTL        int              `json:"ttl,omitempty"`
}

// AddTXTRequest represents parameters for zone/add_txt API method.
// For add_txt, text and subdomain are at the request level, not in domains.
type AddTXTRequest struct {
	BaseRequest
	Domains   []AddAliasDomain `json:"domains"`
	Subdomain string           `json:"subdomain"`
	Text      string           `json:"text"`
	TTL       int              `json:"ttl,omitempty"`
}

// RemoveRecordDomain represents a domain in remove record requests.
type RemoveRecordDomain struct {
	DName string `json:"dname"`
}

// RemoveRecordRequest represents base structure for removing DNS records.
// For remove_record, subdomain, content, and record_type are at the request level.
type RemoveRecordRequest struct {
	BaseRequest
	Domains    []RemoveRecordDomain `json:"domains"`
	Subdomain  string               `json:"subdomain"`
	Content    string               `json:"content"`
	RecordType string               `json:"record_type"`
}

// RemoveNSRequest represents parameters for zone/remove_ns API method.
type RemoveNSRequest struct {
	RemoveRecordRequest
}

// RemoveAliasRequest represents parameters for zone/remove_alias API method (A records).
type RemoveAliasRequest struct {
	RemoveRecordRequest
}

// RemoveAAAARequest represents parameters for zone/remove_aaaa API method.
type RemoveAAAARequest struct {
	RemoveRecordRequest
}

// RemoveCNAMERequest represents parameters for zone/remove_cname API method.
type RemoveCNAMERequest struct {
	RemoveRecordRequest
}

// RemoveMXRequest represents parameters for zone/remove_mx API method.
type RemoveMXRequest struct {
	RemoveRecordRequest
}

// RemoveTXTRequest represents parameters for zone/remove_txt API method.
type RemoveTXTRequest struct {
	RemoveRecordRequest
}

// ServiceListRequest represents parameters for service/get_list API method.
type ServiceListRequest struct {
	BaseRequest
	PageSize int `json:"page_size,omitempty"`
}

// ZoneGetNSRequest represents parameters for zone/get_ns API method.
type ZoneGetNSRequest struct {
	BaseRequest
	Domains []string `json:"domains"`
}

// ZoneGetResourceRecordsRequest represents parameters for zone/get_resource_records API method.
type ZoneGetResourceRecordsRequest struct {
	BaseRequest
	Domains []ZoneGetResourceRecordsDomain `json:"domains"`
}

// ZoneGetResourceRecordsDomain represents a domain in get_resource_records request.
type ZoneGetResourceRecordsDomain struct {
	DName string `json:"dname"`
}
