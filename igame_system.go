package state_system

type IGameSystem interface {
	OnDestroy()
	OnInitialize(initArgs *InitArgs)
	OnActivate()
	OnPostActivate()
	OnDeactivate()
	OnUpdate(updateArgs UpdateArgs)
	OnPostUpdate(updateArgs UpdateArgs)
}
