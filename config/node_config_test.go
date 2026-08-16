package config_test

import (
	"strconv"
	"testing"

	cfacade "github.com/actorgo-game/actorgo/facade"
	cprofile "github.com/actorgo-game/actorgo/profile"
)

func TestExampleNodeConfigs(t *testing.T) {
	tests := []struct {
		file     string
		nodeID   string
		nodeType string
	}{
		{file: "test.json", nodeID: "0.0.1.1", nodeType: "1"},
		{file: "test.json", nodeID: "0.0.3.1", nodeType: "3"},
		{file: "test.json", nodeID: "0.0.5.1", nodeType: "5"},
		{file: "test-discovery.json", nodeID: "0.0.1.1", nodeType: "1"},
		{file: "test-discovery.json", nodeID: "0.0.5.1", nodeType: "5"},
		{file: "test-discovery.json", nodeID: "0.0.5.2", nodeType: "5"},
		{file: "demo-chat.json", nodeID: "0.0.1.1", nodeType: "1"},
		{file: "demo-gorm.json", nodeID: "0.0.5.1", nodeType: "5"},
		{file: "test-http.json", nodeID: "0.0.3.1", nodeType: "3"},
	}

	for _, test := range tests {
		t.Run(test.file+"/"+test.nodeID, func(t *testing.T) {
			profile, err := cprofile.LoadFile(".", test.file)
			if err != nil {
				t.Fatalf("load profile: %v", err)
			}
			numericNodeID, err := cfacade.GenNodeIdByStr(test.nodeID)
			if err != nil {
				t.Fatalf("generate node ID: %v", err)
			}
			node, err := cprofile.GetNodeWithConfig(
				profile,
				strconv.FormatUint(numericNodeID, 10),
				test.nodeType,
			)
			if err != nil {
				t.Fatalf("resolve node: %v", err)
			}
			if !node.Enabled() {
				t.Fatal("node must be enabled")
			}
		})
	}
}
