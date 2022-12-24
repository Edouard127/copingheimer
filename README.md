# Copingheimer

Copingheimer is a simple tool for coping with the fact that you can't use copenheimer :3

You can make huge bot farms that will scan the internet to find minecraft servers\
All the workers are connected to a central server that will handle the requests and synchronise the workers

## Features

- [x] Cope
- [x] Worker network
- [x] Concurrency

## Usage

### Start the central server
```bash
$ go run main.go --type="server" --mongo="mongodb://localhost:27017"
```

### Start a worker
```bash
$ go run main.go --type="worker" --node="192.168.0.200"
```

## Dashboard (WIP)
The dashboard is available at http://node_ip:80/
If you are a participant, you can get access to the servers scanned by the users of
the node at http://node_ip:80/servers

## Reward system (WIP)
The reward system is based on the number of servers scanned by the users of the node\
You get access to 2 servers scanned by server you scanned

## Syncing
All workers have their own offset of IP to scan\
The node will send them the offset when the worker sends the packet SPacketLogin with the number of instances
```go
packet.Marshal(
	packet.SPacketLogin,
	int32(w.Arguments.Instances), 
)
```

If the worker sends a wrong packet, the worker will be kicked, and the node will update all the workers
offset via the packet SPacketUpdate


## Risks
- [ ] You might get banned from your ISP (Although I doubt it)

I am not responsible for any of the above or any other damages caused by this tool.

## Join the public server
If you want to join the public server, you can do so by running the following command:
```bash
$ go run main.go --type="worker" --node="coming soon"
```
