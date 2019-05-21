{{- $txLen := len .TxMap }}
import (
    "errors"

    "github.com/jmoiron/sqlx"
)

// Type definition for SqlBrick. It contains all bricks depend on the number of
// sql files. It also wraps some sqlx func for for convenience.
type SqlBrick struct {
    Db *sqlx.DB
    {{- range $k, $v := .Bricks }}
    {{ $v }} *{{ $v }}Brick
    {{- end }}
}

{{ if $txLen -}}
// Type definition for brick transaction. This aims at sql transaction.
// If you want a transaction, then invoke 'Begin' to get this struct.
type BrickTx struct {
    tx *sqlx.Tx
    {{- range $k, $v := .TxMap }}
    {{- if $v }}
    {{ $k }} *{{ $k }}BrickTx
    {{- end }}
    {{- end }}
}
{{- end }}

// NewSqlBrick create a new SqlBrick to operate all bricks.
func NewSqlBrick(db *sqlx.DB) *SqlBrick {
    return &SqlBrick{
        Db: db,
        {{- range $k, $v := .Bricks }}
        {{ $v }}: new{{ $v }}Brick(db),
        {{- end }}
    }
}

{{ if $txLen -}}
// Begin will start a new transaction for bricks. If any query
// is defined as tx sql, this must be invoked.
func (b *SqlBrick) Begin() (*BrickTx, error) {
    tx, err := b.Db.Beginx()
    if err != nil {
        return nil, err
    }

    return &BrickTx{
        tx: tx,
        {{- range $k, $v := .TxMap }}
        {{- if $v }}
        {{ $k }}: b.{{ $k }}.new{{ $k }}Tx(tx),
        {{- end }}
        {{- end }}
    }, nil
}

// Commit will end a transaction for brick. Begin must be invoked
// before Commit. Otherwise there will be an error.
func (b *BrickTx) Commit() error {
    if b.tx == nil {
        return errors.New("the Begin func must be invoked first")
    }

    return b.tx.Commit()
}
{{- end }}
