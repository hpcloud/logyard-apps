package storage

type FakeFileStorage struct {
	file_path string
}

var (
	IsLoadCalled  bool
	IsWriteCalled bool
)

func NewFakeFileStorage(path string) Storage {
	return &FakeFileStorage{file_path: path}

}

func (f *FakeFileStorage) Write(data interface{}) {
	IsWriteCalled = true

}

func (f *FakeFileStorage) Load(data interface{}) {
	IsLoadCalled = true

}
