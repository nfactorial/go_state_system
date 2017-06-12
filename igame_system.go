package state_system

type IGameSystem interface {
	OnDestroy()
	OnInitialize(initArgs *InitArgs)
	OnActivate()
	OnDeactivate()
	OnUpdate(updateArgs UpdateArgs)
}
