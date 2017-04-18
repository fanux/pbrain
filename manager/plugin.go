package manager

import (
	"context"

	"github.com/Sirupsen/logrus"
)

//plugin commands define
const (
	PluginCommandCreate   = "plugin-create"
	PluginCommandStart    = "plugin-start"
	PluginCommandStop     = "plugin-stop"
	PluginCommandDestroy  = "plugin-destroy"
	PluginCommandOnaction = "plugin-strategy-onaction"
)

//Pluginer Define interface plugins need to inplement
type Pluginer interface {
	/*
		       plugins can get plugin info from context:

			   	if v := ctx.Value(k); v == nil {
				    fmt.Println("key not found:", k)
					return
				}
				if c,ok := v.(*PluginStrategy);ok{
					fmt.Println(c.Name, c.Kind, c.Status)
				}
	*/
	Create(context.Context) error
	Start(context.Context) error
	Stop(context.Context) error
	Destroy(context.Context) error

	//every strategy has "enable" and "disable" action, or other actions plugin define itself
	OnAction(ctx context.Context, name string, action string) error
}

var plugins map[string]Pluginer

//Regist is
func Regist(name string, p Pluginer) {
	if plugins == nil {
		plugins = make(map[string]Pluginer)
		logrus.Infof("plugner is nil, new map")
	}
	plugins[name] = p
}

//GetPlugin is
func GetPlugin(name string) Pluginer {
	p, ok := plugins[name]
	if !ok {
		return nil
	}
	return p
}
