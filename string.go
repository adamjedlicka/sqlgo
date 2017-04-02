package sqlgo

type String struct {
	notNull bool
	isNull  bool
	value   string
}

func NewString(notNull bool, value ...string) String {
	s := String{}
	s.notNull = notNull
	s.isNull = true

	if len(value) > 0 {
		s.isNull = false
		s.value = value[0]
	}

	return s
}

func (s String) NotNull() bool { return s.notNull }

func (s String) Value() string { return s.value }

func (s String) Get() interface{} {
	if s.isNull && !s.notNull {
		return nil
	}

	return s.value
}

func (s *String) Set(val interface{}) {
	switch val.(type) {
	case *interface{}:
		s.isNull = true
		s.value = ""
	case string:
		s.isNull = false
		s.value = val.(string)
	case []uint8:
		s.isNull = false
		s.value = string(val.([]uint8))
	default:
		// TODO
	}
}

func (s String) String() string {
	if s.isNull {
		return "NULL"
	}
	return s.value
}
