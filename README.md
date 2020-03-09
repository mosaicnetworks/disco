# Disco

Disco is a discovery server for Babble groups.

To start it on `localhost:8080` (hardcoded for now):

```bash
cd server/cmd
go run main.go
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
	"Title": "office group",
	"Description": "this group is for office people",
	"Peers": [
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
{"07e42b07-620f-40d0-a062-c15d3e22eb74":{"id":"07e42b07-620f-40d0-a062-c15d3e22eb74","title":"office group","description":"this group is for office people","peers":[{"NetAddr":"thenetaddr","PubKeyHex":"thepubkey","Moniker":"Monica"}]},"a8dd4c0a-f025-4618-8a3b-66fb0553b034":{"id":"a8dd4c0a-f025-4618-8a3b-66fb0553b034","title":"Group1","description":"Useless Group","peers":[{"NetAddr":"alice@localhost","PubKeyHex":"XXX","Moniker":"Alice"},{"NetAddr":"bob@localhost","PubKeyHex":"YYY","Moniker":"Bob"},{"NetAddr":"charlie@localhost","PubKeyHex":"ZZZ","Moniker":"Charlie"}]}}
```

## Get a specific group

```bash
GET localhost:8080/groups/{ID}
```

```json
{"id":"a8dd4c0a-f025-4618-8a3b-66fb0553b034","title":"Group1","description":"Useless Group","peers":[{"NetAddr":"alice@localhost","PubKeyHex":"XXX","Moniker":"Alice"},{"NetAddr":"bob@localhost","PubKeyHex":"YYY","Moniker":"Bob"},{"NetAddr":"charlie@localhost","PubKeyHex":"ZZZ","Moniker":"Charlie"}]}
```

## Update a group

```bash
PATCH localhost:8080/groups/2
```

```bash
curl --location --request PATCH 'localhost:8080/groups/07e42b07-620f-40d0-a062-c15d3e22eb74' \
--header 'Content-Type: application/json' \
--data-raw '{
	"ID":"07e42b07-620f-40d0-a062-c15d3e22eb74",
	"Title": "office group",
	"Description": "this group has been modified",
	"Peers": [
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
DELETE localhost:8080/groups/07e42b07-620f-40d0-a062-c15d3e22eb74
```

```
The group with ID 2 has been deleted successfully
```
