package pid

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Location to write temporary files for testing
var tmpDir string

// Test the that the Path function returns the same path on demand.
func TestPath(t *testing.T) {
	path := Path("test.pid")
	if path == "" {
		t.Error("did not return a path")
	}

	if path != Path("test.pid") {
		t.Error("returned different path on subsequent call")
	}
}

// Test the PID New function
func TestNew(t *testing.T) {
	if err := makeTmpDir(); err != nil {
		t.Fatal(err)
	}
	defer removeTmpDir()

	path := filepath.Join(tmpDir, "test.pid")
	pid := New(path)

	if pid.Path() != path {
		t.Error("path in new doesn't match associated path")
	}
}

// Test that the free function works on a path that doesn't exist
func TestFreeEmpty(t *testing.T) {
	if err := makeTmpDir(); err != nil {
		t.Fatal(err)
	}
	defer removeTmpDir()

	path := filepath.Join(tmpDir, "test.pid")
	pid := New(path)

	// Ensure the path doesn't exist
	if exists, _ := pathExists(path); exists {
		t.Fatal("test path already exists")
	}

	// Ensure Free doesn't return an error
	if err := pid.Free(); err != nil {
		t.Error(err)
	}
}

// Test that a PID can be saved when the pid doesn't already exist
func TestPIDSave(t *testing.T) {
	if err := makeTmpDir(); err != nil {
		t.Fatal(err)
	}
	defer removeTmpDir()

	path := filepath.Join(tmpDir, "test.pid")
	pid := New(path)

	// Ensure the path doesn't exist
	if exists, _ := pathExists(path); exists {
		t.Fatal("test path already exists")
	}

	// Ensure there is no information in the PID until save
	if pid.PID != 0 || pid.PPID != 0 {
		t.Error("data exists in PID before save")
	}

	// Save the PID file
	if err := pid.Save(); err != nil {
		t.Error(err)
	}

	// Ensure the PID is populated
	if pid.PID == 0 || pid.PPID == 0 {
		t.Error("data does not exist in PID after save")
	}

	// Ensure the path exists
	if exists, _ := pathExists(path); !exists {
		t.Fatal("pid file does not exist after save")
	}
}

// Test that a PID can be freed after being saved
func TestPIDSaveFree(t *testing.T) {
	if err := makeTmpDir(); err != nil {
		t.Fatal(err)
	}
	defer removeTmpDir()

	path := filepath.Join(tmpDir, "test.pid")
	pid := New(path)

	// Ensure the path doesn't exist
	if exists, _ := pathExists(path); exists {
		t.Fatal("test path already exists")
	}

	// Save the PID file
	if err := pid.Save(); err != nil {
		t.Error(err)
	}

	// Ensure the path exists
	if exists, _ := pathExists(path); !exists {
		t.Fatal("pid file does not exist after save")
	}

	// Free the PID file
	if err := pid.Free(); err != nil {
		t.Error(err)
	}

	// Ensure the path doesn't exist anymore
	// Ensure the path doesn't exist
	if exists, _ := pathExists(path); exists {
		t.Error("test path exists after free")
	}
}

// Test that a PID cannot be saved when the pid already exists
func TestSingleProcess(t *testing.T) {
	if err := makeTmpDir(); err != nil {
		t.Fatal(err)
	}
	defer removeTmpDir()

	path := filepath.Join(tmpDir, "test.pid")
	pid := New(path)

	// Ensure the path doesn't exist
	if exists, _ := pathExists(path); exists {
		t.Fatal("test path already exists")
	}

	// Save the PID file
	if err := pid.Save(); err != nil {
		t.Error(err)
	}

	// Ensure the path exists
	if exists, _ := pathExists(path); !exists {
		t.Fatal("pid file does not exist after save")
	}

	// Try to save a new PID (should error)
	pid2 := New(path)
	if err := pid2.Save(); err == nil {
		t.Error("second pid was allowed to save without error")
	}
}

// Test that a PID can be loaded from an existing file
func TestLoad(t *testing.T) {
	if err := makeTmpDir(); err != nil {
		t.Fatal(err)
	}
	defer removeTmpDir()

	path := filepath.Join(tmpDir, "test.pid")

	// Write a test PID file
	testData := map[string]int{
		"pid":  23,
		"ppid": 22,
		"port": 50800,
	}

	// Write the test data as JSON
	data, err := json.Marshal(testData)
	if err != nil {
		t.Fatal(err)
	}

	// And spit it out to a file.
	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a new PID and make sure it has no data
	pid := New(path)
	if pid.PID != 0 || pid.PPID != 0 {
		t.Error("data exists in PID before load")
	}

	// Load the PID file
	if err = pid.Load(); err != nil {
		t.Error(err)
	}

	if pid.PID == 0 || pid.PPID == 0 {
		t.Error("data does not exist in PID after load")
	}
}

//===========================================================================
// Test Helper Functions
//===========================================================================

// Make temporary directory
func makeTmpDir() (err error) {
	tmpDir, err = ioutil.TempDir("", "com.bengfort.x.pid")
	return err
}

//  Remove the temporary directory
func removeTmpDir() (err error) {
	return os.RemoveAll(tmpDir)
}

// Check if a path exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}
