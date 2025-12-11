# regru-go

Go client library for working with reg.ru API.

## Description

`regru-go` provides a convenient interface for working with DNS zones and records through reg.ru API. The library allows you to manage DNS zones and records through a simple and intuitive API.

## Installation

```bash
go get github.com/mixanemca/regru-go
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/mixanemca/regru-go"
)

func main() {
    // Create client
    client := regru.NewClient("your-username", "your-password")
    
    ctx := context.Background()
    
    // Get list of zones
    zones, err := client.ListZones(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d zones\n", len(zones))
    for _, zone := range zones {
        fmt.Printf("- %s (ID: %s)\n", zone.Name, zone.ID)
    }
    
    // Add a new DNS record
    createParams := regru.CreateDNSRecordParams{
        Name:    "www",
        Type:    "A",
        Content: "192.0.2.1",
        TTL:     3600,
    }
    
    record, err := client.AddRR(ctx, "example.com", createParams)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created record: %s\n", record.Name)
}
```

### Using Custom Settings

```go
import (
    "net/http"
    "time"
    
    "github.com/mixanemca/regru-go"
)

// Create client with custom HTTP client
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}

client := regru.NewClient(
    "your-username",
    "your-password",
    regru.WithHTTPClient(httpClient),
    regru.WithBaseURL("https://api.reg.ru/api/regru2"),
)
```

## API

### Client

The client provides the following methods for DNS management:

- `AddRR(ctx, zone, params)` - creates a new DNS record
- `DeleteRR(ctx, zone, rr)` - deletes a DNS record
- `GetRRByName(ctx, zone, name)` - gets a DNS record by name
- `ListZones(ctx)` - returns a list of all zones
- `ListZonesByName(ctx, name)` - returns zones by name
- `ListRecords(ctx, params)` - returns a list of DNS records for a zone
- `ListRecordsByZoneID(ctx, id, params)` - returns records by zone ID
- `UpdateRR(ctx, zone, rr)` - updates a DNS record

## Authentication

To work with reg.ru API, you need:
- Username (login)
- Password

**Important**: To work with the API, you need to configure access from trusted IP addresses. Details are available in the [reg.ru documentation](https://www.reg.ru/support/help/api2).

## Error Handling

The library provides typed errors that can be checked using `errors.Is()` and `errors.As()`:

```go
import (
    "errors"
    "github.com/mixanemca/regru-go"
)

// Check for specific error types
record, err := client.AddRR(ctx, "example.com", params)
if err != nil {
    // Check if it's an unsupported record type
    if errors.Is(err, regru.ErrUnsupportedRecordType) {
        // Handle unsupported type
    }
    
    // Check if it's a record not found error
    if errors.Is(err, regru.ErrRecordNotFound) {
        // Handle not found
    }
    
    // Extract typed error for more details
    var apiErr *regru.APIError
    if errors.As(err, &apiErr) {
        // Access API error message
        fmt.Printf("API error: %s\n", apiErr.Message)
    }
    
    var httpErr *regru.HTTPError
    if errors.As(err, &httpErr) {
        // Access HTTP status code and body
        fmt.Printf("HTTP %d: %s\n", httpErr.StatusCode, httpErr.Body)
    }
}
```

### Available Error Types

- `ErrUnsupportedRecordType` - returned when an unsupported DNS record type is used
- `ErrRecordNotFound` - returned when a DNS record is not found
- `ErrZoneNotFound` - returned when a zone is not found
- `APIError` - represents an error returned by the reg.ru API
- `HTTPError` - represents an HTTP error with status code
- `UnsupportedRecordTypeError` - typed error for unsupported record types
- `RecordNotFoundError` - typed error for record not found
- `ZoneNotFoundError` - typed error for zone not found

## API Documentation

Official reg.ru API documentation: https://www.reg.ru/reseller/api2doc

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Author

[mixanemca](https://github.com/mixanemca)