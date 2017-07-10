# Peers

**Provides support for identifying other peer processes to communicate with in the network.**

Package peers provides helpers for defining and synchronizing remote peers on a network. The package essentially surrounds JSON files with remote host definitions in several standard locations. It also provides helper functionality for synchronizing the JSON data with a remote server as well as adding hosts to the file if they are seen on the network or removing hosts if they become offline.

To start using the library, simply create one of the following files:

- `/etc/fluidfs/peers.json`
- `$HOME/.fluidfs/peers.json`
- `$PWD/peers.json`

Or define a path using the `$PEERS_PATH` environment variable. When calling `peers.Load()`, it will look in each of these locations to find and parse the `peers.json` file, returning a Peers object. Alternatively, a new Peers object can be created and a path specified to its `Load()` method to load a specific file not above. The `peers.json` file can also be saved from the Peers object using the `Dump()` method.

The Peers object can also be synchronized from a remote service using the `Sync()` method. Synchronization fetches `peers.json` from a URL that can be specified by the environment, and can also submit an API key along with the request.

Other important helpers include the ability to identify the localhost or peer from the hostname of the system, or to identify all local peer processes. In short, the Peers object is a useful way to manage the configuration of a connected network of communicating devices.
