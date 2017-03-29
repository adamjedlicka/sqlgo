package sqlgo

import "database/sql"

type Query struct {
	db *sql.DB
}

func Use(db *sql.DB) Query {
	return Query{
		db: db,
	}
}
