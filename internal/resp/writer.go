package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Writer struct {
	writer *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: bufio.NewWriter(w),
	}
}
func (w *Writer) Flush() error {
	return w.writer.Flush()
}

func (w *Writer) Write(v Value) error {
	var bytes []byte
	switch v.Type {
	case SimpleString:
		bytes = append(bytes, '+')
		bytes = append(bytes, v.Str...)
		bytes = append(bytes, '\r', '\n')

	case Error:
		bytes = append(bytes, '-')
		bytes = append(bytes, v.Str...)
		bytes = append(bytes, '\r', '\n')

	case Integer:
		bytes = append(bytes, ':')
		bytes = append(bytes, strconv.AppendInt(nil, v.Num, 10)...)
		bytes = append(bytes, '\r', '\n')

	case BulkString:
		bytes = append(bytes, '$')
		bytes = append(bytes, strconv.Itoa(len(v.Str))...)
		bytes = append(bytes, '\r', '\n')
		bytes = append(bytes, v.Str...)
		bytes = append(bytes, '\r', '\n')

	case Array:
		bytes = append(bytes, '*')
		bytes = append(bytes, strconv.Itoa(len(v.Elems))...)
		bytes = append(bytes, '\r', '\n')
		if _, err := w.writer.Write(bytes); err != nil {
			return err
		}
		for _, elem := range v.Elems {
			if err := w.Write(elem); err != nil {
				return err
			}
		}
		return nil

	case Null:
		bytes = append(bytes, '$', '-', '1', '\r', '\n')

	default:
		return fmt.Errorf("unknown RESP type: %d", v.Type)
	}
	_, err := w.writer.Write(bytes)
	return err
}
