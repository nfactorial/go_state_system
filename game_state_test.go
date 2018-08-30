package state_system

import "testing"

const (
	TEST_NAME        = "Test"
	TEST_SYSTEM_NAME = "TestSystem"
)

// This object implements the IGameSystem interface and counts how often each method is invoked. We use this to
// verify the game state forwards the calls onto the appropriate object at the expected times.
type MockSystem struct {
	IGameSystem
	destroyCount    int
	initCount       int
	updateCount     int
	activateCount   int
	deactivateCount int
}

func (system *MockSystem) OnDestroy() {
	system.destroyCount++
}

func (system *MockSystem) OnInitialize(initArgs *InitArgs) {
	system.initCount++
}

func (system *MockSystem) OnUpdate(updateArgs UpdateArgs) {
	system.updateCount++
}

func (system *MockSystem) OnActivate() {
	system.activateCount++
}

func (system *MockSystem) OnDeactivate() {
	system.deactivateCount++
}

func TestNewGameState(t *testing.T) {
	state := NewGameState(TEST_NAME)

	if state.parent != nil {
		t.Error("State contained an unexpected parent reference.")
	}

	if state.name != TEST_NAME {
		t.Error("Test game state contained an unexpected name", state.name)
	}
}

func TestGameState_addSystem(t *testing.T) {
	state := NewGameState(TEST_NAME)
	system := new(MockSystem)

	state.addSystem(TEST_SYSTEM_NAME, system)

	if state.FindSystem(TEST_SYSTEM_NAME) != system {
		t.Error("System could not be found after adding to game state.")
	}
}

func TestGameState_OnInitialize(t *testing.T) {
	state := NewGameState(TEST_NAME)
	system := new(MockSystem)
	initArgs := new(InitArgs)

	state.addSystem(TEST_SYSTEM_NAME, system)

	state.OnInitialize(initArgs)

	if system.initCount != 1 {
		t.Error("Expected initCount to be 1, was", system.initCount)
	}

	if system.destroyCount != 0 {
		t.Error("Expected destroyCount to be 0, was", system.destroyCount)
	}

	if system.destroyCount != 0 {
		t.Error("Expected updateCount to be 0, was", system.updateCount)
	}

	if system.destroyCount != 0 {
		t.Error("Expected activateCount to be 0, was", system.activateCount)
	}

	if system.destroyCount != 0 {
		t.Error("Expected deactivateCount to be 0, was", system.deactivateCount)
	}
}

func TestGameState_OnDestroy(t *testing.T) {
	state := NewGameState(TEST_NAME)
	system := new(MockSystem)

	state.addSystem(TEST_SYSTEM_NAME, system)

	state.OnDestroy()

	if system.destroyCount != 1 {
		t.Error("Expected destroyCount to be 1, was", system.destroyCount)
	}

	if system.initCount != 0 {
		t.Error("Expected initCount to be 0, was", system.initCount)
	}

	if system.updateCount != 0 {
		t.Error("Expected updateCount to be 0, was", system.updateCount)
	}

	if system.activateCount != 0 {
		t.Error("Expected activateCount to be 0, was", system.activateCount)
	}

	if system.deactivateCount != 0 {
		t.Error("Expected deactivateCount to be 0, was", system.deactivateCount)
	}
}
