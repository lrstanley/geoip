// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"fmt"
	"net/http"
)

// ErrHostResolve is an error that is returned when the address didn't match an IP,
// and thus a hostname lookup was attempted, but failed.
type ErrHostResolve struct {
	Err     error
	Address string
}

func (e *ErrHostResolve) Error() string {
	return fmt.Sprintf("failure during lookup of address: %v", e.Address)
}

func (e *ErrHostResolve) Unwrap() error {
	return e.Err
}

// ErrInternalAddress is an error that is returned when the address is an internal
// bogon address.
type ErrInternalAddress struct {
	Address string
}

func (e *ErrInternalAddress) Error() string {
	return fmt.Sprintf("internal address specified: %v", e.Address)
}

type ErrNotFound struct {
	Address string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("address not found: %v", e.Address)
}

type ErrRateLimitExceeded struct {
	Address string
}

func (e *ErrRateLimitExceeded) Error() string {
	return fmt.Sprintf("rate limit exceeded while looking up: %v", e.Address)
}

// IsClientError returns true if the error is related to user input/can be corrected
// by the user.
func IsClientError(err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case *ErrHostResolve,
		*ErrInternalAddress,
		*ErrNotFound,
		*ErrRateLimitExceeded:
		return true
	}

	return false
}

// ErrorResolver is used to map internally defined errors to HTTP status codes for
// chix.
func ErrorResolver(err error) (status int) {
	switch err.(type) {
	case *ErrRateLimitExceeded:
		return http.StatusTooManyRequests
	case *ErrNotFound, *ErrHostResolve:
		return http.StatusNotFound
	}

	if IsClientError(err) {
		return http.StatusBadRequest
	}

	return 0
}
