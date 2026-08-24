package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r)}
}

func readSimpleString(reader *Reader) (Value, error) {
	str, err := reader.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	clearText := strings.TrimSuffix(str, "\r\n")

	return Value{Type: SimpleString, Str: clearText}, nil
}

func readError(reader *Reader) (Value, error) {
	str, err := readSimpleString(reader)
	if err != nil {
		return Value{}, err
	}
	return Value{Type: Error, Str: str.Str}, nil
}

func readInteger(reader *Reader) (Value, error) {
	str, err := readSimpleString(reader)
	if err != nil {
		return Value{}, err
	}
	num, err := strconv.ParseInt(str.Str, 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("invalid integer: %s", err)
	}
	return Value{Type: Integer, Num: num}, nil

}
func readBulkString(reader *Reader) (Value, error) {
	str, err := readSimpleString(reader)
	if err != nil {
		return Value{}, err
	}
	length, err := strconv.Atoi(str.Str)
	if err != nil {
		return Value{}, err
	}
	if length == -1 {
		return Value{Type: Null}, nil
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(reader.reader, data); err != nil {
		return Value{}, err
	}
	if _, err := reader.reader.ReadBytes('\n'); err != nil {
		return Value{}, err
	}
	return Value{Type: BulkString, Str: string(data)}, nil
}

func readArray(reader *Reader) (Value, error) {
	str, err := readSimpleString(reader)
	if err != nil {
		return Value{}, err
	}
	count, _ := strconv.Atoi(str.Str)
	if count == -1 {
		return Value{Type: Null}, nil
	}
	elems := make([]Value, count)
	for i := 0; i < count; i++ {
		val, err := Read(reader)
		if err != nil {
			return Value{}, err
		}
		elems[i] = val
	}
	return Value{Type: Array, Elems: elems}, nil
}

func Read(reader *Reader) (Value, error) {
	firstByte, err := reader.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch firstByte {
	case '+':
		return readSimpleString(reader)
	case '-':
		return readError(reader)
	case ':':
		return readInteger(reader)
	case '$':
		return readBulkString(reader)
	case '*':
		return readArray(reader)
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %c", firstByte)
	}
}
