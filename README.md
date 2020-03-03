# Disco

Disco is a discovery server for Babble groups.

To start it on `localhost:8080` (hardcoded for now):

`go run main.go`

It exposes the following REST API:

## Add a group

```bash
POST localhost:8080/group
```

```bash
curl --location --request POST 'localhost:8080/group' \
--header 'Content-Type: application/json' \
--data-raw '{
	"ID":"2",
	"Title": "office group",
	"Description": "this group is for office people",
	"Peers": {
		"peers": [
			{
				"NetAddr":"thenetaddr",
				"PubKeyHex":"thepubkey",
				"Moniker":"Monica"
			}
			]
	}
}'
```

## List all groups

```bash
GET localhost:8080/groups
```

```json
[{"ID":"1","Title":"Introduction to Golang","Description":"Come join us for a chance to learn how golang works and get to eventually try it out","Peers":{"peers":[{"NetAddr":"Peer0Addr","PubKeyHex":"XXX","Moniker":"Peer0"}]}},{"ID":"","Title":"","Description":"","Peers":null},{"ID":"","Title":"office group","Description":"this group is for office people","Peers":{"peers":[{"NetAddr":"thenetaddr","PubKeyHex":"thepubkey","Moniker":"Monica"}]}},{"ID":"2","Title":"office group","Description":"this group is for office people","Peers":{"peers":[{"NetAddr":"thenetaddr","PubKeyHex":"thepubkey","Moniker":"Monica"}]}}]
```

## Get a specific group

```bash
GET localhost:8080/groups/{ID}
```

```json
{"ID":"1","Title":"Introduction to Golang","Description":"Come join us for a chance to learn how golang works and get to eventually try it out","Peers":{"peers":[{"NetAddr":"Peer0Addr","PubKeyHex":"XXX","Moniker":"Peer0"}]}}
```

## Update a group

```bash
PATCH localhost:8080/groups/2
```

```bash
curl --location --request PATCH 'localhost:8080/groups/2' \
--header 'Content-Type: application/json' \
--data-raw '{
	"ID":"2",
	"Title": "office group",
	"Description": "this group has been modified",
	"Peers": {
		"peers": [
			{
				"NetAddr":"thenetaddr",
				"PubKeyHex":"thepubkey",
				"Moniker":"Monica"
			}
			]
	}
}'
```

## Delete a group

```bash
DELETE localhost:8080/groups/2
```

```
The group with ID 2 has been deleted successfully
```
