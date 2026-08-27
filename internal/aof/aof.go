package aof

import (
	"io"
	"os"
	"sync"

	"github.com/surajsencha/redis-clone/internal/resp"
)

type AOF struct {
	mu   sync.Mutex
	file *os.File
}

func (s *AOF) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.file.Close()
}

func NewAOF(path string) (*AOF, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &AOF{
		file: file,
	}, nil
}

func (a *AOF) Write(value resp.Value) error {

	a.mu.Lock()
	defer a.mu.Unlock()
	w := resp.NewWriter(a.file)
	w.Write(value)
	w.Flush()
	a.file.Sync()
	return nil

}

func (a *AOF) Read(callback func(value resp.Value)) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.file.Seek(0, io.SeekStart)
	r := resp.NewReader(a.file)
	for {
		value, err := resp.Read(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		callback(value)
	}
}
