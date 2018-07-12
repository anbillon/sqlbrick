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

// NullInt wraps sql.NullInt.
type NullInt struct {
	sql.NullInt64
}

// NewNullInt will create a new NullInt.
func NewNullInt(i int64) NullInt {
	return NullInt{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: true,
		}}
}
