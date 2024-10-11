package enum

var id2Enum = make(map[any]map[int]any)
var str2Enum = make(map[any]map[string]any)

type Enum[T any] struct {
	id  int
	str string
}

func Default[T any]() Enum[T] {
	return Enum[T]{}
}

func FromID[T any](id int) Enum[T] {
	var t T
	return id2Enum[t][id].(Enum[T])
}

func FromStr[T any](str string) Enum[T] {
	var t T
	return str2Enum[t][str].(Enum[T])
}

func New[T any](id int, val string) Enum[T] {
	if id == 0 {
		panic("not accept zero enum")
	}

	var t T
	if _, ok := id2Enum[t]; !ok {
		id2Enum[t] = make(map[int]any)
		str2Enum[t] = make(map[string]any)
	}

	enum := Enum[T]{id, val}
	id2Enum[t][id] = enum
	str2Enum[t][val] = enum

	return enum
}

func (e Enum[T]) ID() int {
	return e.id
}

func (e Enum[T]) String() string {
	return e.str
}
