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

func TestAPIError(t *testing.T) {
	err := &APIError{Message: "test error"}
	assert.NotEmpty(t, err.Error(), "APIError.Error() should not return empty string")

	// Test that errors.As works
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr), "errors.As() should work with APIError")
	assert.Equal(t, "test error", apiErr.Message)
}

func TestHTTPError(t *testing.T) {
	err := &HTTPError{StatusCode: 500, Body: "Internal Server Error"}
	assert.NotEmpty(t, err.Error(), "HTTPError.Error() should not return empty string")

	// Test that errors.As works
	var httpErr *HTTPError
	require.True(t, errors.As(err, &httpErr), "errors.As() should work with HTTPError")
	assert.Equal(t, 500, httpErr.StatusCode)
	assert.Equal(t, "Internal Server Error", httpErr.Body)
}

func TestUnsupportedRecordTypeError(t *testing.T) {
	err := &UnsupportedRecordTypeError{RecordType: "UNSUPPORTED"}
	assert.NotEmpty(t, err.Error(), "UnsupportedRecordTypeError.Error() should not return empty string")
	assert.True(t, errors.Is(err, ErrUnsupportedRecordType), "UnsupportedRecordTypeError should be checkable with errors.Is()")

	var unsupportedErr *UnsupportedRecordTypeError
	require.True(t, errors.As(err, &unsupportedErr), "errors.As() should work with UnsupportedRecordTypeError")
	assert.Equal(t, "UNSUPPORTED", unsupportedErr.RecordType)
}

func TestRecordNotFoundError(t *testing.T) {
	err := &RecordNotFoundError{RecordName: "www"}
	assert.NotEmpty(t, err.Error(), "RecordNotFoundError.Error() should not return empty string")
	assert.True(t, errors.Is(err, ErrRecordNotFound), "RecordNotFoundError should be checkable with errors.Is()")

	var notFoundErr *RecordNotFoundError
	require.True(t, errors.As(err, &notFoundErr), "errors.As() should work with RecordNotFoundError")
	assert.Equal(t, "www", notFoundErr.RecordName)
}

func TestZoneNotFoundError(t *testing.T) {
	err := &ZoneNotFoundError{ZoneID: "12345"}
	assert.NotEmpty(t, err.Error(), "ZoneNotFoundError.Error() should not return empty string")
	assert.True(t, errors.Is(err, ErrZoneNotFound), "ZoneNotFoundError should be checkable with errors.Is()")

	var notFoundErr *ZoneNotFoundError
	require.True(t, errors.As(err, &notFoundErr), "errors.As() should work with ZoneNotFoundError")
	assert.Equal(t, "12345", notFoundErr.ZoneID)
}
