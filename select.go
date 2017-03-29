package sqlgo

type SelectQuery struct {
	Query
	query string
	args  []interface{}
}

func (q Query) Select(from string) SelectQuery {
	return SelectQuery{
		Query: q,
		query: " SELECT * FROM " + from,
		args:  make([]interface{}, 0),
	}
}

func (q SelectQuery) Where(condition string, is interface{}) SelectQuery {
	q.query += " WHERE " + condition
	q.args = append(q.args, is)

	return q
}

func (q SelectQuery) ToSQL() string { return q.query }

func (q SelectQuery) Exec() Table {
	ret := make(Table)

	rows, err := q.db.Query(q.query, q.args...)
	if err != nil {
		return ret
	}
	defer rows.Close()

	columns, err := rows.ColumnTypes()
	if err != nil {
		return ret
	}

	type ref interface{}

	data := make([]interface{}, len(columns))
	dataPtr := make([]interface{}, len(columns))
	for k := range dataPtr {
		dataPtr[k] = &data[k]
	}

	for rows.Next() {
		rows.Scan(dataPtr...)

		for k, v := range columns {
			ret[v] = append(ret[v], *dataPtr[k].(*interface{}))
		}
	}

	return ret
}
