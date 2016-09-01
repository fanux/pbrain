package pipeline

import (
	"github.com/fanux/pbrain/common"
	"github.com/robfig/cron"
)

type Job struct {
	Cron     *cron.Cron
	Strategy common.Strategy
}

type Pipeline struct {
	*common.PluginBase
	// the key is strategy name
	Jobs map[string]Job
}

// Start plugin
func (this *Pipeline) Start() error {
	this.Init()

	if this.Plugin.Status == PLUGIN_ENABLE {
	} else if this.Plugin.Status == PLUGIN_DISABLE {
	}
}

// Stop plugin
func (this *Pipeline) Stop() error {
	// disable all strategies
}

func (this *Pipeline) EnableStrategy(strategyName string) error {
	// start the strategy Cron
}

func (this *Pipeline) DisableStrategy(strategyName string) error {
	// stop the strategy Cron
}

func (this *Pipeline) UpdateDocument(strategyName string) error {
	// stop the strategy Cron
	// start the strategy Cron
}
