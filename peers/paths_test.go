package peers

import (
	"os"
	"testing"
)

// Ensure that the peers path returns three items without the environment
// variable and four when it's set. Also ensure that there are no empty
// strings in the peers paths.
func TestPeersPath(t *testing.T) {
	paths := peersPaths()
	if len(paths) != 3 {
		t.Error("expected three peers paths without envvar, but got ", len(paths))
	}

	os.Setenv("PEERS_PATH", "/tmp/peers.json")
	paths = peersPaths()
	if len(paths) != 4 {
		t.Error("expected four peers paths without envvar, but got ", len(paths))
	}

	// TODO: make this better
	for _, path := range paths {
		if path == "" {
			t.Error("found an empty string in the peers paths")
		}
	}
}
