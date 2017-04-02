package sqlgo

import (
	"database/sql"
	"errors"
	"reflect"
)

// ErrBadArgument indicates that wrong type of wrong value was passed into the function
var ErrBadArgument = errors.New("bad argument")

type Query struct {
	db *sql.DB
}

func With(db *sql.DB) Query {
	return Query{
		db: db,
	}
}

type structMap struct {
	fields  []string
	tags    []string
	methods []string
}

func mapper(val reflect.Value) (structMap, error) {
	var elem reflect.Value
	var typ reflect.Type

	if val.Type().Kind() == reflect.Ptr {
		elem = val.Elem()
		typ = val.Type()
	} else {
		return structMap{}, ErrBadArgument
	}

	m := structMap{}
	m.fields = make([]string, elem.NumField())
	m.tags = make([]string, elem.NumField())
	m.methods = make([]string, val.NumMethod())

	for i := 0; i < elem.NumField(); i++ {
		m.fields[i] = typ.Elem().Field(i).Name
		m.tags[i] = typ.Elem().Field(i).Tag.Get("db")
	}

	for i := 0; i < val.NumMethod(); i++ {
		m.methods[i] = typ.Method(i).Name
	}

	return m, nil
}
