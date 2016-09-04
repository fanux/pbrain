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
	Cron *cron.Cron
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
		this.Jobs[strategy.Name] = Job{cron}
		if this.Plugin.Status == common.PLUGIN_ENABLE {
			//cron.Start()
			//this.EnableStrategy(strategy.Name)

			if strategy.Status == common.STRATEGY_ENABLE {
				cron.Start()
			}

			rules := []PipelineRule{}

			err := json.Unmarshal([]byte(strategy.Document), &rules)
			if err != nil {
				return errors.New("unmarshal document failed")
			}

			for _, rule := range rules {
				cronFunc := CronFunc{this.Client, rule}
				cron.AddFunc(rule.Cron, cronFunc.Run)
			}

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
	strategy := this.GetStrategy(strategyName)

	if !ok || strategy.Status != common.STRATEGY_ENABLE {
		log.Printf("enable strategy failed: [%s] status [%s]", strategyName, strategy.Status)
		return errors.New("enable strategy failed")
	}
	job.Cron.Start()

	/*
		rules := []PipelineRule{}

		err := json.Unmarshal([]byte(strategy.Document), &rules)
		if err != nil {
			return errors.New("unmarshal document failed")
		}

		for _, rule := range rules {
			cronFunc := CronFunc{this.Client, rule}
			job.Cron.AddFunc(rule.Cron, cronFunc.Run)
		}

		fmt.Println("enable strategy: ", strategyName)
	*/
	return nil
}

func (this *Pipeline) DisableStrategy(strategyName string) error {
	// stop the strategy Cron
	job, ok := this.Jobs[strategyName]
	strategy := this.GetStrategy(strategyName)

	if ok && strategy.Status == common.STRATEGY_DISABLE {
		job.Cron.Stop()
		fmt.Println("disable strategy: ", strategyName)
	} else {
		fmt.Println("disable strategy failed: ", strategyName, " ", strategy.Status)
	}

	return nil
}

func (this *Pipeline) UpdateDocument(strategyName string) error {
	/*
		// stop the strategy Cron
		this.DisableStrategy(strategyName)
		// start the strategy Cron
		this.EnableStrategy(strategyName)
	*/

	strategy := this.GetStrategy(strategyName)
	job, ok := this.Jobs[strategyName]
	if ok {
		job.Cron.Stop()
	} else {
		return errors.New("get strategy job error")
	}

	// new a new cron replace the old one
	cron := cron.New()
	this.Jobs[strategy.Name] = Job{cron}

	if strategy.Status == common.STRATEGY_ENABLE {
		cron.Start()
	}
	rules := []PipelineRule{}

	err := json.Unmarshal([]byte(strategy.Document), &rules)
	if err != nil {
		return errors.New("unmarshal document failed")
	}

	for _, rule := range rules {
		cronFunc := CronFunc{this.Client, rule}
		cron.AddFunc(rule.Cron, cronFunc.Run)
	}

	fmt.Println("update document: ", strategyName)
	return nil
}
