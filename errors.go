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
	"fmt"
)

// Predefined errors that can be checked with errors.Is().
var (
	// ErrUnsupportedRecordType is returned when an unsupported DNS record type is used.
	ErrUnsupportedRecordType = errors.New("unsupported record type")

	// ErrRecordNotFound is returned when a DNS record is not found.
	ErrRecordNotFound = errors.New("record not found")

	// ErrZoneNotFound is returned when a zone is not found.
	ErrZoneNotFound = errors.New("zone not found")
)

// APIError represents an error returned by the reg.ru API.
type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s", e.Message)
}

// HTTPError represents an HTTP error with status code.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("API returned status %d: %s", e.StatusCode, e.Body)
}

// UnsupportedRecordTypeError represents an error for unsupported record type.
type UnsupportedRecordTypeError struct {
	RecordType string
}

func (e *UnsupportedRecordTypeError) Error() string {
	return fmt.Sprintf("unsupported record type: %s", e.RecordType)
}

func (e *UnsupportedRecordTypeError) Is(target error) bool {
	return target == ErrUnsupportedRecordType
}

// RecordNotFoundError represents an error when a record is not found.
type RecordNotFoundError struct {
	RecordName string
}

func (e *RecordNotFoundError) Error() string {
	return fmt.Sprintf("record not found: %s", e.RecordName)
}

func (e *RecordNotFoundError) Is(target error) bool {
	return target == ErrRecordNotFound
}

// ZoneNotFoundError represents an error when a zone is not found.
type ZoneNotFoundError struct {
	ZoneID string
}

func (e *ZoneNotFoundError) Error() string {
	return fmt.Sprintf("zone not found: %s", e.ZoneID)
}

func (e *ZoneNotFoundError) Is(target error) bool {
	return target == ErrZoneNotFound
}
