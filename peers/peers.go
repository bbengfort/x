// Package peers provides helpers for defining and synchronizing remote
// peers on a network. The package essentially surrounds JSON files with
// remote host definitions in several standard locations. It also provides
// helper functionality for synchronizing the JSON data with a remote server
// as well as adding hosts to the file if they are seen on the network or
// removing hosts if they become offline.
//
// To start using the library, simply create one of the following files:
//
// - /etc/fluidfs/peers.json
// - $HOME/.fluidfs/peers.json
// - $PWD/peers.json
//
// Or define a path using the $PEERS_PATH environment variable. When calling
// peers.Load(), it will look in each of these locations to find and parse
// the peers.json file, returning a Peers object. Alternatively, a new Peers
// object can be created and a path specified to its Load() method to load a
// specific file not above. The peers.json file can also be saved from the
// Peers object using the Dump() method.
//
// The Peers object can also be synchronized from a remote service using the
// Sync() method. Synchronization fetches peers.json from a URL that can be
// specified by the environment, and can also submit an API key along with
// the request.
//
// Other important helpers include the ability to identify the localhost or
// peer from the hostname of the system, or to identify all local peer
// processes. In short, the Peers object is a useful way to manage the
// configuration of a connected network of communicating devices.
package peers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// Load is the primary entry point for the peers package. It uses a list of
// paths, ordered by priority to find the peers.json file and instantiate the
// peers object from it. If it does not find a peers.json file it simply
// returns an empty collection rather than an error.
//
// The lookup paths for the peers.json file are as follows:
//
// - $PEERS_PATH
// - $PWD/peers.json
// - $HOME/.fluidfs/peers.json
// - /etc/fluidfs/peers.json
//
// At the moment, the first path that is available short circuits the load
// process and all remaining paths are ignored.
func Load() *Peers {
	peers := new(Peers)

	for _, path := range peersPaths() {
		// If there is no error loading the peers, then stop trying to
		// load peers paths because the loading was successful!
		if err := peers.Load(path); err == nil {
			break
		}
	}

	return peers
}

//===========================================================================
// Peers Collection
//===========================================================================

// Peers is a collection of network hosts or processes that can be
// communicated with, along with associated metadata. The Peers object is the
// primary interaction with files on disk and exposes methods that select
// relevent hosts and addresses.
type Peers struct {
	Info  map[string]interface{} `json:"info"`     // metadata associated with the collection
	Peers []*Peer                `json:"replicas"` // the network peers (also called replicas)
	path  string                 // the path that was successfully loaded
}

// Load the peers collection from a JSON file on disk. If the peers are
// successfully loaded, the path it was loaded from is stored and no error
// is returned.
func (p *Peers) Load(path string) error {
	// Read the data from disk
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	// Save the path and return nil
	p.path = path
	return nil
}

// Dump the peers collection as a JSON file to disk. If an empty string is
// passed in as an argument, then it will dump to the location on disk it
// was loaded from.
func (p *Peers) Dump(path string) error {

	// Find the correct path to dump to
	if path == "" {
		if p.path == "" {
			return errors.New("no path specified to dump peers.json to")
		}
		path = p.path
	}

	// Marshal the JSON data
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	// Write the data to the file
	return ioutil.WriteFile(path, data, 0644)
}

// Local returns the peers that are local to the specified host by comparing
// the Peer's Host parameter with the hostname. If the hostname is an empty
// string, then the hostname of the system is used. No errors are returned
// from this function, instead returning an empty list of not found. Because
// the Peer's hostname can be a FQDN, the Host is split on "." and the first
// element is used. If no local replicas are found, it returns an empty list.
func (p *Peers) Local(hostname string) []*Peer {
	// Resolve the hostname
	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	peers := make([]*Peer, 0)
	for _, peer := range p.Peers {
		name := strings.Split(peer.Host, ".")[0]
		if name == hostname {
			peers = append(peers, peer)
		}
	}

	return peers

}

// Localhost returns the peer that is defined by the current localhost. Note
// that multiple peers can reside on a single machine, but this method will
// only return one Peer. Filtering multiple local replicas can be done with
// a precedence ID. If no ID is specified (e.g. 0) then the first local peer
// is returned. If no matching peer is found then an error is returned.
func (p *Peers) Localhost(hostname string, pid uint16) (*Peer, error) {
	for _, peer := range p.Local(hostname) {
		if pid == 0 || peer.PID == pid {
			return peer, nil
		}
	}

	return nil, errors.New("could not find a matching localhost")
}

//===========================================================================
// Peer Struct
//===========================================================================

// Peer represents a single instance of another replica process or host on
// the network that can be communicated with.
type Peer struct {
	PID    uint16 `json:"pid"`     // the precedence id of the peer
	Name   string `json:"name"`    // unique name of the peer
	Addr   string `json:"address"` // the network address of the peer
	Host   string `json:"host"`    // the hostname of the peer
	IPAddr string `json:"ipaddr"`  // the ip address of the peer
	Port   uint16 `json:"port"`    // the port the replica is listening on
}
