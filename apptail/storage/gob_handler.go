package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/ActiveState/log"
	"io/ioutil"
	"os"
)

// exposing these for testing
type Storage interface {
	Write(data interface{})
	Load(data interface{})
}

type FileStorage struct {
	file_path string
}

var (
	path = fmt.Sprintf("%s/.apptail.gob", os.Getenv("HOME"))
)

const FILE_MODE = 0666

func NewFileStorage() Storage {
	return &FileStorage{file_path: path}

}

func (s *FileStorage) Write(data interface{}) {
	m := new(bytes.Buffer)
	enc := gob.NewEncoder(m)
	err := enc.Encode(data)
	if err != nil {
		log.Error(err)

	}
	err = ioutil.WriteFile(s.file_path, m.Bytes(), FILE_MODE)
	if err != nil {
		log.Error(err)

	}

	if err = os.Chmod(s.file_path, FILE_MODE); err != nil {
		log.Error(err)

	}
}

func (s *FileStorage) Load(e interface{}) {
	n, err := ioutil.ReadFile(s.file_path)
	if err != nil {
		log.Error(err)

	}
	p := bytes.NewBuffer(n)
	dec := gob.NewDecoder(p)
	err = dec.Decode(e)
	if err != nil {
		log.Error(err)

	}
}