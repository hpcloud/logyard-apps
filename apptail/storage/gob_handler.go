package storage

import (
	"bytes"
	"encoding/gob"
	"github.com/ActiveState/log"
	"io/ioutil"
	"os"
)

// exposing these for testing
type Storage interface {
	Encode(data interface{}) ([]byte, error)
	Load(data interface{})
	Write(buf []byte) error
}

type FileStorage struct {
	file_path string
}

const FILE_MODE = 0666

func NewFileStorage(path string) Storage {
	return &FileStorage{
		file_path: path,
	}

}

func (s *FileStorage) Encode(data interface{}) ([]byte, error) {
	m := new(bytes.Buffer)
	enc := gob.NewEncoder(m)
	err := enc.Encode(data)
	if err != nil {
		return nil, err

	}
	return m.Bytes(), nil
}

func (s *FileStorage) Write(buf []byte) error {

	if err := ioutil.WriteFile(s.file_path, buf, FILE_MODE); err != nil {

		return err

	}
	// this extra step to make the file accessible by stackato user
	if err := os.Chmod(s.file_path, FILE_MODE); err != nil {
		return err

	}
	return nil
}

func (s *FileStorage) Load(e interface{}) {

	if _, err := os.Stat(s.file_path); os.IsNotExist(err) {
		log.Infof("Creating %s since it does not exist", s.file_path)
		_, err = os.Create(s.file_path)

	} else {
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

}
