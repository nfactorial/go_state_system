package state_system

import "testing"

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
