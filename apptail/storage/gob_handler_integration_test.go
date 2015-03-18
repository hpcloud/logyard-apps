package storage

import (
	"fmt"
	"os"
	"testing"
)

var (
	test_path = fmt.Sprintf("%s/.apptail.gob", os.Getenv("HOME"))
)

type TestInterFace struct {
	Value string
}

func TestLoad_CalledWithCorrectArgs_IfFileNotExistItCreatesOne(t *testing.T) {
	storage := NewFileStorage(test_path)

	var toLoad TestInterFace

	storage.Load(&toLoad)

	if _, err := os.Stat(test_path); os.IsNotExist(err) {
		t.Fail()
	} else {
		t.Log("pass")

	}

	// test clean up
	cleanUp()

}

func TestLoad_CalledWhenFileExist_ItTriesToDecodeGivenInterface(t *testing.T) {

	storage := NewFileStorage(test_path)

	testInterface := &TestInterFace{Value: "Hello Wrold"}

	bytes, _ := storage.Encode(testInterface)
	err := storage.Write(bytes)

	if err != nil {
		t.Log(err)

	}

	var toLoad TestInterFace

	storage.Load(&toLoad)

	if toLoad.Value == testInterface.Value {
		t.Log("pass")

	} else {
		t.Fail()

	}
	cleanUp()

}

func TestWrite_CalledWithACorrectArgs_ItShouldCreateFileIfNotExist(t *testing.T) {
	storage := NewFileStorage(test_path)

	testInterface := &TestInterFace{Value: "Hello Wrold"}

	bytes, _ := storage.Encode(testInterface)

	storage.Write(bytes)

	if _, err := os.Stat(test_path); os.IsNotExist(err) {
		t.Fail()
	} else {
		t.Log("pass")

	}

	cleanUp()
}

func cleanUp() {
	os.Remove(test_path)

}
