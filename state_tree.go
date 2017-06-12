package state_system

const MAXIMUM_STATE_CHANGES = 16

type StateTree struct {
	stateMap map[string] *GameState
	pendingState *GameState
	activeState *GameState
}

func NewStateTree() *StateTree {
	tree := new(StateTree)

	tree.stateMap = make(map[string] *GameState)

	return tree
}

//! \brief Invoked when the state tree is ready for use and game systems may be prepared for processing.
//! \param initArgs [in] -
//!        Initialization information for use by the state tree.
func (d *StateTree) OnInitialize(initArgs InitArgs) {
	initArgs.StateTree = d

	for _, state := range d.stateMap {
		// We only pass the call onto root states, they will pass it onto their children
		if state.parent == nil {
			state.OnInitialize(initArgs)
		}
	}
}

//! \brief Invoked when the state tree is about to be removed from the running title.
func (d *StateTree) OnDestroy() {
	for _, state := range d.stateMap {
		// We only pass the call onto root states, they will pass it onto their children
		if state.parent == nil {
			state.OnDestroy()
		}
	}
}

//! \brief Called each frame the state tree should be processed.
//! \param updateArgs [in] -
//!        Details about the current frame being processed.
func (d *StateTree) OnUpdate(updateArgs UpdateArgs) {
	if d.activeState != nil {
		d.activeState.OnUpdate(updateArgs)
	}

	d.commitStateChange()
}

func (d *StateTree) requestStateChange(name string) {
	state, exists := d.stateMap[name]
	if exists {
		d.pendingState = state
	}
	// TODO: Should we raise/log an error if the requested state was invalid?
}

func (d *StateTree) commitStateChange() {
	changeCounter := 0
	for d.pendingState != nil && changeCounter < MAXIMUM_STATE_CHANGES {
		var pending = d.pendingState
		changeCounter++
		d.pendingState = nil

		if pending != d.activeState {
			rootState := d.findCommonAncestor(d.activeState, pending)

			if d.activeState != nil {
				d.activeState.OnExit(rootState)
			}

			d.activeState = pending
			pending.OnEnter(rootState)
		}
	}
}

func (d *StateTree) findGameState(name string) *GameState {
	state, exists := d.stateMap[name]
	if exists {
		return state
	}

	return nil
}

func (d *StateTree) findCommonAncestor(stateA *GameState, stateB *GameState) *GameState {
	if stateA != nil && stateB != nil {
		if stateA == stateB {
			return stateA
		}

		scanA := stateA.parent
		scanB := stateB.parent

		for scanA != nil && scanB != nil {
			if stateB.checkParentHierarchy(scanA) {
				return scanA
			}

			if stateA.checkParentHierarchy(scanB) {
				return scanB
			}

			scanA = scanA.parent
			scanB = scanB.parent
		}
	}

	return nil
}
