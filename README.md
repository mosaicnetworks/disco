# Disco

Disco is a discovery and WebRTC-signaling server for Babble groups.

To start it on `localhost:8080` (hardcoded for now):

```bash
make run
```

It exposes the following REST API:

## Add a group

```bash
POST localhost:8080/group
```

```bash
curl --location --request POST 'localhost:8080/group' \
--header 'Content-Type: application/json' \
--data-raw '{
	"GroupUID": "very unique id",
	"GroupName": "office group",
	"AppID": "BabbleChat",
	"Peers": [
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	],
	"InitialPeers": [
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	]
}'
```

## List all groups

```bash
GET localhost:8080/groups
```

```json
{"very unique id":{"GroupUID":"very unique id","GroupName":"office group","AppID":"BabbleChat","PubKey":"","LastUpdated":1583773505,"Peers":[{"NetAddr":"thenetaddr","PubKeyHex":"thepubkey","Moniker":"Monica"}],"InitialPeers":[{"NetAddr":"thenetaddr","PubKeyHex":"thepubkey","Moniker":"Monica"}]}}
```

## Get a specific group

```bash
GET localhost:8080/groups/{ID}
```

```json
{"very unique id":{"GroupUID":"very unique id","GroupName":"office group","AppID":"BabbleChat","PubKey":"","LastUpdated":1583773505,"Peers":[{"NetAddr":"thenetaddr","PubKeyHex":"thepubkey","Moniker":"Monica"}],"InitialPeers":[{"NetAddr":"thenetaddr","PubKeyHex":"thepubkey","Moniker":"Monica"}]}}
```

## Update a group

```bash
PATCH localhost:8080/groups/2
```

```bash
curl --location --request POST 'localhost:8080/groups/very unique id' \
--header 'Content-Type: application/json' \
--data-raw '{
	"GroupUID": "very unique id",
	"GroupName": "office group modified",
	"AppID": "BabbleChat",
	"Peers": [
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	],
	"InitialPeers": [
		{
			"NetAddr":"thenetaddr",
			"PubKeyHex":"thepubkey",
			"Moniker":"Monica"
		}
	]
}'
```

## Delete a group

```bash
DELETE localhost:8080/groups/2
```

```
The group with ID 2 has been deleted successfully
```
