package state_system

const MAXIMUM_STATE_CHANGES = 16

type StateTree struct {
	Name         string       `json:"name"`
	Main         string       `json:"main"`
	StateList    []*GameState `json:"states"`
	stateMap     map[string]*GameState
	pendingState *GameState
	activeState  *GameState
}

func NewStateTree() *StateTree {
	tree := new(StateTree)

	tree.stateMap = make(map[string]*GameState)

	return tree
}

// Adds a new GameState instance to the state tree.
func (tree *StateTree) AddState(name string, state *GameState) bool {
	if state != nil {
		if tree.findGameState(name) == nil {
			// TODO: Rather than have a slice and a map, can we just have the one?
			tree.StateList = append(tree.StateList, state)
			tree.stateMap[name] = state
			return true
		}
	}

	return false
}

//! \brief Invoked when the state tree is ready for use and game systems may be prepared for processing.
//! \param initArgs [in] -
//!        Initialization information for use by the state tree.
func (tree *StateTree) OnInitialize(initArgs *InitArgs) {
	initArgs.StateTree = tree

	for _, state := range tree.stateMap {
		// We only pass the call onto root states, they will pass it onto their children
		if state.parent == nil {
			state.OnInitialize(initArgs)
		}
	}
}

//! \brief Invoked when the state tree is about to be removed from the running title.
func (tree *StateTree) OnDestroy() {
	for _, state := range tree.stateMap {
		// We only pass the call onto root states, they will pass it onto their children
		if state.parent == nil {
			state.OnDestroy()
		}
	}
}

//! \brief Called each frame the state tree should be processed.
//! \param updateArgs [in] -
//!        Details about the current frame being processed.
func (tree *StateTree) OnUpdate(updateArgs UpdateArgs) {
	if tree.activeState != nil {
		tree.activeState.OnUpdate(updateArgs)
	}

	tree.commitStateChange()
}

func (tree *StateTree) OnPostUpdate(updateArgs UpdateArgs) {
	if tree.activeState != nil {
		tree.activeState.OnPostUpdate(updateArgs)
	}

	tree.commitStateChange()
}

func (tree *StateTree) requestStateChange(name string) {
	state, exists := tree.stateMap[name]
	if exists {
		tree.pendingState = state
	}
	// TODO: Should we raise/log an error if the requested state was invalid?
}

func (tree *StateTree) commitStateChange() {
	changeCounter := 0
	for tree.pendingState != nil && changeCounter < MAXIMUM_STATE_CHANGES {
		pending := tree.pendingState

		changeCounter++
		tree.pendingState = nil

		if pending != tree.activeState {
			rootState := tree.findCommonAncestor(tree.activeState, pending)

			if tree.activeState != nil {
				tree.activeState.OnExit(rootState)
			}

			tree.activeState = pending
			pending.OnEnter(rootState)
			pending.OnPostEnter(rootState)
		}
	}
}

func (tree *StateTree) findGameState(name string) *GameState {
	state, exists := tree.stateMap[name]
	if exists {
		return state
	}

	return nil
}

func (tree *StateTree) findCommonAncestor(stateA *GameState, stateB *GameState) *GameState {
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
