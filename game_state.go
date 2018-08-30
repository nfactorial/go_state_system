package state_system

type GameSystemDesc struct {
	Name       string                 `json:"name"`
	SystemType string                 `json:"type"`
	Params     map[string]interface{} `json:"params"`
}

type GameStateDesc struct {
	Name     string           `json:"name"`
	Children []string         `json:"children"`
	Systems  []GameSystemDesc `json:"systems"`
}

type StateTreeDesc struct {
	Name   string          `json:"name"`
	Main   string          `json:"main"`
	States []GameStateDesc `json:"states"`
}

type GameState struct {
	parent     *GameState
	name       string
	childList  []*GameState
	systemMap  map[string]IGameSystem
	systemList []IGameSystem
	updateList []IGameSystem
}

func NewGameState(name string) *GameState {
	state := new(GameState)

	state.name = name
	state.systemMap = make(map[string]IGameSystem)

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
	if state.parent != nil && state.parent != root {
		state.parent.OnEnter(root)
	}

	for _, system := range state.systemList {
		system.OnActivate()
	}
}

//! \brief Invoked after the game becomes active and has completed the OnEnter loop.
//! \param root [in] -
//!		   The game state at the root of the state switch, activation will not be passed up-to the root state.
func (state *GameState) OnPostEnter(root *GameState) {
	if state.parent != nil && state.parent != root {
		state.parent.OnPostEnter(root)
	}

	for _, system := range state.systemList {
		system.OnPostActivate()
	}
}

//! \brief Invoked when the game state is no longer active within the running title.
//! \oaram root [in] -
//!        The game state at the root of the state switch, de-activation will not be passed up-to the root state.
func (state *GameState) OnExit(root *GameState) {
	for loop := len(state.systemList) - 1; loop >= 0; loop-- {
		state.systemList[loop].OnDeactivate()
	}

	if state.parent != nil && state.parent != root {
		state.parent.OnExit(root)
	}
}

//! \brief Called each frame while the game state is active.
//! \param updateArgs [in] -
//!        Details about the current frame being processed.
func (state *GameState) OnUpdate(updateArgs UpdateArgs) {
	if state.parent != nil {
		state.parent.OnUpdate(updateArgs)
	}

	for _, system := range state.updateList {
		system.OnUpdate(updateArgs)
	}
}

//! \brief	Called each frame while the game state is active and after the main OnUpdate method has completed.
//! \param	updateArgs [in] -
//!			Details about the current frame being processed.
func (state *GameState) OnPostUpdate(updateArgs UpdateArgs) {
	if state.parent != nil {
		state.parent.OnPostUpdate(updateArgs)
	}

	for _, system := range state.updateList {
		system.OnPostUpdate(updateArgs)
	}
}

//! \brief	Retrieves an IGameSystem object that is associated with the specified name.
//!			Game systems can only be found within the current game state or in the parent hierarchy.
//! \param  name [in] - The name of the IGameSystem instance to be retrieved.
//! \return The IGameSystem associated with the specified name or nil if one could not be found.
func (state *GameState) FindSystem(name string) IGameSystem {
	system, exists := state.systemMap[name]
	if exists {
		return system
	}

	if state.parent != nil {
		return state.parent.FindSystem(name)
	}

	return nil
}

func (state *GameState) addSystem(name string, system IGameSystem) {
	_, exists := state.systemMap[name]
	if !exists && system != nil {
		state.systemMap[name] = system
		state.systemList = append(state.systemList, system)
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

	if state.parent != nil {
		return state.parent.checkParentHierarchy(root)
	}

	return false
}
