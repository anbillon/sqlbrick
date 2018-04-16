import (
    "github.com/jmoiron/sqlx"
)

// Type definition for SqlBrick. It contains all bricks depends on the number of
// sql files. It alson wrap some sqlx func for for convenience.
type SqlBrick struct {
    Db *sqlx.DB{{ range $k, $v := .Bricks }}
    {{ $v }} *{{ $v }}Brick{{ end }}
}

// NewSqlBrick create a new SqlBrick to operate all bricks.
func NewSqlBrick(db *sqlx.DB) *SqlBrick {
    return &SqlBrick{
        Db: db,{{ range $k, $v := .Bricks }}
        {{ $v }}: new{{ $v }}Brick(db),{{ end }}
    }
}

// Begin wraps the transaction of sqlx.Beginx(). It returns sqlx.Tx to
// operate transaction as sqlx do.
func (s *SqlBrick) Begin() (*sqlx.Tx, error) {
    return s.Db.Beginx()
}
