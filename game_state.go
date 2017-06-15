package state_system

type GameState struct {
	Name string `json:"name"`
	Parent string `json:"parent"`
	Children []string `json:"children"`
	Systems []string `json:"systems"`
	parentState *GameState
	childList []*GameState
	systemMap map[string] IGameSystem	// TODO: Should be a 'SystemRef' object array?
	systemList []IGameSystem
	updateList []IGameSystem
}

func NewGameState(name string) *GameState {
	state := new(GameState)

	state.Name = name
	state.systemMap = make(map[string] IGameSystem)

	return state
}

//! \brief Invoked when the state tree is ready for use and game systems may be prepared for processing.
//! \param initArgs [in] -
//!        Initialization information for use by the state tree.
func (state *GameState) OnInitialize(initArgs *InitArgs) {
	initArgs.GameState = state

	for _, system := range state.systemList {
		system.OnInitialize(initArgs)
	}

	for _, child := range state.childList {
		child.OnInitialize(initArgs)
	}
}

//! \brief Invoked when the game state is about to be removed from the running title.
func (state *GameState) OnDestroy() {
	// Child states must be destroyed before our own systems are destroyed
	for loop := len(state.childList) - 1; loop >= 0; loop-- {
		state.childList[loop].OnDestroy()
	}

	for loop := len(state.systemList) - 1; loop >= 0; loop-- {
		state.systemList[loop].OnDestroy()
	}
}

//! \brief Invoked when the game state becomes active within the running title.
//! \param root [in] -
//!        The game state at the root of the state switch, activation will not be passed up-to the root state.
func (state *GameState) OnEnter(root *GameState) {
	if state.parentState != nil && state.parentState != root {
		state.parentState.OnEnter(root)
	}

	for _, system := range state.systemList {
		system.OnActivate()
	}
}

//! \brief Invoked when the game state is no longer active within the running title.
//! \oaram root [in] -
//!        The game state at the root of the state switch, de-activation will not be passed up-to the root state.
func (state *GameState) OnExit(root *GameState) {
	for loop := len(state.systemList) - 1; loop >= 0; loop-- {
		state.systemList[loop].OnDeactivate()
	}

	if state.parentState != nil && state.parentState != root {
		state.parentState.OnExit(root)
	}
}

//! \brief Called each frame while the game state is active.
//! \param updateArgs [in] -
//!        Details about the current frame being processed.
func (state *GameState) OnUpdate(updateArgs UpdateArgs) {
	if state.parentState != nil {
		state.parentState.OnUpdate(updateArgs)
	}

	for _, system := range state.updateList {
		system.OnUpdate(updateArgs)
	}
}

func (state *GameState) FindSystem(name string) IGameSystem {
	system, exists := state.systemMap[name]
	if exists {
		return system
	}

	if state.parentState != nil {
		return state.parentState.FindSystem(name)
	}

	return nil
}

func (state *GameState) addSystem(name string, system IGameSystem) {
	_, exists := state.systemMap[name]
	if !exists && system != nil {
		state.systemMap[name] = system
		state.systemList = append(state.systemList, system)
		state.Systems = append(state.Systems, name)
	}
}

//! \brief Determines whether or not the specified state exists within our parent branch of the state tree.
//! \param state [in] -
//!        The state to be looked for within the parent hierarchy.
//! \return <em>True</em> if the supplied state exists within our parent hierarchy otherwise <em>false</em>.
func (state *GameState) checkParentHierarchy(root *GameState) bool {
	if state == root {
		return true
	}

	if state.parentState != nil {
		return state.parentState.checkParentHierarchy(root)
	}

	return false
}
