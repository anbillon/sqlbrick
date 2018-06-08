// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

type Parser struct {
	line         string
	currentBlock string
	next         bool
	newBlock     bool
	sqlBlocks    []string
	syntaxes     []Syntax
	definitions  []Statement
}

// NewParser create a new parser to parse sql file.
func NewParser() *Parser {
	return &Parser{
		next: true,
	}
}

func (p *Parser) validStatement(block string) bool {
	_, headOk := p.matchDefineHead(block)
	_, tailOk := p.matchDefineTail(block)
	if headOk && tailOk {
		return true
	}

	reg := regexp.MustCompile(`TABLE|INSERT|DELETE|UPDATE|SELECT`)
	m := reg.FindStringSubmatch(strings.ToUpper(block))
	if m != nil {
		return true
	}

	return false
}

func (p *Parser) appendLine() error {
	currentIndex := len(p.sqlBlocks) - 1
	current := p.line + "\n"

	if p.newBlock {
		p.sqlBlocks = append(p.sqlBlocks, current)
	} else {
		p.sqlBlocks[currentIndex] += current
	}

	if _, ok := p.matchDefineHead(current); ok && currentIndex > 0 {
		if !p.validStatement(p.sqlBlocks[currentIndex]) {
			return errors.Errorf(
				"error definition found for new block:\n%v",
				p.sqlBlocks[currentIndex])
		}
	}

	return nil
}

// matchDefineHead will match {define ...} header
func (p *Parser) matchDefineHead(line string) ([]string, bool) {
	reg := regexp.MustCompile(`{\s*define(.*?)}`)
	matches := reg.FindStringSubmatch(line)

	if matches == nil || len(matches) == 0 {
		return nil, false
	}

	return matches, true
}

// matchDefineTail will match {end define} tail
func (p *Parser) matchDefineTail(line string) ([]string, bool) {
	reg := regexp.MustCompile(`{\s*end \s*define\s*}`)
	matches := reg.FindStringSubmatch(line)

	if matches == nil || len(matches) == 0 {
		return nil, false
	}

	return matches, true
}

func (p *Parser) searchFieldAndType(source []string) (string, string) {
	reg := regexp.MustCompile(`\([0-9]*\)|"`)
	fieldName := ""
	for i := 0; i < len(source); i++ {
		input := strings.ToLower(strings.TrimSpace(source[i]))
		if len(input) == 0 {
			continue
		}

		if len(fieldName) > 0 {
			return fieldName, reg.ReplaceAllString(input, "")
		} else {
			fieldName = reg.ReplaceAllString(input, "")
		}
	}

	return "", ""
}

// parseFields will parse all fields from create table statement.
func (p *Parser) parseFields(block string) error {
	leftIndex := strings.Index(block, "(")
	rightIndex := strings.LastIndex(block, ")")
	if leftIndex <= 0 || rightIndex <= 0 || leftIndex >= rightIndex {
		return errors.Errorf("create ddl is not correct: \n%v", block)
	}

	specialRegex := regexp.MustCompile(`^\s*(UNIQUE|CONSTRAINT|PRIMARY KEY|FOREIGN KEY|REFERENCES|INDEX|CHECK) \s*(\S+)`)
	fieldsSyntax := strings.Split(block[leftIndex+1:rightIndex], ",")
	for _, value := range fieldsSyntax {
		if specialRegex.FindStringSubmatch(strings.TrimSpace(value)) != nil {
			continue
		}

		definition := strings.Split(value, " ")
		if definition == nil || len(definition) == 0 {
			return errors.Errorf("invalid defintion: %v", value)
		}

		fieldName, typeKey := p.searchFieldAndType(definition)
		if len(typeKey) == 0 {
			continue
		}

		var fieldType string
		var found bool
		nullable := !(strings.Contains(strings.ToUpper(value), "NOT") &&
			strings.Contains(strings.ToUpper(value), "NULL"))
		unsigned := strings.Contains(strings.ToUpper(value), "UNSIGNED")
		if nullable {
			fieldType, found = nullableFieldTypes[typeKey]
			if !found || len(fieldType) == 0 {
				return errors.Errorf("invalid type for: %v", value)
			}
		} else {
			fieldType, found = fieldTypes[typeKey]
			if !found || len(fieldType) == 0 {
				return errors.Errorf("invalid type for: %v", value)
			}
			if unsigned {
				fieldType = "u" + fieldType
			}
		}

		syntax := Syntax{
			DbFieldName: fieldName,
			FieldName:   strcase.ToCamel(fieldName),
			FieldType:   fieldType,
		}
		p.syntaxes = append(p.syntaxes, syntax)
	}

	return nil
}

// parseComment will parse all comment for sql statement
func (p *Parser) parseComment() string {
	reg := regexp.MustCompile(`--(.*)`)
	matches := reg.FindAllStringSubmatch(p.currentBlock, -1)
	if matches == nil || len(matches) == 0 {
		return ""
	}

	var comment string
	// remove comsumed line
	for i := 0; i < len(matches); i++ {
		p.currentBlock = strings.TrimSpace(
			strings.Replace(p.currentBlock, matches[i][0], "", 1))
		comment += matches[i][1]
	}

	return strings.TrimSpace(comment)
}

// parseDynamicQueries will parse definition for sql statement
func (p *Parser) parseDefinition() (*Definition, error) {
	nameRegex := regexp.MustCompile(`name (.*)`)
	mapperRegex := regexp.MustCompile(`mapper (.*)`)
	txRegex := regexp.MustCompile(`tx (true|false)`)
	var definition Definition
	block := p.currentBlock

	headMatches, _ := p.matchDefineHead(block)
	tailMatches, _ := p.matchDefineTail(block)
	if headMatches == nil && len(headMatches) < 2 ||
		tailMatches == nil || len(tailMatches) == 0 {
		return nil, errors.Errorf("error definition:\n %v", block)
	}

	// remove comsumed line
	p.currentBlock = strings.TrimSpace(
		strings.Replace(p.currentBlock, headMatches[0], "", 1))
	p.currentBlock = strings.TrimSpace(
		strings.Replace(p.currentBlock, tailMatches[0], "", 1))
	definitionHead := strings.TrimSpace(headMatches[1])
	left := definitionHead
	tags := strings.Split(definitionHead, ",")
	for _, value := range tags {
		nm := nameRegex.FindStringSubmatch(value)
		mm := mapperRegex.FindStringSubmatch(value)
		tm := txRegex.FindStringSubmatch(value)
		if nm != nil && mm != nil {
			return nil, errors.Errorf("definition format error:\n%v ", block)
		}
		if nm != nil && len(nm) >= 2 {
			left = strings.Replace(left, value, "", 1)
			definition.Name = strings.TrimSpace(nm[1])
		}
		if mm != nil && len(mm) >= 2 {
			left = strings.Replace(left, value, "", 1)
			definition.Mapper = mappers[strings.TrimSpace(mm[1])]
			if definition.Mapper == MapperDefault {
				return nil, errors.Errorf("error mapper type %v\n", block)
			}
		}
		if tm != nil && len(tm) >= 1 {
			left = strings.Replace(left, value, "", 1)
			definition.IsTx = true
		}
	}

	if len(strings.TrimSpace(strings.Replace(left, ",", "", -1))) != 0 {
		return nil, errors.Errorf("definition format error:\n%v ", block)
	}

	if len(definition.Name) == 0 {
		return nil, errors.Errorf("no name definition found:\n %v", block)
	}

	if definition.Mapper == MapperDefault {
		definition.Mapper = MapperArray
	}

	return &definition, nil
}

// parsePlaceholder parse queries's placeholder to sql that sqlx can execute.
func (p *Parser) parsePlaceholder(block string) ([]string, string, error) {
	reg := regexp.MustCompile(`\${(.*?)}`)
	errorReg := regexp.MustCompile(`\$(.*?)`)
	matches := reg.FindAllStringSubmatch(block, -1)
	if matches == nil {
		if errorReg.FindStringSubmatch(block) != nil {
			return nil, block, errors.Errorf("invalid placeholder:\n %v", block)
		}
		return nil, block, nil
	}

	var args []string
	for _, v := range matches {
		block = strings.Replace(
			block, v[0], fmt.Sprintf(":%s", v[1]), -1)
		args = append(args, strcase.ToLowerCamel(v[1]))
	}

	return args, block, nil
}

func (p *Parser) convertExpression(expression string, fieldName string) string {
	if len(fieldName) == 0 {
		return expression
	}

	var fieldType string
	for _, v := range p.syntaxes {
		if strings.ToUpper(v.FieldName) == strings.ToUpper(strcase.ToCamel(fieldName)) {
			fieldType = v.FieldType
			break
		}
	}

	var extra string
	switch fieldType {
	case "typex.NullBool":
		extra = ".Bool"
	case "typex.NullInt":
		extra = ".Int64"
	case "typex.NullFloat":
		extra = ".Float64"
	case "typex.NullString":
		extra = ".String"
	case "typex.NullTime":
		extra = ".Time"
	case "typex.NullBytes":
		extra = ".Bytes"
	}

	return strings.Replace(expression, fieldName,
		strcase.ToCamel(fieldName)+extra, -1)
}

func (p *Parser) parseDynamicQuery() (*DynamicQuery, error) {
	statement := strings.Replace(p.currentBlock, ";", "", -1)
	var dynamicQuery = new(DynamicQuery)
	splitRegexp := regexp.MustCompile(`({\s*if (.*?)})[\s\S]*?({\s*end \s*if\s*})`)
	headRegexp := regexp.MustCompile(`({\s*if (.*?)})`)
	endRegexp := regexp.MustCompile(`({\s*end \s*if\s*})`)
	fieldRegexp := regexp.MustCompile(`:([A-Za-z0-9_-]*)`)

	args, consumedQuery, err := p.parsePlaceholder(statement)
	if err != nil {
		return nil, err
	}
	matches := splitRegexp.FindAllStringSubmatch(consumedQuery, -1)
	// no condition found
	if matches == nil || len(matches) == 0 {
		dynamicQuery.Segments = append(dynamicQuery.Segments, consumedQuery)
		dynamicQuery.QueryType = p.parseQueryType(statement)
		dynamicQuery.Args = args
		return dynamicQuery, nil
	}

	queryType := p.parseQueryType(statement)
	if queryType == QueryTypeInvalid {
		return nil, errors.Errorf("invalid query %v", statement)
	}

	// TODO: add select support
	if queryType != QueryTypeUpdate {
		return nil, errors.Errorf(
			"condition statement only support update now:\n%v", statement)
	}

	var conditions []Condition
	for index, value := range matches {
		query := value[0]
		// append next condition or not when generate souce code
		appendNext := false

		// to check if next was condition or segment
		i := strings.Index(consumedQuery, query)
		if i >= 0 {
			left := strings.TrimSpace(consumedQuery[i+len(query):])
			lm := headRegexp.FindStringSubmatch(left)
			if lm != nil && len(lm) >= 2 && strings.HasPrefix(left, lm[0]) {
				appendNext = true
			}
		}

		// remove condition head and end
		query = headRegexp.ReplaceAllString(query, "")
		query = strings.TrimSpace(endRegexp.ReplaceAllString(query, ""))
		if index < len(matches)-1 && len(matches) != 1 && !strings.HasSuffix(query, ",") {
			return nil, errors.Errorf("invalid statement, missing comma: %v", query)
		}

		hm := headRegexp.FindStringSubmatch(value[0])
		if hm == nil || len(hm) != 3 {
			return nil, errors.Errorf("invalid condition: %v", statement)
		}

		fm := fieldRegexp.FindStringSubmatch(query)
		if fm == nil {
			return nil, errors.Errorf("invalid condition: %v", statement)
		}
		fieldName := fm[1]
		conditions = append(conditions, Condition{
			Expression: p.convertExpression(hm[2], fieldName),
			Query:      query,
			AppendNext: appendNext,
		})
	}

	segments := splitRegexp.Split(consumedQuery, -1)
	var realSegments []string
	var removeLastComma bool
	var indexOfWhere int
	var spaceNumber int
	for index, v := range segments {
		segment := strings.Replace(strings.TrimSpace(v), "\n", " ", -1)
		if len(segment) == 0 {
			spaceNumber += 1
			continue
		}
		realSegments = append(realSegments, segment)
		if strings.Contains(segment, "WHERE") {
			indexOfWhere = index - spaceNumber
			// to avoid previous 'WHERE' prefix check
			if removeLastComma {
				removeLastComma = false
			}
		}
		if strings.HasPrefix(segment, "WHERE") {
			removeLastComma = true
		}
	}
	if len(conditions) < len(realSegments)-1 {
		return nil, errors.Errorf("invalid condition: %v", statement)
	}

	dynamicQuery.QueryType = queryType
	dynamicQuery.Args = args
	dynamicQuery.Conditions = conditions
	dynamicQuery.Segments = realSegments
	dynamicQuery.IndexOfWhere = indexOfWhere
	dynamicQuery.RemoveLastComma = removeLastComma

	return dynamicQuery, nil
}

// parseQueriesType will parse all query type in parsed queries
func (p *Parser) parseQueryType(block string) QueryType {
	q := strings.Split(block, " ")
	if q == nil || len(q) < 2 {
		return QueryTypeInvalid
	}

	queryType, found := queryTypes[strings.ToUpper(q[0])]
	if !found {
		return QueryTypeInvalid
	}

	return queryType
}

func (p *Parser) parseSqlBlocks() error {
	// parse field first
	createDDLIndex := -1
	for index, block := range p.sqlBlocks {
		if p.isCreateDDL(block) {
			createDDLIndex = index
			break
		}
	}

	if createDDLIndex == -1 {
		return errors.New("create table statement not found")
	}

	err := p.parseFields(p.sqlBlocks[createDDLIndex])
	if err != nil {
		return err
	}

	for _, block := range p.sqlBlocks {
		p.currentBlock = block
		comment := p.parseComment()
		definition, err := p.parseDefinition()
		if err != nil {
			return err
		}
		dymanicQuery, err := p.parseDynamicQuery()
		if err != nil {
			return err
		}

		p.definitions = append(p.definitions, Statement{
			Definition: definition,
			Query:      dymanicQuery,
			Comment:    comment,
		})
	}

	return nil
}

func (p *Parser) isCreateDDL(block string) bool {
	reg := regexp.MustCompile(`\s*(\S+) TABLE`)
	matches := reg.FindStringSubmatch(strings.ToUpper(block))
	if matches == nil || len(matches) < 2 {
		return false
	}

	return strings.TrimSpace(matches[1]) == "CREATE"
}

// LoadSqlFile will load a sql file with given path and split it into
// different part for database usage, such as DDL and CURD.
func (p *Parser) LoadSqlFile(sqlFilePath string) ([]Statement, []Syntax, error) {
	f, err := os.Open(sqlFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p.line = scanner.Text()
		if len(strings.TrimSpace(p.line)) == 0 {
			continue
		}

		if p.next {
			p.next = false
			p.newBlock = true
		} else {
			p.newBlock = false
			if _, ok := p.matchDefineTail(p.line); ok {
				p.next = true
			}
		}

		if err := p.appendLine(); err != nil {
			return nil, nil, err
		}
	}

	// walk all blocks to parse
	if err := p.parseSqlBlocks(); err != nil {
		return nil, nil, err
	}

	return p.definitions, p.syntaxes, nil
}
