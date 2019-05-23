module github.com/anbillon/sqlbrick

go 1.12

require (
	github.com/anbillon/sqlbrick/typex v0.0.0-00010101000000-000000000000
	github.com/gobuffalo/packd v0.1.0 // indirect
	github.com/gobuffalo/packr v1.25.0
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.1.1
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	golang.org/x/tools v0.0.0-20190521203540-521d6ed310dd
	google.golang.org/appengine v1.6.0 // indirect
)

replace github.com/anbillon/sqlbrick/typex => ./typex
