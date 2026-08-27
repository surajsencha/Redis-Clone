package resp

type ValueType int

const (
	SimpleString ValueType = iota
	Error
	Integer
	BulkString
	Array
	Null
)

type Value struct {
	Type  ValueType
	Str   string
	Num   int64
	List  []string
	Elems []Value
}
