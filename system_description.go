package state_system

type systemDescription struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Parameters map[string]SystemParameter `json:"params"`
	instance IGameSystem
}

// TODO: Should get SystemFactory from another package?
func (desc *systemDescription) create(factory ISystemFactory) {
	desc.instance = factory.Create(desc.Type)
}
