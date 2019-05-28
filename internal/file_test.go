// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package internal

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileGetBrickName(t *testing.T) {
	name := GetBrickName(filepath.Join("test", "test.sqb"))
	assert.Equal(t, "Test", name)
}

func TestFileGetSourceName(t *testing.T) {
	name := GetSourceName(filepath.Join("test", "test.sqb"))
	assert.Equal(t, "test", name)
}

func TestFileGetFileName(t *testing.T) {
	name := GetFileName(filepath.Join("test", "test.sqb"))
	assert.Equal(t, "test.sqb", name)
}
