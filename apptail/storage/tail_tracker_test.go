package storage

import (
	"testing"
)

var (
	instanceKey = "fakeDockerId12345"
	childKey    = "path/to/some/file/stderr.log"
	debug       = false
)

func TestRegisterInstance_PassedNonExsitingInstanceKey_ItShouldAddToMapOfInstances(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage, debug)

	tracker.RegisterInstance(instanceKey)

	if tracker.IsInstanceRegistered(instanceKey) {
		t.Log("Passed")

	} else {
		t.Fail()

	}

}

func Test_InitializeChildNode_CalledWithCorrectArgs_IfNotRegisteredItShouldAssociateToInstance(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage, debug)

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
	tracker := NewTracker(fakeFileStorage, debug)

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
	tracker := NewTracker(fakeFileStorage, debug)

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
	tracker := NewTracker(fakeFileStorage, debug)

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
	tracker := NewTracker(fakeFileStorage, debug)
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
	tracker := NewTracker(fakeFileStorage, debug)
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
	tracker := NewTracker(fakeFileStorage, debug)

	tracker.LoadTailers()

	if IsLoadCalled == true {
		t.Log("passed")

	} else {
		t.Fail()

	}

}

func TestCommit_WhenCalled_ItCallsExpectedMethods(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage, debug)

	tracker.Commit()

	if IsWriteCalled && IsEncodeCalled {
		t.Log("passed")
	} else {
		t.Fail()

	}
}

func TestCommit_UnderlyingCallReturnsError_CommitBubbleUpTheError(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage, debug)

	ThrowError = true

	err := tracker.Commit()

	if err != nil {
		t.Log("passed")

	} else {
		t.Fail()
	}

}

func TestCleanUp_WhenCalledWithAListOfValidIds_ShouldRemoveOldContainerId(t *testing.T) {
	fakeFileStorage := NewFakeFileStorage("somepath")
	tracker := NewTracker(fakeFileStorage, debug)

	validIds := make(map[string]bool)
	validIds["dockerId1"] = true
	validIds["dockerId2"] = true
	validIds["dockerId3"] = true

	tracker.RegisterInstance("dockerId1")
	tracker.RegisterInstance("dockerId2")

	tracker.CleanUp(validIds)

	if tracker.IsInstanceRegistered("dockerId3") {
		t.Fail()

	} else {
		t.Log("pass")

	}

}

func Test_getEntriesToCleanUp_WhenCalledWithTwoValidMaps_ItShouldReturnInvalidInstances(t *testing.T) {

	map_one := map[string]bool{
		"docker1": true,
		"docker2": true,
		"docker3": true,
	}

	map_two := map[string]bool{
		"docker1": true,
		"docker2": true,
	}

	invalidInstances := getInvalidInstances(map_one, map_two)

	if len(invalidInstances) < 1 {
		t.Fail()

	} else {
		t.Log("passed")

	}

}
