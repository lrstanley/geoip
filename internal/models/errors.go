// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lrstanley/chix/v2"
)

// ErrHostResolve is an error that is returned when the address didn't match an IP,
// and thus a hostname lookup was attempted, but failed.
type ErrHostResolve struct { //nolint:errname // stable exported API
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
type ErrInternalAddress struct { //nolint:errname // stable exported API
	Address string
}

func (e *ErrInternalAddress) Error() string {
	return fmt.Sprintf("internal address specified: %v", e.Address)
}

type ErrNotFound struct { //nolint:errname // stable exported API
	Address string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("address not found: %v", e.Address)
}

type ErrRateLimitExceeded struct { //nolint:errname // stable exported API
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

	_, ok0 := errors.AsType[*ErrHostResolve](err)
	_, ok1 := errors.AsType[*ErrInternalAddress](err)
	_, ok2 := errors.AsType[*ErrNotFound](err)
	_, ok3 := errors.AsType[*ErrRateLimitExceeded](err)

	return ok0 || ok1 || ok2 || ok3
}

// ErrorResolver maps internally defined errors to HTTP status codes for chix.
// Returns nil when the error is not one of the known application errors.
func ErrorResolver(oerr *chix.ResolvedError) *chix.ResolvedError {
	if oerr == nil || oerr.Err == nil {
		return nil
	}

	err := oerr.Err

	if _, ok := errors.AsType[*ErrRateLimitExceeded](err); ok {
		oerr.StatusCode = http.StatusTooManyRequests
		oerr.Visibility = chix.ErrorPublic
		return oerr
	}
	if _, ok := errors.AsType[*ErrNotFound](err); ok {
		oerr.StatusCode = http.StatusNotFound
		oerr.Visibility = chix.ErrorPublic
		return oerr
	}
	if _, ok := errors.AsType[*ErrHostResolve](err); ok {
		oerr.StatusCode = http.StatusNotFound
		oerr.Visibility = chix.ErrorPublic
		return oerr
	}

	if IsClientError(err) {
		oerr.StatusCode = http.StatusBadRequest
		oerr.Visibility = chix.ErrorPublic
		return oerr
	}

	return nil
}
