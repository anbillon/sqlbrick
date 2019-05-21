import (
    "database/sql"
    "errors"
    "time"

    "github.com/anbillon/sqlbrick/typex"
    "github.com/jmoiron/sqlx"
)

// Type definition for {{ .BrickName }} which defined in sql file.
// This can be used as a model in database operation.
type {{ .BrickName }} struct {
    {{- range $k, $v := .Syntaxes }}
    {{ $v.FieldName }} {{ $v.FieldType }}   `db:"{{ $v.DbFieldName }}"`
    {{- end }}
}

// Type definition for {{ .BrickName }}Brick. This brick will contains all database
// operation from given sql file. Each sql file will have only one brick.
type {{ .BrickName }}Brick struct {
    db *sqlx.DB
}

{{ if .HasTx -}}
// Type definition for {{ .BrickName }} transaction. This aims at sql transaction.
type {{ .BrickName }}BrickTx struct {
    tx *sqlx.Tx
}
{{- end }}

// new{{ .BrickName }}Brick will create a {{ .BrickName }} brick. This is used
// invoke the query function generated from sql file.
func new{{ .BrickName }}Brick(db *sqlx.DB) *{{ .BrickName }}Brick {
    return &{{ .BrickName }}Brick{db: db}
}

{{ if .HasTx -}}
// new{{ .BrickName }}Tx will create a new transaction for {{ .BrickName }}.
func (b *{{ .BrickName }}Brick) new{{ .BrickName }}Tx(tx *sqlx.Tx) *{{ .BrickName }}BrickTx {
    return &{{ .BrickName }}BrickTx{tx: tx}
}

// checkTx will check if tx is available.
func (b *{{ .BrickName }}BrickTx) checkTx() error {
    if b.tx == nil {
        return errors.New("the Begin func must be invoked first")
    }
    return nil
}
{{ end }}