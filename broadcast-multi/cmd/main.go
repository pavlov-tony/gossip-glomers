package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type App struct {
	sync.RWMutex
	node     *maelstrom.Node
	store    []int
	topology map[string][]string
}

func main() {
	app := &App{
		node:  maelstrom.NewNode(),
		store: make([]int, 0),
	}
	app.node.Handle("read", app.ReadHandler)
	app.node.Handle("topology", app.TopologyHandler)
	app.node.Handle("broadcast", app.BroadcastHandler)
	if err := app.node.Run(); err != nil {
		log.Fatalf("Failed to run the node: %v", err)
	}
}

func (a *App) ReadHandler(msg maelstrom.Message) error {
	a.RLock()
	ids := make([]int, len(a.store))
	copy(ids, a.store)
	a.RUnlock()
	return a.node.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": ids,
	})
}

func (a *App) TopologyHandler(msg maelstrom.Message) error {
	t := struct {
		Topology map[string][]string
	}{}
	if err := json.Unmarshal(msg.Body, &t); err != nil {
		return err
	}
	a.Lock()
	a.topology = t.Topology
	a.Unlock()
	return a.node.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
}

func (a *App) BroadcastHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	m, ok := body["message"].(float64)
	if !ok {
		return fmt.Errorf("invalid message type")
	}
	a.Lock()
	a.store = append(a.store, int(m))
	a.Unlock()
	a.RLock()
	temp := a.topology[a.node.ID()]
	topology := make([]string, len(temp))
	copy(topology, temp)
	a.RUnlock()
	for _, node := range topology {
		a.node.RPC(node)
	}
	return a.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}
