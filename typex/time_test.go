// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package typex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeScanValue(t *testing.T) {
	var nullTime NullTime
	err := nullTime.Scan(time.Now())
	assert.NoError(t, err)
	assert.True(t, nullTime.Valid)

	err = nullTime.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, nullTime.Valid)
	if v, err := nullTime.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong NullTime
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assert.False(t, nullTime.Valid)
}
