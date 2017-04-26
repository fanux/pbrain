package manager

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
)

//Plugin is
type Plugin struct {
	Name        string
	Kind        string
	Status      string
	Description string
	Manual      string
	Spec        string
}

//Strategy is
type Strategy struct {
	Name     string
	Status   string
	Actions  []string
	Document string
}

//PluginStrategy is
type PluginStrategy struct {
	*Plugin
	Strategy map[string]*Strategy
}

//PluginMap is
type PluginMap struct {
	plugins map[string]*PluginStrategy
	mutex   sync.Mutex
}

//Command is
type Command struct {
	Name string
	//plugin name
	Pname string
	//strategy name
	Sname string
	//other argument, like action name
	Body string
}

var (
	pluginMap PluginMap
	pdb       Dber
)

//GetPlugin is
func (ps *PluginMap) GetPlugin(name string) *Plugin {
	temp, ok := ps.plugins[name]
	if ok {
		return temp.Plugin
	}
	logrus.Errorf("get plugin failed: %s", name)
	return nil

}

//GetStrategy is
func (ps *PluginMap) GetStrategy(pname string, name string) *Strategy {
	temp, ok := ps.plugins[pname]
	if !ok {
		logrus.Errorf("get plugin failed: %s", pname)
		return nil
	}
	stemp, ok := temp.Strategy[name]
	if !ok {
		logrus.Errorf("get plugin [%s] failed, strategy: %s", pname, name)
		return nil
	}
	return stemp
}

//SetPlugin is
func (ps *PluginMap) SetPlugin(name string, p Plugin) {
	ps.mutex.Lock()
	defer ps.metex.Unlock()

	ps.plugins[name].Plugin = &p
	//Save to db
	pdb.Save(ps.plugins[name])
}

//SetStrategy is
func (ps *PluginMap) SetStrategy(pname string, name string, s Strategy) {
	ps.mutex.Lock()
	defer ps.metex.Unlock()

	ps.plugins[pname].Strategy[name] = &s
	//Save to db
	pdb.Save(ps.plugins[name])
}

//StartManager is
func StartManager(opts Opts) {
	pdb = NewDb("", "")
	pluginMap.plugins = pdb.LoadAll()
	if plugins == nil {
		plugins = make(map[string]Pluginer)
	}
	events := make(map[string]chan Command)

	for k, v := range pluginMap.plugins {
		events[k] = make(chan Command)
		ctx := context.Background()
		ctx = context.WithValue(ctx, k, v)

		p = GetPlugin(k)
		go func(ctx context.Context, command chan Command, p Pluginer) {
			for {
				//TODO each gorutine communicate with ctx, set a chan into ctx
				cmd := <-command
				switch cmd.Name {
				case PluginCommandCreate:
				//	go p.Create(ctx)
				case PluginCommandStart:
					go p.Start(ctx)
				case PluginCommandStop:
				//	go p.Stop(ctx)
				case PluginCommandDestroy:
				//	go p.Destroy(ctx)
				case PluginCommandOnaction:
					go p.OnAction(ctx, cmd.Pname, cmd.Sname, cmd.Body)
				default:
					logrus.Errorf("Unknow command type: %s", cmd.Name)
				}
			}
		}(ctx, events[k], p)

		events[k] <- Command{PluginCommandStart, k, "", ""}

		for sname, s := range v.Strategy {
			if s.Status == StrategyActionEnable {
				events[k] <- Command{PluginCommandOnaction, k, sname, StrategyActionEnable}
			}
		}
	}

	//just for test
	http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
