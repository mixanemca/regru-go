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

// AddRecordRequest represents base structure for adding DNS records.
type AddRecordRequest struct {
	BaseRequest
	Domains []AddRecordDomain `json:"domains"`
	TTL     int               `json:"ttl,omitempty"`
}

// AddNSRequest represents parameters for zone/add_ns API method.
type AddNSRequest struct {
	AddRecordRequest
}

// AddAliasRequest represents parameters for zone/add_alias API method (A records).
type AddAliasRequest struct {
	AddRecordRequest
}

// AddAAAARequest represents parameters for zone/add_aaaa API method.
type AddAAAARequest struct {
	AddRecordRequest
}

// AddCNAMERequest represents parameters for zone/add_cname API method.
type AddCNAMERequest struct {
	AddRecordRequest
}

// AddMXRequest represents parameters for zone/add_mx API method.
type AddMXRequest struct {
	AddRecordRequest
	Priority int `json:"priority,omitempty"`
}

// AddTXTRequest represents parameters for zone/add_txt API method.
type AddTXTRequest struct {
	AddRecordRequest
}

// RemoveRecordDomain represents a domain in remove record requests.
type RemoveRecordDomain struct {
	DName     string `json:"dname"`
	Subdomain string `json:"subdomain"`
	Content   string `json:"content"`
}

// RemoveRecordRequest represents base structure for removing DNS records.
type RemoveRecordRequest struct {
	BaseRequest
	Domains []RemoveRecordDomain `json:"domains"`
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
