package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var (
	test_path = fmt.Sprintf("%s/.apptail.gob.test", os.Getenv("HOME"))
)

type TestInterface struct {
	Value string
}

func TestLoad_CalledWithCorrectArgs_IfFileNotExistItCreatesOne(t *testing.T) {

	// test clean up
	defer cleanUp()

	storage := NewFileStorage(test_path)

	var toLoad TestInterface

	storage.Load(&toLoad)

	if _, err := os.Stat(test_path); os.IsNotExist(err) {
		t.Fail()
	} else {
		t.Log("pass")

	}
}

func TestLoad_CalledWhenFileExist_ItTriesToDecodeGivenInterface(t *testing.T) {

	// test clean up
	defer cleanUp()

	storage := NewFileStorage(test_path)

	testInterface := &TestInterface{Value: "Hello Wrold"}

	bytes, _ := storage.Encode(testInterface)
	err := storage.Write(bytes)

	if err != nil {
		t.Log(err)

	}

	var toLoad TestInterface

	storage.Load(&toLoad)

	if toLoad.Value == testInterface.Value {
		t.Log("pass")

	} else {
		t.Fail()

	}

}

func TestWrite_CalledWithACorrectArgs_ItShouldCreateFileIfNotExist(t *testing.T) {
	// test clean up
	defer cleanUp()

	storage := NewFileStorage(test_path)

	testInterface := &TestInterface{Value: "Hello Wrold"}

	bytes, _ := storage.Encode(testInterface)

	storage.Write(bytes)

	if _, err := os.Stat(test_path); os.IsNotExist(err) {
		t.Fail()
	} else {
		t.Log("pass")

	}
}

func cleanUp() {
	os.Remove(test_path)

}
