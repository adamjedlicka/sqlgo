package sqlgo

import (
	"bytes"
	"reflect"
)

type InsertQuery struct {
	Query
	table   string
	columns []string
	values  []string
	data    []interface{}
}

func (q Query) Insert(into string) InsertQuery {
	return InsertQuery{
		Query: q,
		table: into,
	}
}

func (q InsertQuery) Columns(cols ...string) InsertQuery {
	q.columns = cols
	return q
}

func (q InsertQuery) Values(vals ...string) InsertQuery {
	q.values = make([]string, len(q.columns))

	for i := 0; i < len(q.columns); i++ {
		if i < len(vals) {
			q.values[i] = vals[i]
		} else {
			q.values[i] = "?"
		}
	}

	return q
}

func (q InsertQuery) Data(data ...interface{}) InsertQuery {
	q.data = data
	return q
}

func (q InsertQuery) ToSQL() string {
	var buf bytes.Buffer

	buf.WriteString("INSERT INTO ")
	buf.WriteString(q.table)

	buf.WriteString(" (")
	for i := 0; i < len(q.columns); i++ {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(q.columns[i])
	}
	buf.WriteString(")")

	buf.WriteString(" VALUES")
	for i := 0; i < len(q.data); i++ {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" (")
		for j := 0; j < len(q.values); j++ {
			if j != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(q.values[j])
		}
		buf.WriteString(")")
	}

	return buf.String()
}

func (q InsertQuery) Exec() error {
	flatData := make([]interface{}, 0)

	for _, data := range q.data {
		val := reflect.ValueOf(data)
		m, err := mapper(val)
		if err != nil {
			return err
		}

		for valueIndex, value := range q.values {
			if value == "?" {
				// Check for matching fields
				for i := 0; i < len(m.fields); i++ {
					if m.fields[i] == q.columns[valueIndex] {
						fn := reflect.ValueOf(data).Field(i).MethodByName("Get")
						res := fn.Call([]reflect.Value{})

						flatData = append(flatData, res[0].Interface())

						continue
					}
				}

				// Check for matching tags
				for i := 0; i < len(m.tags); i++ {
					if m.tags[i] == q.columns[valueIndex] {
						fn := reflect.ValueOf(data).Field(i).MethodByName("Get")
						res := fn.Call([]reflect.Value{})

						flatData = append(flatData, res[0].Interface())

						continue
					}
				}

				// Check for matching methods
				for i := 0; i < len(m.methods); i++ {
					if m.methods[i] == "Get"+q.columns[valueIndex] {
						fn := reflect.ValueOf(data).Method(i)
						res := fn.Call([]reflect.Value{})

						flatData = append(flatData, res[0].Interface())

						continue
					}
				}
			}
		}
	}

	_, err := q.db.Exec(q.ToSQL(), flatData...)

	return err
}
