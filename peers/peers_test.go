package peers

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// Test that the Peers collection can be loaded from a file.
func TestPeersLoad(t *testing.T) {
	peers := new(Peers)
	if err := peers.Load("testdata/peers.json"); err != nil {
		t.Error(err)
	}

	if len(peers.Peers) != 6 {
		t.Errorf("expected to load 6 peers but got %d", len(peers.Peers))
	}

	if peers.Info["num_replicas"].(float64) != 6 {
		t.Error("peers metadata not successfully loaded")
	}

	if peers.Info["updated"].(string) == "" {
		t.Error("peers metadata not successfully loaded")
	}
}

// Test that the Peers collection can be dumped to disk.
func TestPeersDump(t *testing.T) {
	peers := new(Peers)
	peers.Info = make(map[string]interface{})
	peers.Info["num_replicas"] = 2
	peers.Info["updated"] = time.Now()
	peers.Peers = []*Peer{
		{
			PID: 1, Name: "alpha", IPAddr: "10.10.10.1", Port: 3264,
		},
		{
			PID: 2, Name: "bravo", IPAddr: "10.10.10.2", Port: 3264,
		},
	}

	// Create a temp file
	tmpfile, err := ioutil.TempFile("", "peers.json")
	if err != nil {
		t.Fatal(err)
	}
	path := tmpfile.Name()
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Ensure tempfile is cleaned up
	defer os.Remove(path)

	// Dump the peers
	if err := peers.Dump(path); err != nil {
		t.Error(err)
	}

	// Ensure the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file not created during peers dump")
	}

}

// Test that the local finds all the local hosts
func TestLocal(t *testing.T) {
	peers := new(Peers)
	if err := peers.Load("testdata/peers.json"); err != nil {
		t.Error(err)
	}

	local := peers.Local("bravo")
	if len(local) != 3 {
		t.Error("could not find all local replicas for bravo")
	}
}

// Test that we find a single localhost by pid
func TestLocalhost(t *testing.T) {
	peers := new(Peers)
	if err := peers.Load("testdata/peers.json"); err != nil {
		t.Error(err)
	}

	// Find second bravo
	local, err := peers.Localhost("bravo", uint32(11))
	if err != nil {
		t.Error(err)
	}

	if local.Name != "bravo-11" {
		t.Error("did not find right localhost with arguments")
	}

	// Find first charlie
	local, err = peers.Localhost("charlie", 0)
	if err != nil {
		t.Error(err)
	}

	if local.Name != "charlie" {
		t.Error("did not find right localhost with arguments")
	}
}
