package common

import "fmt"

type PluginBase struct {
	Client     Client
	Plugin     Plugin
	Strategies []Strategy
}

func (this *PluginBase) Init() {
	// get plugin and strategy from plugin manager
	// new a client
	this.Plugin, _ = this.Client.GetPluginInfo(this.Plugin.Name)
	this.Strategies, _ = this.Client.GetPluginStrategies(this.Plugin.Name)
	fmt.Printf("init plugin [%s] strategies [%s]\n", this.Plugin, this.Strategies)
}

func (this *PluginBase) GetPluginName() string {
	fmt.Println("get plugin name: ", this.Plugin.Name)
	return this.Plugin.Name
}

func (this *PluginBase) GetStrategy(strategyName string) Strategy {
	for _, strategy := range this.Strategies {
		if strategy.Name == strategyName {
			return strategy
		}
	}

	return Strategy{}
}

// using the plugin manager host and port
func GetBasePlugin(host, port, pluginName string) *PluginBase {
	pluginBase := PluginBase{}

	client := NewClient(host, port)
	pluginBase.Client = client
	pluginBase.Plugin.Name = pluginName
	pluginBase.Init()

	return &pluginBase
}
