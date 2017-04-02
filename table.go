package sqlgo

import (
	"reflect"

	"database/sql"
)

type Table struct {
	Types []*sql.ColumnType
	Data  [][]interface{}
}

func NewTable() *Table {
	t := new(Table)
	t.Types = make([]*sql.ColumnType, 0)
	t.Data = make([][]interface{}, 0)

	return t
}

func (t *Table) Scan(dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return ErrBadArgument
	}

	val := reflect.ValueOf(dst)
	typ := reflect.TypeOf(dst).Elem()

	if typ.Kind() == reflect.Struct {
		t.scanRow(val, 0)
	} else if typ.Kind() == reflect.Slice {
		for i := 0; i < t.NrRows(); i++ {
			user := reflect.New(typ.Elem())
			t.scanRow(user, i)
			val.Elem().Set(reflect.Append(val.Elem(), user.Elem()))
		}
	}

	return nil
}

func (t *Table) scanRow(val reflect.Value, row int) error {
	elem := val.Elem()

	m, err := mapper(val)
	if err != nil {
		return ErrBadArgument
	}

	for k, v := range t.Types {
		// Check for matching fields
		for i := 0; i < len(m.fields); i++ {
			if m.fields[i] == v.Name() {
				fn := elem.Field(i).Addr().MethodByName("Set")
				if !fn.IsValid() {
					continue
				}

				if reflect.ValueOf(t.Data[row][k]).IsValid() {
					fn.Call([]reflect.Value{reflect.ValueOf(t.Data[row][k])})
				} else {
					fn.Call([]reflect.Value{reflect.ValueOf((*interface{})(nil))})
				}

				continue
			}
		}

		// Check for matching tags
		for i := 0; i < len(m.tags); i++ {
			if m.tags[i] == v.Name() {
				fn := elem.Field(i).Addr().MethodByName("Set")
				if !fn.IsValid() {
					continue
				}

				if reflect.ValueOf(t.Data[row][k]).IsValid() {
					fn.Call([]reflect.Value{reflect.ValueOf(t.Data[row][k])})
				} else {
					fn.Call([]reflect.Value{reflect.ValueOf((*interface{})(nil))})
				}

				continue
			}
		}

		// Check for matching methods
		for i := 0; i < len(m.methods); i++ {
			if m.methods[i] == "Set"+v.Name() {
				if reflect.ValueOf(t.Data[row][k]).IsValid() {
					val.Method(i).Call([]reflect.Value{reflect.ValueOf(t.Data[row][k])})
				} else {
					val.Method(i).Call([]reflect.Value{reflect.ValueOf((*interface{})(nil))})
				}

				continue
			}
		}
	}

	return nil
}

func (t Table) NrRows() int {
	return len(t.Data)
}

func (t Table) NrColumns() int {
	return len(t.Types)
}
