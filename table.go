package sqlgo

import (
	"database/sql"
	"reflect"
)

type Table map[*sql.ColumnType][]interface{}

func (t Table) Scan(i interface{}) {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		for k, v := range t {
			val := reflect.ValueOf(i).MethodByName("Set" + k.Name())
			if val != reflect.ValueOf(nil) {
				if reflect.ValueOf(v[0]).IsValid() {
					val.Call([]reflect.Value{reflect.ValueOf((interface{})(v[0]))})
				} else {
					val.Call([]reflect.Value{reflect.ValueOf((*interface{})(nil))})
				}
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
