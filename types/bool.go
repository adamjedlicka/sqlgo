package types

import "strconv"

type Bool struct {
	notNull bool
	isNull  bool
	value   bool
}

func NewBool(notNull bool, value ...bool) Bool {
	b := Bool{}
	b.notNull = notNull
	b.isNull = true

	if len(value) > 0 {
		b.isNull = false
		b.value = value[0]
	}

	return b
}

func (b Bool) NotNull() bool { return b.notNull }

func (b Bool) Value() bool { return b.value }

func (b Bool) Get() interface{} {
	if b.isNull && !b.notNull {
		return nil
	}

	return b.value
}

func (b *Bool) Set(val interface{}) {
	switch val.(type) {
	case *interface{}:
		b.isNull = true
		b.value = false
	case int64:
		b.isNull = false
		b.value = val.(int64) == 1
	case []uint8:
		tmp, err := strconv.ParseInt(string(val.([]uint8)), 10, 64)
		if err == nil {
			b.isNull = false
			b.value = tmp == 1
		}
	default:
		// TODO
	}
}

func (b Bool) String() string {
	if b.isNull {
		return "NULL"
	} else if b.value == true {
		return "TRUE"
	} else {
		return "FALSE"
	}
}
