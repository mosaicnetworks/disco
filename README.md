# Disco

Disco is a discovery server for Babble groups. It offers an API to advertise and
manage Babble groups, as well as a WebRTC signaling mechanism and TURN server,
which can be used by Babble in `webrtc` mode to establish p2p connections in the
face of NATs. 

## Table of Contents
 * [Usage](#usage)
 * [Discovery](#discovery)
	+ [Add a group](#add-a-group)
	+ [List groups](#list-groups)
	+ [Get a specific group](#get-a-specific-group)
	+ [Update a group](#update-a-group)
	+ [Delete a group](#delete-a-group)
	+ [TTL](#ttl)
 * [WebRTC Signaling](#webrtc-signaling)
 * [TURN](#turn)
 * [Caveats](#caveats)

## Usage

```bash
Discovery service for Babble

Usage:
  disco [flags]

Flags:
      --address string          Advertise address (use public address) (default "0.0.0.0")
      --cert-file string        File containing TLS certificate (default "cert.pem")
      --disco-port string       Discovery API port (default "1443")
  -h, --help                    help for disco
      --ice-password string     ICE server password corresponding to username (default "test")
      --ice-port string         ICE server port (default "3478")
      --ice-username string     ICE server userame. Only this user will be allowed to use the ICE server (default "test")
      --key-file string         File containing certificate key (default "key.pem")
      --realm string            Administrative routing domain within the WebRTC signaling (default "main")
      --signal-port string      WebRTC-Signaling port (default "2443")
      --ttl duration            Group Time To Live, after which groups will be deleted (default 5m0s)
      --ttl-hearbeat duration   Ticker frequency for checking group TTL (default 1m0s)
```

The discovery API, WebRTC-signaling router, and TURN server are exposed on 
different ports, `disco-port` (default 1443), `signal-port` (default 2443), and
`ice-port` respectively. 

The discovery API and WebRTC-signaling router bind to all interfaces `0.0.0.0`
and are secured with TLS, with the same underlying certificate. The certificate 
and key files are specified with `cert-file` and `key-file` options.

The `TURN` server also binds to `0.0.0.0` but it advertises itself at the 
address specified by `--address`. This must be the public IP of the machine 
running the server, as it will be used as a TURN relay address.

The `ice-username` and `ice-password` options define the credentials of a single
user allowed to authenticate and use the TURN server. `Babble` has homonymous 
config options to match the username and password fields when using the Disco
TURN server.  

To get started with a localhost disco server, simply run:

```bash
make run
```

## Discovery

The discovery API offers a mechanism to create, discover, and manage Babble
groups. A Babble group defines a set of peers engaged in a Babble consensus
network.

The following sections describe the API calls using curl, but first we need curl
to trust the self-signed certificate in `test_data/cert.pem`:

```bash
export CURL_CA_BUNDLE=test_data/cert.pem
```

### Add a group

```bash
POST https://localhost:1443/group
```

```bash
 curl --location --request POST 'https://localhost:1443/group' \
--header 'Content-Type: application/json' \
--data-binary @new_group.json
```

where `new_group.json` contains the json for a new group:

```json
{
	"Name": "office group",
	"AppID": "BabbleChat",
	"Peers": [ 
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	],
	"GenesisPeers": [
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	]
}
```

Note that the `ID` field of the group is omitted. This is because it will be
randomly generated by the server. If the `ID` is specified, the server will 
either update the group in place, or create a new one with that ID.

### List groups

```bash
GET https://localhost:1443/groups?app-id=
```

Use the `app-id` query parameter to return only the groups belonging to a 
certain application.

```json
{
	"8f41c928-360b-4202-90f1-a7efa6b7ffd3": {
		"ID":"8f41c928-360b-4202-90f1-a7efa6b7ffd3",
		"Name":"office group",
		"AppID":"BabbleChat",
		"PubKey":"",
		"LastUpdated":1583773505,
		"Peers":[
			{
				"NetAddr":"thenetaddr",
				"PubKeyHex":"thepubkey",
				"Moniker":"Monica"
			}
		],
		"InitialPeers":[
			{
				"NetAddr":"thenetaddr",
				"PubKeyHex":"thepubkey",
				"Moniker":"Monica"
			}
		]
	}
}
```

### Get a specific group

```bash
GET https://localhost:1443/groups/{ID}
```

```json
{
	"ID":"8f41c928-360b-4202-90f1-a7efa6b7ffd3",
	"Name":"office group",
	"AppID":"BabbleChat",
	"PubKey":"",
	"LastUpdated":1583773505,
	"Peers":[
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	],
	"InitialPeers":[
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	]
}
```

### Update a group

```bash
PATCH https://localhost:1443/groups/2
```

```bash
 curl --location --request PATCH 'https://localhost:1443/group' \
--header 'Content-Type: application/json' \
--data-binary @updated_group.json
```

### Delete a group

```bash
DELETE https://localhost:1443/groups/2
```

```
The group with ID 2 has been deleted successfully
```

### TTL

It is possible to set a `Time To Live` (`--ttl`), and a ticker frequency 
(`--ttl-hearbeat`), to ensure that groups get deleted from the server after 
their TTL has expired. 

## WebRTC Signaling

The WebRTC Signaling Server enables Babble nodes to exchange connection 
information (SDP) prior to establishing a direct P2P link. It is implemented as 
a WAMP server (Web Application Messaging Protocol) which is basically RPC over 
secure websockets. 

When using Babble with `webrtc` enabled, you can point to the Disco server's
`signal-port`. If the server's TLS certificate is self-signed, you can copy the
`cert.pem` file in Babble's data-directory.

## TURN

The TURN server offers STUN/TURN services which can be used by Babble to 
establish P2P connections. It helps devices punch holes through NATs when this
is possible, or relays packets on their behalf.

When using Babble with `webrtc` enabled, you can point to the Disco server's 
`ice-port` with the appropriate `ice-username` and `ice-password` to 
authenticate to, and use the Disco server's TURN services.

## Improvements

Ideally, we would like the same disco server to be used by multiple apps. Group 
discovery, signaling, and TURN services, should occur within separate 
administrative domains (realms) for each application. This scheme would require 
developers to register their applications with the disco server, which would 
spin-up new instances of the services for each application and somehow 
authenticate application clients as to prevent them from accessing groups from 
other applications.

This separation is not implemented yet. Applications are sharing the same 
signaling router (same realm), and TURN server, and are not authenticated with 
the disco server. Groups are indexed by AppID and it is possible to query only 
groups pertaining to an AppID. This allows developers to only display those 
groups that are relevant to their application, but it is definitely not a secure
system. 

The group database is not persisted, meaning that all groups are lost when the
server is restarted.

Updating and deleting groups should be protected to require authorisation from
enough group members. Ultimately we should implement a Babble light-client for
every group.

TTL should be configurable at the group level, to enable group creators to set
their own TTL.