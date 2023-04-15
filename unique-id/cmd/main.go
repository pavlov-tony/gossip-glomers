package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	var counter uint64
	n.Handle("generate", func(msg maelstrom.Message) error {
		body := map[string]any{}
		body["type"] = "generate_ok"
		buf := bytes.NewBuffer(make([]byte, 0, 136))
		if err := binary.Write(buf, binary.LittleEndian, time.Now().UnixNano()); err != nil {
			return err
		}
		ui64, err := strconv.ParseUint(strings.TrimPrefix(n.ID(), "n"), 10, 8)
		if err != nil {
			return err
		}
		nodeID := uint8(ui64)
		if err := binary.Write(buf, binary.LittleEndian, nodeID); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.LittleEndian, atomic.AddUint64(&counter, 1)); err != nil {
			return err
		}
		body["id"] = buf.Bytes()
		return n.Reply(msg, body)
	})
	if err := n.Run(); err != nil {
		log.Fatalf("Failed to run the node: %v", err)
	}
}
