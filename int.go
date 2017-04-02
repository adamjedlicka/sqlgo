package sqlgo

import "strconv"

type Int struct {
	notNull bool
	isNull  bool
	value   int64
}

func NewInt(notNull bool, value ...int64) Int {
	i := Int{}
	i.notNull = notNull
	i.isNull = true

	if len(value) > 0 {
		i.isNull = false
		i.value = value[0]
	}

	return i
}

func (i Int) NotNull() bool { return i.notNull }

func (i Int) Value() int64 { return i.value }

func (i Int) Get() interface{} {
	if i.isNull && !i.notNull {
		return nil
	}

	return i.value
}

func (i *Int) Set(val interface{}) {
	switch val.(type) {
	case *interface{}:
		i.isNull = true
		i.value = 0
	case int64:
		i.isNull = false
		i.value = val.(int64)
	case []uint8:
		tmp, err := strconv.ParseInt(string(val.([]uint8)), 10, 64)
		if err == nil {
			i.isNull = false
			i.value = tmp
		}
	default:
		// TODO
	}
}

func (i Int) String() string {
	if i.isNull {
		return "NULL"
	}
	return strconv.FormatInt(i.value, 10)
}
