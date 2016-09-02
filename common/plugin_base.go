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
	fmt.Println("init plugin")
}

func (this *PluginBase) GetPluginName() string {
	fmt.Println("get plugin name: ", this.Plugin.Name)
	return this.Plugin.Name
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
