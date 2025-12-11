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

// APIResponse represents the base structure of reg.ru API response.
type APIResponse struct {
	Answer    interface{} `json:"answer,omitempty"`
	ErrorText string      `json:"error_text,omitempty"`
}

// ServiceListResponse represents the response for service/get_list.
type ServiceListResponse struct {
	Answer ServiceListAnswer `json:"answer,omitempty"`
}

// ServiceListAnswer contains the list of services.
type ServiceListAnswer struct {
	Services []Service `json:"services,omitempty"`
}

// Service represents a service in reg.ru API.
type Service struct {
	ServiceType string `json:"service_type,omitempty"`
	Domain      string `json:"domain,omitempty"`
	ServiceID   string `json:"service_id,omitempty"`
}

// ZoneListResponse represents the response for zone/get_ns.
type ZoneListResponse struct {
	Answer ZoneListAnswer `json:"answer,omitempty"`
}

// ZoneListAnswer contains the list of domains with their DNS records.
type ZoneListAnswer struct {
	Domains []DomainWithRecords `json:"domains,omitempty"`
}

// DomainWithRecords represents a domain with its DNS records.
type DomainWithRecords struct {
	DName  string     `json:"dname,omitempty"`
	NSList []NSRecord `json:"ns_list,omitempty"`
}

// NSRecord represents a DNS record in reg.ru API format.
type NSRecord struct {
	Subdomain string `json:"subdomain,omitempty"`
	Type      string `json:"type,omitempty"`
	Content   string `json:"content,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	DNSID     string `json:"dns_id,omitempty"`
}

// AddNSResponse represents the response for zone/add_ns.
type AddNSResponse struct {
	Answer AddNSAnswer `json:"answer,omitempty"`
}

// AddNSAnswer contains the result of adding a DNS record.
type AddNSAnswer struct {
	Domains []DomainResult `json:"domains,omitempty"`
}

// DomainResult represents the result of an operation on a domain.
type DomainResult struct {
	DName  string `json:"dname,omitempty"`
	Result string `json:"result,omitempty"`
	DNSID  string `json:"dns_id,omitempty"`
}
