package storage

import "errors"

type FakeFileStorage struct {
	file_path string
}

var (
	IsEncodeCalled bool
	IsLoadCalled   bool
	IsWriteCalled  bool
	ThrowError     bool
)

func NewFakeFileStorage(path string) Storage {
	return &FakeFileStorage{file_path: path}

}

func (s *FakeFileStorage) Encode(data interface{}) ([]byte, error) {
	IsEncodeCalled = true
	var byte []byte

	if ThrowError {
		return nil, errors.New("something went wrong while trying to Encode")

	}
	return byte, nil
}

func (f *FakeFileStorage) Write(buf []byte) error {
	IsWriteCalled = true

	if ThrowError {
		return errors.New("something went wrong while trying to Write to the file")

	}
	return nil
}

func (f *FakeFileStorage) Load(data interface{}) {
	IsLoadCalled = true
}
