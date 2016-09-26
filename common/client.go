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
	MetricalScaleApps(appscales []MetricalAppScale) ([]MetricalAppScaleHosts, error)
	/*
		{
		    "App":"hadoop:latest",
		    "Number":-2,             # 会释放2个实例
		    "MinNumber":1,
		    "Hosts":[
		        "192.168.0.2",       # 将释放哪些节点上的实例
		        "192.168.0.3",
		    ],
		    "ScaleUp":"ats:latest"    # 释放的实例用于启动哪个业务，业务不需要关心这个字段
		}
	*/
	MetricalScaleAppsAction(message InformScaleDownAppMessage) error
}

func NewClient(host, port string) Client {
	log.Printf("new api client info:%s%s", host, port)
	return Api{host, port}
}
