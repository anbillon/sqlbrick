import (
    "time"

    "github.com/anbillon/sqlbrick/typex"
)

// Type definition for {{ .BrickName }} which defined in sqb file.
// This can be used as a model in database operation.
type {{ .BrickName }} struct {
    {{- range $k, $v := .Syntaxes }}
    {{ $v.FieldName }} {{ $v.FieldType }}   `db:"{{ $v.DbFieldName }}"`
    {{- end }}
}