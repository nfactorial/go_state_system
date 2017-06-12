package state_system

type GameState struct {
	parent *GameState
	childList []GameState
	systemList []IGameSystem
	updateList []IGameSystem
}

//! \brief Invoked when the state tree is ready for use and game systems may be prepared for processing.
//! \param initArgs [in] -
//!        Initialization information for use by the state tree.
func (d *GameState) OnInitialize(initArgs InitArgs) {
	initArgs.GameState = d

	for _, system := range d.systemList {
		system.OnInitialize(initArgs)
	}

	for _, child := range d.childList {
		child.OnInitialize(initArgs)
	}
}

//! \brief Invoked when the game state is about to be removed from the running title.
func (d *GameState) OnDestroy() {
	// Child states must be destroyed before our own systems are destroyed
	for loop := len(d.childList) - 1; loop >= 0; loop-- {
		d.childList[loop].OnDestroy()
	}

	for loop := len(d.systemList) - 1; loop >= 0; loop-- {
		d.systemList[loop].OnDestroy()
	}
}

//! \brief Invoked when the game state becomes active within the running title.
//! \param root [in] -
//!        The game state at the root of the state switch, activation will not be passed up-to the root state.
func (d *GameState) OnEnter(root *GameState) {
	if d.parent != nil && d.parent != root {
		d.parent.OnEnter(root)
	}

	for _, system := range d.systemList {
		system.OnActivate()
	}
}

//! \brief Invoked when the game state is no longer active within the running title.
//! \oaram root [in] -
//!        The game state at the root of the state switch, de-activation will not be passed up-to the root state.
func (d *GameState) OnExit(root *GameState) {
	for loop := len(d.systemList) - 1; loop >= 0; loop-- {
		d.systemList[loop].OnDeactivate()
	}

	if d.parent != nil && d.parent != root {
		d.parent.OnExit(root)
	}
}

//! \brief Called each frame while the game state is active.
//! \param updateArgs [in] -
//!        Details about the current frame being processed.
func (d *GameState) OnUpdate(updateArgs UpdateArgs) {
	if d.parent != nil {
		d.parent.OnUpdate(updateArgs)
	}

	for _, system := range d.updateList {
		system.OnUpdate(updateArgs)
	}
}

//! \brief Determines whether or not the specified state exists within our parent branch of the state tree.
//! \param state [in] -
//!        The state to be looked for within the parent hierarchy.
//! \return <em>True</em> if the supplied state exists within our parent hierarchy otherwise <em>false</em>.
func (d *GameState) checkParentHierarchy(root *GameState) bool {
	if d == root {
		return true
	}

	if d.parent != nil {
		return d.parent.checkParentHierarchy(root)
	}

	return false
}
