package state_system

import "testing"

const TEST_FACTORY_SYSTEM_NAME = "TestFactory"

type FactoryTestMockGameSystem struct {
	counter int
}

func (system *FactoryTestMockGameSystem) OnDestroy() {

}

func (system *FactoryTestMockGameSystem) OnInitialize(initArgs *InitArgs) {
	system.counter++
}

func (system *FactoryTestMockGameSystem) OnActivate() {

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

func TestSystemFactory_Create(t *testing.T) {

}
