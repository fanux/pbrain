package pipeline

import (
	"fmt"

	"github.com/fanux/pbrain/common"
	"github.com/robfig/cron"
)

const PLUGIN_NAME = "plugin_pipeline"

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
	for _, strategy := range this.Strategies {
	}

	if this.Plugin.Status == common.PLUGIN_ENABLE {
	} else if this.Plugin.Status == common.PLUGIN_DISABLE {
	}
	return nil
}

// Stop plugin
func (this *Pipeline) Stop() error {
	// disable all strategies
	fmt.Println("stop plugin")
	return nil
}

func (this *Pipeline) EnableStrategy(strategyName string) error {
	// start the strategy Cron
	fmt.Println("enable strategy: ", strategyName)
	return nil
}

func (this *Pipeline) DisableStrategy(strategyName string) error {
	// stop the strategy Cron
	fmt.Println("disable strategy: ", strategyName)
	return nil
}

func (this *Pipeline) UpdateDocument(strategyName string) error {
	// stop the strategy Cron
	// start the strategy Cron
	fmt.Println("update document: ", strategyName)
	return nil
}
