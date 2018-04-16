// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package typex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesScanValue(t *testing.T) {
	var nullBytes NullBytes
	err := nullBytes.Scan([]byte("Test bytes"))
	assert.NoError(t, err)
	assert.True(t, nullBytes.Valid, "assert error")

	err = nullBytes.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, nullBytes.Valid)
}
