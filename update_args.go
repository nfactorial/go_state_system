package state_system

type UpdateArgs struct {
	DeltaTime float32

	stateTree *StateTree
}

func (d *UpdateArgs) RequestStateChange(state string) {
	if d.stateTree != nil {
		d.stateTree.requestStateChange(state)
	}
}
