package state_system

import "testing"

const TEST_FACTORY_SYSTEM_NAME = "TestFactory"
const TEST_INVALID_SYSTEM_NAME = "InvalidSystem"

type FactoryTestMockGameSystem struct {
	counter int
}

func (system *FactoryTestMockGameSystem) OnDestroy() {

}

func (system *FactoryTestMockGameSystem) OnInitialize(initArgs *InitArgs) {
}

func (system *FactoryTestMockGameSystem) OnActivate() {
	system.counter++
}

func (system *FactoryTestMockGameSystem) OnDeactivate() {

}

func (system *FactoryTestMockGameSystem) OnUpdate(updateArgs UpdateArgs) {

}

func MockGameSystemConstructor() IGameSystem {
	return new(FactoryTestMockGameSystem)
}

func TestSystemFactory_Register(t *testing.T) {
	factory := CreateSystemFactory()

	if factory.Register(TEST_FACTORY_SYSTEM_NAME, nil) {
		t.Error("Unexpected successful registration of nil constructor.")
	}

	if !factory.Register(TEST_FACTORY_SYSTEM_NAME, MockGameSystemConstructor) {
		t.Error("Unable to register constructor.")
	}

	if factory.Register(TEST_FACTORY_SYSTEM_NAME, MockGameSystemConstructor) {
		t.Error("Should not be able to register multiple systems with a single key")
	}
}

func TestSystemFactory_Unregister(t *testing.T) {
	factory := CreateSystemFactory()

	if factory.Unregister(TEST_FACTORY_SYSTEM_NAME) {
		t.Error("Unregister succeeded with empty factory.")
	}

	factory.Register(TEST_FACTORY_SYSTEM_NAME, MockGameSystemConstructor)
	if !factory.Unregister(TEST_FACTORY_SYSTEM_NAME) {
		t.Error("Unregister failed when attempting to remove registered system.")
	}
	if factory.Unregister(TEST_FACTORY_SYSTEM_NAME) {
		t.Error("Unregister succeeded when removing the same system twice.")
	}
}

func TestSystemFactory_exists(t *testing.T) {
	factory := CreateSystemFactory()

	if factory.exists(TEST_FACTORY_SYSTEM_NAME) {
		t.Error("Newly constructed factory reported presense of non-existant object.")
	}

	factory.Register(TEST_FACTORY_SYSTEM_NAME, MockGameSystemConstructor)

	if !factory.exists(TEST_FACTORY_SYSTEM_NAME) {
		t.Error("factory.exists return false after system object was registered.")
	}

	factory.Unregister(TEST_FACTORY_SYSTEM_NAME)

	if factory.exists(TEST_FACTORY_SYSTEM_NAME) {
		t.Error("factory.exists returned true after system object was removed.")
	}
}

func TestSystemFactory_Create(t *testing.T) {
	factory := CreateSystemFactory()

	system := factory.Create(TEST_FACTORY_SYSTEM_NAME)
	if system != nil {
		t.Error("System factory successfully created object before it was registered.")
	}

	factory.Register(TEST_FACTORY_SYSTEM_NAME, MockGameSystemConstructor)

	system = factory.Create(TEST_INVALID_SYSTEM_NAME)
	if system != nil {
		t.Error("Expected invalid system name creation to fail, but it succeeded.")
	}

	system = factory.Create(TEST_FACTORY_SYSTEM_NAME)
	if system == nil {
		t.Error("Failed to create factory object after registration.")
	}

	testInstance, ok := system.(*FactoryTestMockGameSystem)
	if !ok {
		t.Error("Returned object from factory was not of expected type.")
	}

	if testInstance.counter != 0 {
		t.Error("Created instance contained an invalid counter value.")
	}
	// Invoke our counter method, so we can check they are the same object
	system.OnActivate()

	if testInstance.counter != 1 {
		t.Error("Counter was an unexpected value, expected 1 but found", testInstance.counter)
	}
}
