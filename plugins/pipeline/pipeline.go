package pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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

type PipelineRule struct {
	Cron string

	Apps []common.AppScale
}

type CronFunc struct {
	Client common.Client
	Rule   PipelineRule
}

func (this CronFunc) Run() {
	log.Printf("Cron: [%s] time: [%s]", this.Rule.Cron,
		time.Now().Format("2016-01-02 03:04:05 PM"))

	this.Client.ScaleApps(this.Rule.Apps)
}

// Start plugin
func (this *Pipeline) Start() error {
	// init jobs
	this.Jobs = make(map[string]Job)

	for _, strategy := range this.Strategies {
		cron := cron.New()
		this.Jobs[strategy.Name] = Job{cron, strategy}
		if this.Plugin.Status == common.PLUGIN_ENABLE {
			cron.Start()

			this.EnableStrategy(strategy.Name)
		}
	}

	fmt.Println("start plugin")

	return nil
}

// Stop plugin
func (this *Pipeline) Stop() error {
	// disable all strategies
	fmt.Println("stop plugin")

	for _, v := range this.Jobs {
		v.Cron.Stop()
	}

	return nil
}

func (this *Pipeline) EnableStrategy(strategyName string) error {
	// start the strategy Cron
	job, ok := this.Jobs[strategyName]
	if !ok || job.Strategy.Status != common.STRATEGY_ENABLE {
		log.Printf("enable strategy failed: [%s]", strategyName)
		return errors.New("enable strategy failed")
	}

	rules := []PipelineRule{}

	err := json.Unmarshal([]byte(job.Strategy.Document), &rules)
	if err != nil {
		return errors.New("unmarshal document failed")
	}

	for _, rule := range rules {
		cronFunc := CronFunc{this.Client, rule}
		job.Cron.AddFunc(rule.Cron, cronFunc.Run)
	}

	fmt.Println("enable strategy: ", strategyName)
	return nil
}

func (this *Pipeline) DisableStrategy(strategyName string) error {
	// stop the strategy Cron
	fmt.Println("disable strategy: ", strategyName)

	job, ok := this.Jobs[strategyName]

	if ok && job.Strategy.Status == common.STRATEGY_DISABLE {
		job.Cron.Stop()
	}

	return nil
}

func (this *Pipeline) UpdateDocument(strategyName string) error {
	// stop the strategy Cron
	this.DisableStrategy(strategyName)
	// start the strategy Cron
	this.EnableStrategy(strategyName)

	fmt.Println("update document: ", strategyName)
	return nil
}
