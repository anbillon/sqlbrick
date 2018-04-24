// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.
//
// Package typex implements the Scanner interface and the driver
// Valuer interface so we can use them in database operation.
package typex

import (
	"database/sql/driver"
)

// NullTime represents a []byte that may be null. NullBytes implements the
// sql.Scanner interface so it can be used as a scan destination, similar to
// sql.NullString.
type NullBytes struct {
	Bytes []byte
	Valid bool
}

// NewNullBytes create a new NullBytes. It will check if the given bytes is valid.
func NewNullBytes(b []byte) NullBytes {
	valid := true
	if b == nil {
		valid = false
	}
	return NullBytes{Bytes: b, Valid: valid}
}

// Scan implements the Scanner interface.
func (nb *NullBytes) Scan(value interface{}) error {
	if value == nil {
		nb.Bytes, nb.Valid = []byte{}, false
		return nil
	}
	nb.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (nb NullBytes) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bytes, nil
}
