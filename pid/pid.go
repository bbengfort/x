// Package pid manages pid files for background process management.
package pid

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

//===========================================================================
// Helper Methods
//===========================================================================

// Path is a helper function that computes the best possible PID file for the
// current system, by first attempting to get the user directory then
// resorting to /var/run on Linux systems and elsewhere on other systems.
func Path(filename string) string {
	usr, err := user.Current()
	if err == nil {
		return filepath.Join(usr.HomeDir, ".run", filename)
	}

	return filepath.Join("/", "var", "run", filename)
}

// New PID file at the given location. Note that this function only creates
// an empty PID, which can then be loaded or saved in order to obtain
// process information.
func New(path string) *PID {
	return &PID{path: path}
}

//===========================================================================
// PID File Management
//===========================================================================

// PID describes the server process and is accessed by both the server and the
// command line client in order to facilitate cross-process communication.
type PID struct {
	PID  int    `json:"pid"`  // The process id assigned by the OS
	PPID int    `json:"ppid"` // The parent process id
	path string // The path to the pid file
}

// Save the PID file to disk after first determining the process ids.
// NOTE: This method will fail if the PID file already exists.
func (pid *PID) Save() error {
	var err error

	// Get the currently running Process ID and Parent ID
	pid.PID = os.Getpid()
	pid.PPID = os.Getppid()

	// Marshall the JSON representation
	data, err := json.Marshal(pid)
	if err != nil {
		return err
	}

	path := pid.Path()
	// Ensure that a PID file does not exist (race possible)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Make sure the directory exists.
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// Write the JSON representation of the PID file to disk
		return ioutil.WriteFile(path, data, 0644)
	}

	return fmt.Errorf("PID file exists already at '%s'", path)
}

// Load the PID file -- used by the command line client to populate the PID.
func (pid *PID) Load() error {
	data, err := ioutil.ReadFile(pid.Path())
	if err != nil {
		return fmt.Errorf("no PID file exists at %s; process not running?", pid.Path())
	}

	return json.Unmarshal(data, &pid)
}

// Free the PID file (delete it) -- used by the server on shutdown to cleanup
// and ensure that stray process information isn't just lying about.
// Does not return an error if the PID file does not exist.
func (pid *PID) Free() error {
	// If the PID file doesn't exist, just ignore and return.
	if _, err := os.Stat(pid.Path()); os.IsNotExist(err) {
		return nil
	}

	// Delete the PID file
	return os.Remove(pid.Path())
}

// Path is a getter method to return the location of the PID file on disk.
func (pid *PID) Path() string {
	return pid.path
}

//===========================================================================
// PID Process Management
//===========================================================================

// Process finds and returns the associated operating system process so that
// it can be signaled from a client application.
func (pid *PID) Process() (*os.Process, error) {
	if pid.PID == 0 {
		return nil, errors.New("PID has not yet been saved or loaded")
	}

	return os.FindProcess(pid.PID)
}

// Kill causes the process identified by the PID file to exit immediately
func (pid *PID) Kill() error {
	proc, err := pid.Process()
	if err != nil {
		return err
	}

	return proc.Kill()
}

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (pid *PID) Signal(sig os.Signal) error {
	proc, err := pid.Process()
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}
