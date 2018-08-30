package state_system

import (
	"encoding/json"
	"testing"
)

func TestNewStateTree(t *testing.T) {
	tree := NewStateTree()

	if tree.pendingState != nil {
		t.Error("Pending state was not nil!")
	}

	if tree.activeState != nil {
		t.Error("Active state was not nil!")
	}

	if tree.stateMap == nil {
		t.Error("Tree did not contain a valid Map.")
	}
}

func TestReadJson(t *testing.T) {
	var jsonBlob = []byte(`{"name": "game", "main": "game_start", "states": []}`)

	var stateDesc StateTreeDesc

	json.Unmarshal(jsonBlob, &stateDesc)

	if stateDesc.Name != "game" {
		t.Error("State tree name was incorrect.")
	}

	if stateDesc.Main != "game_start" {
		t.Error("State tree main entry was incorrect.")
	}

	if len(stateDesc.States) != 0 {
		t.Error("State tree did not contain the expected states.")
	}
}
