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

package regru_test

import (
	"context"
	"fmt"
	"time"

	"github.com/mixanemca/regru-go"
)

func ExampleNewClient() {
	// Create a new client with default settings
	_ = regru.NewClient("your-username", "your-password")

	// Or create with custom options
	_ = regru.NewClient(
		"your-username",
		"your-password",
		regru.WithBaseURL("https://api.reg.ru/api/regru2"),
		regru.WithTimeout(60*time.Second),
	)
}

func ExampleClient_ListZones() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	zones, err := client.ListZones(ctx)
	if err != nil {
		// Handle error (e.g., check for HTTPError, APIError)
		return
	}

	for _, zone := range zones {
		fmt.Printf("Zone: %s (ID: %s)\n", zone.Name, zone.ID)
	}
}

func ExampleClient_ListZonesByName() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	zones, err := client.ListZonesByName(ctx, "example.com")
	if err != nil {
		return
	}

	for _, zone := range zones {
		fmt.Printf("Found zone: %s\n", zone.Name)
	}
}

func ExampleClient_AddRR() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	// Add an A record
	params := regru.CreateDNSRecordParams{
		Name:    "www",
		Type:    "A",
		Content: "192.0.2.1",
		TTL:     3600,
	}

	record, err := client.AddRR(ctx, "example.com", params)
	if err != nil {
		return
	}

	fmt.Printf("Created %s record: %s -> %s\n", record.Type, record.Name, record.Content)
}

func ExampleClient_AddRR_cname() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	// Add a CNAME record
	params := regru.CreateDNSRecordParams{
		Name:    "blog",
		Type:    "CNAME",
		Content: "example.github.io",
		TTL:     3600,
	}

	record, err := client.AddRR(ctx, "example.com", params)
	if err != nil {
		return
	}

	fmt.Printf("Created %s record: %s -> %s\n", record.Type, record.Name, record.Content)
}

func ExampleClient_ListRecords() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	params := regru.ListDNSRecordsParams{
		ZoneName: "example.com",
		Type:     "A", // Filter by record type
	}

	records, err := client.ListRecords(ctx, params)
	if err != nil {
		return
	}

	for _, record := range records {
		fmt.Printf("%s %s -> %s\n", record.Type, record.Name, record.Content)
	}
}

func ExampleClient_GetRRByName() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	record, err := client.GetRRByName(ctx, "example.com", "www")
	if err != nil {
		return
	}

	fmt.Printf("Record: %s %s -> %s\n", record.Type, record.Name, record.Content)
}

func ExampleClient_UpdateRR() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	// Update existing record
	record := regru.DNSRecord{
		Name:    "www",
		Type:    "A",
		Content: "192.0.2.3", // New IP address
		TTL:     3600,
	}

	updated, err := client.UpdateRR(ctx, "example.com", record)
	if err != nil {
		return
	}

	fmt.Printf("Updated record: %s -> %s\n", updated.Name, updated.Content)
}

func ExampleClient_DeleteRR() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	record := regru.DNSRecord{
		Name:    "www",
		Type:    "A",
		Content: "192.0.2.1",
	}

	err := client.DeleteRR(ctx, "example.com", record)
	if err != nil {
		return
	}

	fmt.Println("Record deleted successfully")
}

func ExampleClient_ListRecordsByZoneID() {
	client := regru.NewClient("your-username", "your-password")
	ctx := context.Background()

	params := regru.ListDNSRecordsParams{
		Type: "A", // Optional: filter by type
	}

	records, err := client.ListRecordsByZoneID(ctx, "zone-id-123", params)
	if err != nil {
		return
	}

	fmt.Printf("Found %d records\n", len(records))
}
