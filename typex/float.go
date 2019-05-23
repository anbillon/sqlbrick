// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package typex

import (
	"database/sql"
)

// NullFloat wraps sql.NullFloat64.
type NullFloat struct {
	sql.NullFloat64
}

// NewNullFloat will create a new NullFloat.
func NewNullFloat(f float64) NullFloat {
	return NullFloat{
		NullFloat64: sql.NullFloat64{
			Float64: f,
			Valid:   true,
		}}
}
