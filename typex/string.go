// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.
//
// Package typex implements the Scanner interface and the driver
// Valuer interface so we can use them in database operation.
package typex

import (
	"database/sql"
)

// NullString wraps sql.NullString.
type NullString struct {
	sql.NullString
}

// NewNullString will create a new NullString.
func NewNullString(s string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: s,
			Valid:  true,
		}}
}
