// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

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
