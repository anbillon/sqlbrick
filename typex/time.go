// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.
//
// Package typex implements the Scanner interface and the driver
// Valuer interface so we can use them in database operation.

package typex

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// NullTime represents a time.Time that may be null. NullTime implements the
// sql.Scanner interface so it can be used as a scan destination, similar to
// sql.NullString.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		nt.Time = x
	case nil:
		nt.Valid = false
		return nil
	default:
		err = fmt.Errorf("cannot scan type %T into null.Time: %v", value, value)
	}
	nt.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
