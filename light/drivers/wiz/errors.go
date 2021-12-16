// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package wiz

import "fmt"

// QueryErrorCode represents a code that was returned by the device as a response to a query.
type QueryErrorCode int64

const (
	QueryErrorCodeInvalidParams = -32602 // Parameter is outside the valid range or not available at all.
)

// ErrQueryFailed is returned if the device responds with an error message.
type ErrQueryFailed struct {
	errorCode QueryErrorCode
	message   string
}

func (e *ErrQueryFailed) Error() string {
	return fmt.Sprintf("light bulb returned error code %d: %v", e.errorCode, e.message)
}

// QueryErrorCode returns the error code that was returned by the device in response to a query.
func (e *ErrQueryFailed) QueryErrorCode() QueryErrorCode {
	return e.errorCode
}

// Code returns the error message that was returned by the device.
func (e *ErrQueryFailed) Message() string {
	return e.message
}
