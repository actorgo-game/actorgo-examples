package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type nodeConfigRow struct {
	NodeID  any    `json:"node_id"`
	Address string `json:"address"`
}

func TestClusterNodeTypes(t *testing.T) {
	configBytes, err := os.ReadFile(filepath.Join("..", "..", "..", "config", "demo-cluster.json"))
	if err != nil {
		t.Fatal(err)
	}

	var config struct {
		Cluster struct {
			NATS struct {
				MasterNodeID string `json:"master_node_id"`
			} `json:"nats"`
		} `json:"cluster"`
		Node map[string][]nodeConfigRow `json:"node"`
	}
	if err = json.Unmarshal(configBytes, &config); err != nil {
		t.Fatal(err)
	}

	if config.Cluster.NATS.MasterNodeID != "262145" {
		t.Fatalf("unexpected master_node_id: %s", config.Cluster.NATS.MasterNodeID)
	}

	for _, nodeType := range []string{"1", "2", "3", "4", "5"} {
		t.Run(nodeType, func(t *testing.T) {
			rows, ok := config.Node[nodeType]
			if !ok || len(rows) == 0 {
				t.Fatalf("node type %s is not configured", nodeType)
			}
			if rows[0].NodeID != nil {
				t.Fatalf("node type %s should use its single config without node_id", nodeType)
			}
		})
	}
	if got := config.Node["3"][0].Address; got != ":8081" {
		t.Fatalf("web address = %q, want %q", got, ":8081")
	}
}

func TestAreaServerNodeIDMapping(t *testing.T) {
	configBytes, err := os.ReadFile(filepath.Join("..", "..", "..", "config", "data", "areaServerConfig.json"))
	if err != nil {
		t.Fatal(err)
	}

	var rows []AreaServerRow
	if err = json.Unmarshal(configBytes, &rows); err != nil {
		t.Fatal(err)
	}

	expected := map[int32]string{
		10001: encodeNodeID(1, 1, 5, 1),
		10002: encodeNodeID(1, 2, 5, 1),
		10003: encodeNodeID(1, 3, 5, 1),
		10004: encodeNodeID(1, 4, 5, 1),
	}
	for _, row := range rows {
		want, ok := expected[row.ServerId]
		if !ok {
			t.Fatalf("unexpected serverId: %d", row.ServerId)
		}
		if row.NodeID != want {
			t.Fatalf("serverId %d: nodeId = %s, want %s", row.ServerId, row.NodeID, want)
		}
		delete(expected, row.ServerId)
	}
	if len(expected) != 0 {
		t.Fatalf("missing server mappings: %v", expected)
	}
}

func encodeNodeID(bigWorldID, worldID, nodeType, nodeInst uint64) string {
	const (
		bigWorldIDMask = uint64(0x3FF)
		commonMask     = uint64(0x3FFFF)
	)
	nodeID := (bigWorldID&bigWorldIDMask)<<54 |
		(worldID&commonMask)<<36 |
		(nodeType&commonMask)<<18 |
		(nodeInst & commonMask)
	return fmt.Sprint(nodeID)
}
