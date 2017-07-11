package peers

import (
	"os"
	"os/user"
	"path/filepath"
)

// The expected direcotry and file name of the peers file.
const dirname = "fluidfs"
const filename = "peers.json"

//===========================================================================
// Path Lookup
//===========================================================================

// Path returns the recommended location to store a peers.json file for
// external services to use. It returns the environment peers path if it
// exists, otherwise the user peers path. If there is no user directory, then
// it falls back on the system peers path.
func Path() string {
	if path := envPeers(); path != "" {
		return path
	}

	if path := userPeers(); path != "" {
		return path
	}

	return systemPeers()
}

//===========================================================================
// Functions that return lookup paths
//===========================================================================

// Returns the peers paths in priority order, filtering out empty strings
// from the various lookup functions.
func peersPaths() []string {
	paths := make([]string, 0, 4)

	if env := envPeers(); env != "" {
		paths = append(paths, env)
	}

	if cwd := cwdPeers(); cwd != "" {
		paths = append(paths, cwd)
	}

	if user := userPeers(); user != "" {
		paths = append(paths, user)
	}

	if system := systemPeers(); system != "" {
		paths = append(paths, system)
	}

	return paths
}

// Returns the path to the peers file in the system directory. For now it
// just looks up the file in /etc but could be extended for other systems.
func systemPeers() string {
	return filepath.Join("/", "etc", dirname, filename)
}

// Returns the path to the peers file in the user directory. Returns an empty
// string if there is no user directory.
func userPeers() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(usr.HomeDir, "."+dirname, filename)
}

// Returns the path to the peers file in the current working directory. It
// returns an empty string if there is no current working directory.
func cwdPeers() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(cwd, filename)
}

// Returns the path to the peers file specified by the environment or an
// empty string if there is no path in the environment.
func envPeers() string {
	return os.Getenv("PEERS_PATH")
}
