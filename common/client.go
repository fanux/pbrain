package common

import "log"

type Client interface {
	GetPluginInfo(pluginName string) (Plugin, error)
	GetPluginStrategies(pluginName string) ([]Strategy, error)

	/*
		[
			{
				"App":"ats",
				"Number":20,
			},
			{
				"App":"hadoop",
				"Number":10
			}
		]
	*/
	ScaleApps(appscale []AppScale) error
	/*
		[
			{
				"App":"ats:latest",
				"Number":5,
				"MinNum":2
			},
			{
				"App":"hadoop:latest",
				"Number":-2,
				"MinNum":1
			},
			{
				"App":"redis:latest",
				"Number":-3,
				"MinNum":3
			}
		]
	*/
	MetricalScaleApps(appscale []MetricalAppScale) ([]MetricalAppScaleHosts, error)
}

func NewClient(host, port string) Client {
	log.Printf("new api client info:%s%s", host, port)
	return Api{host, port}
}
