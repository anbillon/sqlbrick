// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"path/filepath"
)

func TestFileGetBrickName(t *testing.T) {
	name := getBrickName(filepath.Join("test", "test.sql"))
	assert.Equal(t, "Test", name)
}

func TestFileGetSourceName(t *testing.T) {
	name := getSourceName(filepath.Join("test", "test.sql"))
	assert.Equal(t, "test", name)
}

func TestFileGetFileName(t *testing.T) {
	name := getFileName(filepath.Join("test", "test.sql"))
	assert.Equal(t, "test.sql", name)
}
