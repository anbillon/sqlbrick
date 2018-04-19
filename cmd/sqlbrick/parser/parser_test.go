// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var parser *Parser

func init() {
	parser = NewParser()
}

func TestParserValidStatement(t *testing.T) {
	tests := []string{
		`{define name Create}
			CREATE TABLE book
		{end define}`,
		`{define name Update}
			update book SET
		{end define}`,
		"SELECT * FROM book",
		"This is non sense",
	}
	for k, v := range tests {
		ok := parser.validStatement(v)
		if k == 3 {
			assert.False(t, ok)
		} else {
			assert.True(t, ok)
		}
	}
}

func TestParserMatchDefineHead(t *testing.T) {
	tests := []string{
		`{define name Create}
			CREATE TABLE book
		{end define}`,
		"This is non sense",
	}

	for k, v := range tests {
		_, ok := parser.matchDefineHead(v)
		if k == 1 {
			assert.False(t, ok)
		} else {
			assert.True(t, ok)
		}
	}
}

func TestParserMatchDefineTail(t *testing.T) {
	tests := []string{
		`{define name Create}
			CREATE TABLE book
		{end define}`,
		"This is non sense",
	}

	for k, v := range tests {
		_, ok := parser.matchDefineTail(v)
		if k == 1 {
			assert.False(t, ok)
		} else {
			assert.True(t, ok)
		}
	}
}

func TestParserParseComment(t *testing.T) {
	tests := []string{
		`-- This is a comment example.`,
		"This is non sense",
	}

	for k, v := range tests {
		parser.currentBlock = v
		c := parser.parseComment()
		if k == 1 {
			assert.Equal(t, "", c)
		} else {
			assert.Equal(t, "This is a comment example.", c)
		}
	}
}

func TestParserParsePlaceholder(t *testing.T) {
	tests := []string{
		`INSERT book(name, price) VALUES(${name}, ${price})`,
		`INSERT book(name, price) VALUES($name}, $price)`,
		"This is non sense",
	}

	for k, v := range tests {
		_, r, err := parser.parsePlaceholder(v)
		if k == 2 {
			assert.Equal(t, "This is non sense", r)
		} else if k == 1 {
			assert.Error(t, err)
		} else {
			assert.Equal(t, "INSERT book(name, price) VALUES(:name, :price)", r)
		}
	}
}

func TestParserParseQueryType(t *testing.T) {
	tests := []string{
		`CREATE TABLE book`,
		`SELECT * FROM book`,
		"This is non sense",
	}

	for k, v := range tests {
		tp := parser.parseQueryType(v)
		if k == 2 {
			assert.Equal(t, QueryTypeInvalid, tp)
		} else if k == 1 {
			assert.Equal(t, QueryTypeSelect, tp)
		} else {
			assert.Equal(t, QueryTypeCreate, tp)
		}
	}
}

func TestParserIsCreateDDL(t *testing.T) {
	tests := []string{
		`CREATE TABLE book`,
		`SELECT * FROM book`,
		"This is non sense",
	}

	for k, v := range tests {
		ok := parser.isCreateDDL(v)
		if k == 0 {
			assert.True(t, ok)
		} else {
			assert.False(t, ok)
		}
	}
}

func TestParserLoadSqlFile(t *testing.T) {
	tests := []string{
		"../../../examples/sql/book.sql",
	}

	for k, v := range tests {
		_, _, err := parser.LoadSqlFile(v)
		if k == 0 {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
