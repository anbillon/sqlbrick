// Copyright (c) 2019 Anbillon Team (anbillonteam@gmail.com).

package entity

import (
	"github.com/anbillon/sqlbrick/typex"
)

type CustomBook struct {
	Id      int32            `db:"id"`
	Uid     int32            `db:"uid"`
	Name    string           `db:"name"`
	Content typex.NullString `db:"content"`
}
