package main

import (
	"path/filepath"
	"strconv"
	"testing"

	cfacade "github.com/actorgo-game/actorgo/facade"
	cprofile "github.com/actorgo-game/actorgo/profile"
)

func TestDemoGORMNodeConfig(t *testing.T) {
	config, err := cprofile.LoadFile(filepath.Join("..", "config"), "demo-gorm.json")
	if err != nil {
		t.Fatalf("load demo-gorm config: %v", err)
	}

	numericNodeID, err := cfacade.GenNodeIdByStr("0.0.5.1")
	if err != nil {
		t.Fatalf("generate node ID: %v", err)
	}
	node, err := cprofile.GetNodeWithConfig(config, strconv.FormatUint(numericNodeID, 10), "5")
	if err != nil {
		t.Fatalf("resolve demo-gorm node: %v", err)
	}
	if !node.Enabled() {
		t.Fatal("demo-gorm node must be enabled")
	}
	if got := node.Settings().GetConfig("db_id_list").GetString("center_db_id"); got != "center_db_1" {
		t.Fatalf("center_db_id = %q, want center_db_1", got)
	}
}
