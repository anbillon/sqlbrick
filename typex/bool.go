// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package typex

import (
	"database/sql"
)

// NullString wraps sql.NullBool.
type NullBool struct {
	sql.NullBool
}

// NewNullBool will create a new NullBool.
func NewNullBool(b bool) NullBool {
	return NullBool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: true,
		}}
}
