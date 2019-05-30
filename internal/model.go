// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package internal

var fieldTypes = map[string]string{
	"bool":        "bool",
	"boolean":     "bool",
	"tinyint":     "int8",
	"smallint":    "int16",
	"integer":     "int32",
	"int":         "int",
	"bigint":      "int64",
	"int2":        "int16",
	"int4":        "int32",
	"int8":        "int64",
	"smallserial": "int16",
	"serial":      "int32",
	"bigserial":   "int64",
	"serial2":     "int16",
	"serial4":     "int32",
	"serial8":     "int64",
	"float":       "float32",
	"real":        "float32",
	"double":      "float64",
	"numeric":     "float64",
	"decimal":     "float64",
	"float4":      "float32",
	"float8":      "float64",
	"text":        "string",
	"varchar":     "string",
	"char":        "string",
	"bit":         "string",
	"varbit":      "string",
	"interval":    "string",
	"money":       "string",
	"cidr":        "string",
	"inet":        "string",
	"macaddr":     "string",
	"uuid":        "string",
	"json":        "typex.JsonText",
	"xml":         "string",
	"year":        "string",
	"date":        "time.Time",
	"datetime":    "time.Time",
	"timestamp":   "time.Time",
	"time":        "time.Time",
	"timetz":      "time.Time",
	"timestamptz": "time.Time",
	"binary":      "[]byte",
	"varbinary":   "[]byte",
	"tinyblob":    "[]byte",
	"blob":        "[]byte",
	"mediumblob":  "[]byte",
	"longblob":    "[]byte",
	"bytea":       "[]byte",
}

var nullableFieldTypes = map[string]string{
	"bool":        "typex.NullBool",
	"boolean":     "typex.NullBool",
	"tinyint":     "typex.NullInt",
	"smallint":    "typex.NullInt",
	"integer":     "typex.NullInt",
	"int":         "typex.NullInt",
	"bigint":      "typex.NullInt",
	"int2":        "typex.NullInt",
	"int4":        "typex.NullInt",
	"int8":        "typex.NullInt",
	"smallserial": "typex.NullInt",
	"serial":      "typex.NullInt",
	"bigserial":   "typex.NullInt",
	"serial2":     "typex.NullInt",
	"serial4":     "typex.NullInt",
	"serial8":     "typex.NullInt",
	"float":       "typex.NullFloat",
	"double":      "typex.NullFloat",
	"real":        "typex.NullFloat",
	"numeric":     "typex.NullFloat",
	"decimal":     "typex.NullFloat",
	"float4":      "typex.NullFloat",
	"float8":      "typex.NullFloat",
	"text":        "typex.NullString",
	"varchar":     "typex.NullString",
	"char":        "typex.NullString",
	"bit":         "typex.NullString",
	"varbit":      "typex.NullString",
	"interval":    "typex.NullString",
	"money":       "typex.NullString",
	"cidr":        "typex.NullString",
	"inet":        "typex.NullString",
	"macaddr":     "typex.NullString",
	"uuid":        "typex.NullString",
	"json":        "typex.NullJsonText",
	"xml":         "typex.NullString",
	"year":        "typex.NullString",
	"date":        "typex.NullTime",
	"datetime":    "typex.NullTime",
	"timestamp":   "typex.NullTime",
	"time":        "typex.NullTime",
	"timetz":      "typex.NullTime",
	"timestamptz": "typex.NullTime",
	"binary":      "typex.NullBytes",
	"varbinary":   "typex.NullBytes",
	"tinyblob":    "typex.NullBytes",
	"blob":        "typex.NullBytes",
	"mediumblob":  "typex.NullBytes",
	"longblob":    "typex.NullBytes",
	"bytea":       "typex.NullBytes",
}

type QueryType int8

const (
	QueryTypeInvalid QueryType = iota
	QueryTypeCreate
	QueryTypeInsert
	QueryTypeDelete
	QueryTypeUpdate
	QueryTypeSelect
)

var queryTypes = map[string]QueryType{
	"CREATE": QueryTypeCreate,
	"INSERT": QueryTypeInsert,
	"DELETE": QueryTypeDelete,
	"UPDATE": QueryTypeUpdate,
	"SELECT": QueryTypeSelect,
}

type DataType int8

const (
	DataPrimitive DataType = iota
	DataStruct
	DataArray
)

type Syntax struct {
	DbFieldName string
	FieldName   string
	FieldType   string
}

type Condition struct {
	Expression string
	Query      string
	AppendNext bool
}

type DynamicQuery struct {
	QueryType       QueryType
	Conditions      []Condition
	Segments        []string
	Args            []string
	IndexOfWhere    int
	RemoveLastComma bool
}

type TypeDetail struct {
	Name string
	Type DataType
}

type Definition struct {
	Name     string
	ArgTypes []TypeDetail
	Mapper   *TypeDetail
	IsTx     bool
}

type Statement struct {
	Definition *Definition
	Query      *DynamicQuery
	Comment    string
}
