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

import "fmt"

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
	ServiceType string      `json:"service_type,omitempty"`
	ServType    string      `json:"servtype,omitempty"` // Alternative field name used by some API methods
	Domain      string      `json:"domain,omitempty"`
	DName       string      `json:"dname,omitempty"`      // Alternative field name for domain name
	ServiceID   interface{} `json:"service_id,omitempty"` // Can be int or string depending on API method
}

// GetServiceType returns the service type, checking both possible field names.
func (s *Service) GetServiceType() string {
	if s.ServiceType != "" {
		return s.ServiceType
	}
	return s.ServType
}

// GetDomain returns the domain name, checking both possible field names.
func (s *Service) GetDomain() string {
	if s.Domain != "" {
		return s.Domain
	}
	return s.DName
}

// GetServiceID returns the service ID as a string.
func (s *Service) GetServiceID() string {
	switch v := s.ServiceID.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ZoneListResponse represents the response for zone/get_ns (for backward compatibility).
type ZoneListResponse struct {
	Answer ZoneListAnswer `json:"answer,omitempty"`
}

// ZoneListAnswer contains the list of domains with their DNS records.
type ZoneListAnswer struct {
	Domains []DomainWithRecords `json:"domains,omitempty"`
}

// DomainWithRecords represents a domain with its DNS records (for zone/get_ns).
type DomainWithRecords struct {
	DName  string     `json:"dname,omitempty"`
	NSList []NSRecord `json:"ns_list,omitempty"`
}

// NSRecord represents a DNS record in reg.ru API format (for zone/get_ns).
type NSRecord struct {
	Subdomain string `json:"subdomain,omitempty"`
	Type      string `json:"type,omitempty"`
	Content   string `json:"content,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	DNSID     string `json:"dns_id,omitempty"`
}

// ZoneGetResourceRecordsResponse represents the response for zone/get_resource_records.
type ZoneGetResourceRecordsResponse struct {
	Answer ZoneGetResourceRecordsAnswer `json:"answer,omitempty"`
	Result string                       `json:"result,omitempty"`
}

// ZoneGetResourceRecordsAnswer contains the list of domains with their resource records.
type ZoneGetResourceRecordsAnswer struct {
	Domains []DomainWithResourceRecords `json:"domains,omitempty"`
}

// DomainWithResourceRecords represents a domain with its resource records.
type DomainWithResourceRecords struct {
	DName     string           `json:"dname,omitempty"`
	Result    string           `json:"result,omitempty"`
	RRList    []ResourceRecord `json:"rrs,omitempty"`
	ServiceID string           `json:"service_id,omitempty"`
	SOA       *SOAInfo         `json:"soa,omitempty"`
}

// ResourceRecord represents a DNS resource record in zone/get_resource_records format.
type ResourceRecord struct {
	Content string      `json:"content,omitempty"`
	Prio    interface{} `json:"prio,omitempty"` // Can be number or string
	Rectype string      `json:"rectype,omitempty"`
	State   string      `json:"state,omitempty"`
	Subname string      `json:"subname,omitempty"`
}

// GetPrio returns the priority as a string.
func (r *ResourceRecord) GetPrio() string {
	if r.Prio == nil {
		return ""
	}
	switch v := r.Prio.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// SOAInfo represents SOA record information.
type SOAInfo struct {
	MinimumTTL string `json:"minimum_ttl,omitempty"`
	TTL        string `json:"ttl,omitempty"`
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
