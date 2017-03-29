package sqlgo

import (
	"database/sql"
	"reflect"
)

type Table map[*sql.ColumnType][]interface{}

func (t Table) Scan(i interface{}) {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		var fn reflect.Value

		for k, v := range t {

			field := reflect.ValueOf(i).Elem().FieldByName(k.Name())
			if field.IsValid() {
				fn = field.Addr().MethodByName("Set")

				if reflect.ValueOf(v[0]).IsValid() {
					fn.Call([]reflect.Value{reflect.ValueOf((interface{})(v[0]))})
				} else {
					fn.Call([]reflect.Value{reflect.ValueOf((*interface{})(nil))})
				}

				continue
			}

			fn = reflect.ValueOf(i).MethodByName("Set" + k.Name())
			if fn != reflect.ValueOf(nil) {
				if reflect.ValueOf(v[0]).IsValid() {
					fn.Call([]reflect.Value{reflect.ValueOf((interface{})(v[0]))})
				} else {
					fn.Call([]reflect.Value{reflect.ValueOf((*interface{})(nil))})
				}

				continue
			}

		}
	} else if reflect.TypeOf(i).Kind() == reflect.Slice {

	}
}

func (t Table) NrRows() int {
	for k := range t {
		return len(t[k])
	}

	return 0
}

func (t Table) NrColumns() int {
	return len(t)
}
