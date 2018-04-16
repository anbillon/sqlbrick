import (
    "database/sql"
    "time"

    "anbillon.com/sqlbrick/typex"
    "github.com/jmoiron/sqlx"
)

// Type definition for {{ .BrickName }} which defined in sql file.
// This can be used as a model in database operation.
type {{ .BrickName }} struct {
    {{ range $k, $v := .Syntaxes }}
    {{ $v.FieldName }} {{ $v.FieldType }}   `db:"{{ $v.DbFieldName }}"`{{ end }}
}

// Type definition for {{ .BrickName }}Brick. This brick will contains all database
// operation from given sql file. Each sql file will have only one brick.
type {{ .BrickName }}Brick struct {
    db *sqlx.DB
}

// new{{ .BrickName }}Brick will create a {{ .BrickName }} brick. This is used
// invoke the query function generated from sql file.
func new{{ .BrickName }}Brick(db *sqlx.DB) *{{ .BrickName }}Brick {
    return &{{ .BrickName }}Brick{db: db}
}
