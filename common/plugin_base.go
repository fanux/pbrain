package common

type PluginBase struct {
	Client     Client
	Plugin     Plugin
	Strategies []Strategy
}

func (this *PluginBase) Init() {
	// get plugin and strategy from plugin manager
	// new a client
}
