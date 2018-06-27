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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// LoadFrom is a secondary entry point for the peers package. It requires a
// specific path and if there are any errors loading, the error is returned.
func LoadFrom(path string) (*Peers, error) {
	peers := new(Peers)
	err := peers.Load(path)
	return peers, err
}

// Sync is a helper function that performs a SyncFrom() but looks up the
// url and api key from the environment, expecting the following:
//
// - $PEERS_SYNC_URL: url endpoint for sync GET request
// - $PEERS_SYNC_APIKEY: key to add to headers as X-Api-Key
//
// See the SyncFrom function for more details.
func Sync() (*Peers, error) {
	url := os.Getenv("PEERS_SYNC_URL")
	if url == "" {
		return nil, errors.New("could not find $PEERS_SYNC_URL")
	}

	key := os.Getenv("PEERS_SYNC_APIKEY")
	if key == "" {
		return nil, errors.New("could not find $PEERS_SYNC_APIKEY")
	}

	return SyncFrom(url, key)
}

// SyncFrom is a remote entry point for the peers package. It uses an HTTP
// request to synchronize the peers from a remote host and instantiate the
// peers collection. It expects a url and an api key to perform the GET
// request, adding the api key to the headers as "X-Api-Key".
func SyncFrom(url, apikey string) (*Peers, error) {
	// Conduct the request with a 5 second timeout
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", apikey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Ensure connection is closed on complete
	defer resp.Body.Close()

	// Check the status from the client
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"could not synchronize peers: %s", resp.Status,
		)
	}

	// Parse the body of the response
	peers := new(Peers)
	if err := json.NewDecoder(resp.Body).Decode(&peers); err != nil {
		return nil, err
	}
	return peers, nil
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

	// Make sure the directory exists.
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
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
		name := strings.Split(peer.Hostname, ".")[0]
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

// Get a specific peer by name (which should be unique). If the named peer is
// not found, then an error is returned.
func (p *Peers) Get(hostname string) (*Peer, error) {
	for _, peer := range p.Peers {
		if peer.Name == hostname {
			return peer, nil
		}
	}

	return nil, fmt.Errorf("could not find a peer named '%s'", hostname)
}

//===========================================================================
// Peer Struct
//===========================================================================

// Peer represents a single instance of another replica process or host on
// the network that can be communicated with.
type Peer struct {
	PID         uint16 `json:"pid"`                   // the precedence id of the peer
	Name        string `json:"name"`                  // unique name of the peer
	Description string `json:"description,omitempty"` // a text description of the peer
	Hostname    string `json:"hostname,omitempty"`    // the hostname of the peer
	IPAddr      string `json:"ip_address"`            // the ip address of the peer
	Domain      string `json:"domain,omitempty"`      // the domain name of hte host
	Port        uint16 `json:"port"`                  // the port the replica is listening on

	// Extra information that may be associated with the host
	AWSInstance map[string]string `json:"aws_instance,omitempty"`
}

// IsLocal returns True if the Peer has the same hostname as the localhost.
// Because the host can be specified as a FQDN, this method splits the name
// on "." and inspects the first element of the name.
func (p *Peer) IsLocal() bool {
	hostname, err := os.Hostname()
	if err != nil {
		return false
	}

	name := strings.Split(p.Hostname, ".")[0]
	return name == hostname
}

// Endpoint returns an string with the ip address and the port (or the domain)
// to connect to the peer using TCP.
func (p *Peer) Endpoint(dns bool) string {
	if dns && p.Domain != "" {
		return fmt.Sprintf("%s:%d", p.Domain, p.Port)
	}

	return fmt.Sprintf("%s:%d", p.IPAddr, p.Port)
}

// ZMQEndpoint returns an endpoint to bind or connect on specifically for the
// ZMQ protocol, that is a TCP protocol socket on the specified port.
//
// If server is true, then a bind address of tcp://*:port is returned so that
// messages received on any IP address at that port are handled. If server is
// false then one of two things happen. If the peer IsLocal() then
// tcp://localhost:port is returned. Otherwise tcp://IPAddr:port is returned.
func (p *Peer) ZMQEndpoint(server bool) string {
	if server {
		return fmt.Sprintf("tcp://*:%d", p.Port)
	}

	if p.IsLocal() {
		return fmt.Sprintf("tcp://localhost:%d", p.Port)
	}

	return fmt.Sprintf("tcp://%s:%d", p.IPAddr, p.Port)
}
