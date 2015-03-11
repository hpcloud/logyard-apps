package storage

import (
	"testing"
)

var (
	instanceKey = "fakeDockerId12345"
	childKey    = "path/to/some/file/stderr.log"
)

func TestRegisterInstance_PassedNonExsitingInstanceKey_ItShouldAddToMapOfInstances(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)

	tracker.RegisterInstance(instanceKey)

	if tracker.IsInstanceRegistered(instanceKey) {
		t.Log("Passed")

	} else {
		t.Fail()

	}

}

func Test_InitializeChildNode_CalledWithCorrectArgs_IfNotRegisteredItShouldAssociateToInstance(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)

	tracker.RegisterInstance(instanceKey)
	tracker.InitializeChildNode(instanceKey, childKey, 1024)

	if tracker.IsChildNodeInitialized(instanceKey, childKey) {
		t.Log("passed")

	} else {
		t.Fail()

	}

}

func TestGetFileCachedOffset_CalledWithCorrectArgs_ItShouldReturnCorrectCachedOffsetFromTailNodeMap(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)

	var offset int64 = 1024

	tracker.RegisterInstance(instanceKey)
	tracker.InitializeChildNode(instanceKey, childKey, offset)

	cached_offset := tracker.GetFileCachedOffset(instanceKey, childKey)

	if cached_offset == offset {
		t.Log("passed")

	} else {
		t.Fail()

	}

}

func TestUpdate_CalledFrequently_ItShouldIncrementOffsetPosition(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)

	var offset int64 = 1024

	tracker.RegisterInstance(instanceKey)
	tracker.InitializeChildNode(instanceKey, childKey, offset)

	for i := 0; i <= 7; i++ {
		offset++
		tracker.Update(instanceKey, childKey, offset)

	}

	currentOffset := tracker.GetFileCachedOffset(instanceKey, childKey)

	if currentOffset == offset {
		t.Log("passed")

	} else {
		t.Fail()

	}

}

func TestUpdate_CalledWithNonExistingInstanceKey_ItSholdNotIncrementOffset(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)

	var offset int64 = 1024

	tracker.RegisterInstance(instanceKey)
	tracker.InitializeChildNode(instanceKey, childKey, offset)

	currentOffset := tracker.GetFileCachedOffset("badInstanceKey", childKey)

	if currentOffset == offset {
		t.Fail()

	} else {
		t.Log("passed")

	}

}

func TestRemove_CalledWithaCorrectKey_ItshouldRemoveInstanceFromCached(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)
	tracker.RegisterInstance(instanceKey)

	tracker.Remove(instanceKey)

	if tracker.IsInstanceRegistered(instanceKey) {
		t.Fail()

	} else {
		t.Log("Passed")

	}
}

func TestRemove_CalledWithCorrectKey_ItShouldAlsoRemoveFromFile(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)
	tracker.RegisterInstance(instanceKey)

	tracker.Remove(instanceKey)

	if IsWriteCalled == true {
		t.Log("passed")

	} else {
		t.Fail()

	}

}

func TestLoadTailers_WhenCalled_ItLoadsFromGobFile(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)

	tracker.LoadTailers()

	if IsLoadCalled == true {
		t.Log("passed")

	} else {
		t.Fail()

	}

}

func TestSubmit_WhenCalled_ItCallsStorageWriteToWriteToFile(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage)

	tracker.Submit()

	if IsWriteCalled == true {
		t.Log("passed")
	} else {
		t.Fail()

	}
}
