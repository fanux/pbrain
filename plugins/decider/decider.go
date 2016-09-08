package decider

import "github.com/fanux/pbrain/common"

const (
	COMMAND_APP_METRICAL = "app-metrical"
)

type AppConf struct {
	App      string
	Priority int
	MinNum   int
	Spec     []struct {
		Metrical []int
		Number   int
	}
}

type AppInfo struct {
}

type Decider struct {
	*common.PluginBase

	AppInfo map[string]AppInfo
}

func (this *Decider) Start() error {
	common.RegistCommand(COMMAND_APP_METRICAL, this.OnScale)
}

func (this *Decider) Stop() error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
}

func (this *Decider) DisableStrategy(strategyName string) error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
}

func (this *Decider) EnableStrategy(strategyName string) error {
	common.RegistCommand(COMMAND_APP_METRICAL, this.OnScale)
}

func (this *Decider) UpdateDocument(strategyName string) error {
}

func (this *Decider) OnScale(command common.Command) error {
}
