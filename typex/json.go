// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.
//
// Package typex implements the Scanner interface and the driver
// Valuer interface so we can use them in database operation.

package typex

import (
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
)

// JsonText is a types.JSONText, just for alias
type JsonText types.JSONText

// NullJsonText is a types.NullJSONText just for alias
type NullJsonText types.NullJSONText

// NewNullJsonText will create a new NullJsonText.
func NewNullJsonText(m json.RawMessage) NullJsonText {
	return NullJsonText{
		JSONText: types.JSONText(m),
		Valid:    m != nil,
	}
}
