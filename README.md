# Disco

Disco is a discovery and WebRTC-signaling server for Babble groups.

## Table of Contents
 * [Usage](#usage)
 * [Discovery](#discovery)
	+ [Add a group](#add-a-group)
	+ [List all groups](#list-all-groups)
	+ [Get a specific group](#get-a-specific-group)
	+ [Update a group](#update-a-group)
	+ [Delete a group](#delete-a-group)
 * [Caveats](#caveats)

## Usage

```bash
Discovery service for Babble

Usage:
  disco [flags]

Flags:
      --address string       Address of the server (default "localhost")
      --cert-file string     File containing TLS certificate (default "cert.pem")
      --disco-port string    Discovery API port (default "1443")
  -h, --help                 help for disco
      --key-file string      File containing certificate key (default "key.pem")
      --realm string         Administrative routing domain within the WebRTC signaling (default "main")
      --signal-port string   WebRTC-Signaling port (default "2443")
```

The discovery API and WebRTC-signaling router are exposed on different ports,
`disco-port` (default 1443) and `signal-port` (default 2443) respectively.

Both services are secured with TLS, and the same underlying certificate. The
certificate and key files are specified with `cert-file` and `key-file` options.

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

## Improvements

Ideally, we would like the same disco server to be used by multiple apps. Group 
discovery and WebRTC signaling should occur within separate administrative 
domains (realms) for each application. This scheme would require application 
developers to register their applications with the disco server, which would 
spin-up new instances of signaling servers for each application and somehow 
authenticate application clients as to prevent them from accessing groups from 
other applications.

This separation is not implemented yet. Applications are sharing the same 
signaling router (same realm), and are not authenticated with the disco server. 
Groups are indexed by AppID and it is possible to query only groups pertaining 
to an AppID. This allows developers to only display those groups that are 
relevant to their application, but it is definitely not a secure system. 

The group database is not persisted, meaning that all groups are lost when the
server is restarted.

Updating and deleting groups should be protected to require authorisation from
enough group members. Ultimately we should implement a Babble light-client for
every group.