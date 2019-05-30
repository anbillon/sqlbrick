// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package internal

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
		if k == 2 || k == 3 {
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

func TestParseDefinition(t *testing.T) {
	tests := []string{
		`{define name Create}
			CREATE TABLE book
		{end define}`,
		`{define name Select, mapper}
			SELECT * FROM book;
		{end define}`,
		`{define name Select, basicType}
			SELECT * FROM book;
		{end define}`,
		`{define Select, basicType}
			SELECT * FROM book;
		{end define}`,
		"This is non sense",
	}

	for k, v := range tests {
		parser.currentBlock = v
		_, err := parser.parseDefinition()
		if k == 0 {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParserParseParamType(t *testing.T) {
	tests := []string{
		"int int",
		"[]int",
		"entity.CustomBook",
	}

	for k, v := range tests {
		parser.currentBlock = "SELECT * FROM book WHERE id = ${id} and uid = ${uid};"
		types, err := parser.parseArgTypes(v)
		if k == 0 {
			assert.Equal(t, 2, len(types))
		} else if k == 1 {
			assert.Error(t, err)
		} else {
			assert.Equal(t, 1, len(types))
		}
	}
}

func TestParseTypeDetail(t *testing.T) {
	tests := []string{
		"int",
		"entity.Book",
		"[]string",
	}

	for k, v := range tests {
		detail, _ := parser.parseTypeDetail(v)
		if k == 0 {
			assert.Equal(t, DataPrimitive, detail.Type)
		} else if k == 1 {
			assert.Equal(t, DataStruct, detail.Type)
		} else {
			assert.Equal(t, DataArray, detail.Type)
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
			assert.Equal(t, "this is a comment example.", c)
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
		"../examples/sqb/book.sqb",
	}

	for k, v := range tests {
		_, _, _, err := parser.LoadSqbFile(v)
		if k == 0 {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
