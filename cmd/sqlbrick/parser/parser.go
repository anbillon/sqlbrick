// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

type Parser struct {
	line         string
	currentTag   string
	createDDLTag string
	next         bool
	syntaxes     []Syntax
	keys         []string
	statements   map[string]Statement
}

// NewParser create a new parser to parse sql file.
func NewParser() *Parser {
	return &Parser{
		next:       true,
		statements: make(map[string]Statement),
	}
}

func (p *Parser) getDDLTag(line string) (string, bool) {
	reg := regexp.MustCompile("^\\s*(\\S+) TABLE")
	matches := reg.FindStringSubmatch(strings.ToUpper(line))
	if matches == nil {
		return "", false
	}

	return defaultDDLTag[matches[1]], strings.TrimSpace(matches[1]) == "CREATE"
}

func (p *Parser) getCURDTag(line string) string {
	reg := regexp.MustCompile("^\\s*--\\s*name:\\s*(\\S+)")
	matches := reg.FindStringSubmatch(line)
	if matches == nil {
		return ""
	}

	return matches[1]
}

func (p *Parser) validateQuery() error {
	for _, statement := range p.statements {
		reg := regexp.MustCompile("\\s*(\\S+)\\s*--\\s*name:\\s*(\\S+)")
		matches := reg.FindStringSubmatch(statement.Query)
		if len(matches) >= 2 {
			return errors.New(fmt.Sprintf("validate query error, check if you "+
				"have lost semicolon in your sql query: %v", statement.Query))
		}
	}

	return nil
}

func (p *Parser) appendQueryLine() {
	statement := p.statements
	current := statement[p.currentTag].Query
	line := strings.Trim(p.line, " \t")
	if len(line) == 0 {
		return
	}

	if len(current) > 0 {
		current += "\n"
	}

	current += line

	p.statements[p.currentTag] = Statement{Query: current}
}

func (p *Parser) isDDL(query string) bool {
	reg := regexp.MustCompile("^\\s*(\\S+) TABLE")
	matches := reg.FindStringSubmatch(strings.ToUpper(query))
	if matches != nil {
		return true
	}
	return false
}

// parseQueriesType will parse all query type in parsed queries
func (p *Parser) parseQueriesType(key string, statement Statement) QueryType {
	query := statement.Query
	q := strings.Split(query, " ")
	if q == nil || len(q) < 2 {
		log.Printf("invalid query %v", query)
		return QueryTypeInvalid
	}

	queryType, found := queryTypes[strings.ToUpper(q[0])]
	if !found {
		log.Printf("invalid query %v", query)
		return QueryTypeInvalid
	}

	return queryType
}

// parseWildcards parse queries's wildcards to sql that sqlx can execute.
func (p *Parser) parseWildcards(key string, statement Statement) (string, []string) {
	var args []string
	reg := regexp.MustCompile(`\${(.*?)}`)
	query := statement.Query
	matches := reg.FindAllStringSubmatch(query, -1)
	if matches == nil {
		return query, args
	}

	namedQuery := query
	for _, v := range matches {
		namedQuery = strings.Replace(namedQuery, v[0], fmt.Sprintf(":%s", v[1]), -1)
		args = append(args, strcase.ToLowerCamel(v[1]))
	}

	return namedQuery, args
}

func (p *Parser) parseQueries() {
	for key, value := range p.statements {
		if p.isDDL(value.Query) {
			continue
		}

		queryType := p.parseQueriesType(key, value)
		query, args := p.parseWildcards(key, value)

		p.statements[key] = Statement{
			QueryType: queryType,
			Query:     query,
			Args:      args,
		}
	}
}

func (p *Parser) searchType(source []string) string {
	var i int
	for i = 1; i < len(source); i ++ {
		typeKey := strings.ToLower(strings.TrimSpace(source[i]))
		if len(typeKey) != 0 {
			return typeKey
		}
	}

	return ""
}

func (p *Parser) parseFields() {
	if len(p.createDDLTag) == 0 {
		return
	}

	createDDL := p.statements[p.createDDLTag].Query
	leftIndex := strings.Index(createDDL, "(")
	rightIndex := strings.LastIndex(createDDL, ")")
	if leftIndex <= 0 || rightIndex <= 0 || leftIndex >= rightIndex {
		log.Printf("ddl is not correct: %v", createDDL)
		return
	}
	fieldsSyntax := strings.Split(createDDL[leftIndex+1:rightIndex], ",")
	for _, value := range fieldsSyntax {
		definition := strings.Split(value, " ")
		if definition == nil || len(definition) < 2 {
			log.Printf("invalid defintion: %v", value)
			continue
		}

		typeKey := p.searchType(definition)
		if len(typeKey) == 0 {
			continue
		}

		index := strings.Index(typeKey, "(")
		if index > 0 {
			typeKey = typeKey[0:index]
		}

		var fieldType string
		var found bool
		nullable := !(strings.Contains(strings.ToUpper(value), "NOT") &&
			strings.Contains(strings.ToUpper(value), "NULL"))
		unsigned := strings.Contains(strings.ToUpper(value), "UNSIGNED")
		if nullable {
			fieldType, found = nullableFieldTypes[typeKey]
			if !found || len(fieldType) == 0 {
				log.Printf("invalid type for: %v", value)
				continue
			}
		} else {
			fieldType, found = fieldTypes[typeKey]
			if !found || len(fieldType) == 0 {
				log.Printf("invalid type for: %v", value)
				continue
			}
			if unsigned {
				fieldType = "u" + fieldType
			}
		}

		fieldName := strings.Trim(strings.TrimSpace(definition[0]), `"`)

		syntax := Syntax{
			DbFieldName: fieldName,
			FieldName:   strcase.ToCamel(fieldName),
			FieldType:   fieldType,
		}
		p.syntaxes = append(p.syntaxes, syntax)
	}
}

// LoadSqlFile will load a sql file with given path and split it into
// different part for database usage, such as DDL and CURD.
func (p *Parser) LoadSqlFile(sqlFilePath string) (
	[]string, map[string]Statement, []Syntax, error) {
	f, err := os.Open(sqlFilePath)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p.line = scanner.Text()
		if p.next && strings.TrimSpace(p.line) != "" {
			if tag := p.getCURDTag(p.line); len(tag) > 0 {
				p.currentTag = tag
			} else if tag, isCreateDDL := p.getDDLTag(p.line); len(tag) > 0 {
				p.currentTag = tag
				if isCreateDDL {
					p.createDDLTag = tag
				}
				p.appendQueryLine()
			} else {
				log.Printf("no tag found for current line: %v", p.line)
				continue
			}
			p.keys = append(p.keys, p.currentTag)
			p.next = false
		} else {
			if strings.HasSuffix(p.line, ";") {
				p.line = p.line[:len(p.line)-1]
				p.next = true
			}
			p.appendQueryLine()
		}
	}

	// validate if query was correct
	if err := p.validateQuery(); err != nil {
		return nil, nil, nil, err
	}

	p.parseQueries()

	// parse fields if create ddl tag existed
	p.parseFields()

	return p.keys, p.statements, p.syntaxes, nil
}
