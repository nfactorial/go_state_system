package state_system

type SystemConstructor func() IGameSystem

// This interface provides support for creating game systems that have been registered with a specific type name.
type ISystemFactory interface {
	Create(typeName string) IGameSystem
}

// Factory object that allows creation of a game system based on the string it has been associated with.
type SystemFactory struct {
	systemMap map[string] SystemConstructor
}

func CreateSystemFactory() *SystemFactory {
	factory := new(SystemFactory)

	factory.systemMap = make(map[string] SystemConstructor)

	return factory
}

// Registers a new game system with the factory object. If the specified name is already associated with a system
// object this method return false. If the system is successfully registered this method returns true.
func (factory *SystemFactory) Register(typeName string, ctor SystemConstructor) bool {
	if ctor != nil {
		_, exists := factory.systemMap[typeName]
		if !exists {
			factory.systemMap[typeName] = ctor
			return true
		}
	}

	return false
}

// Removes a previously registered game system from the factory object. This method returns true if a game system
// was successfully unregistered otherwise it returns false.
func (factory *SystemFactory) Unregister(typeName string) bool {
	_, exists := factory.systemMap[typeName]
	if exists {
		delete(factory.systemMap, typeName)
		return true
	}

	return false
}

// Attempts to create a game system that has been associated with the specified string.
func (factory *SystemFactory) Create(typeName string) IGameSystem {
	ctor, exists := factory.systemMap[typeName]
	if exists {
		return ctor()
	}

	return nil
}

func (factory *SystemFactory) exists(typeName string) bool {
	_, exists := factory.systemMap[typeName]
	return exists
}
