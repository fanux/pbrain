package judge

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fanux/pbrain/manager"
)

//Document is
/*
 	{
		Items:[
			{
				"Metric":[60,80], # if 60 < metric <=80
				"Cmds":["scale -f app==online -n 5","scale -f app==offline -n -5"]
			},
			{
				"Metric":[40.60], # if 40 < metric <=60
				"Cmds":["scale -f app==online -n 5","scale -f app==offline -n -5"]
			},
			{
				"Metric":[20,40], # if 20 < metric <= 40
				"Cmds":["scale -f app==online -n 5","scale -f app==offline -n -5"]
			},
			{
				"Metric":[0,20], # if 0 < metric <= 20
				"Cmds":["scale -f app==online -n 5","scale -f app==offline -n -5"]
			},
		]
	}
*/
type Document struct {
	Items []struct {
		Metric []int
		Cmds   []string
	}
}

const judgePluginName = "plugin-judge"

var totalCount = 60

//PluginJudge is
type PluginJudge struct {
	*Document
	MetricsCount int
	TimesCount   int
	Particle     time.Second
}

func init() {
	manager.Regist(judgePluginName, &PluginJudge{Particle: 60})
}

//Judge return the excut command
func (j *PluginJudge) Judge() []string {
	metric := MetricsCount / TimesCount
	for _, i := range j.Items {
		if metric <= i.Metric[0] && metric > i.Metric[1] {
			return i.Cmds
		}
	}

	return []string{}
}

//Run is
func (j *PluginJudge) Run(p *manager.PluginStrategy, sname string) {
	if p.Status != manager.PluginStatusEnable {
		logrus.Infof("plugin not enable, run strategy failed: %s", p.Name)
		return
	}

	j.RunStrategy(p.Strategy[sname])
}

//RunStrategy is
func (j *PluginJudge) RunStrategy(s *manager.Strategy) {
	if s.Status != manager.StrategyStatusEnable {
		logrus.Infof("strategy not enable, run strategy failed: %s", s.Name)
		return
	}

	err := json.Unmarshal([]byte(s.Document), j.Document)
	if err != nil {
		logrus.Errorf("load strategy document error: %s", s.Name)
		return
	}

	c := NewCollector("")

	for {
		for ; j.TimesCount < totalCount; j.TimesCount++ {
			j.MetricsCount += c.Collect()
			time.Sleep(j.Particle)
		}
		j.TimesCount = 0
		j.MetricsCount = 0
		cmds := j.Judge()
		if len(cmds) != 0 {
			if err := Excute(cmds); err != nil {
				logrus.Errorf("Excute cmds [ %s ] error: %s", cmds, err)
			}
		}
	}
}

//Create is
func (j *PluginJudge) Create(ctx context.Context) error {
}

//Start is
func (j *PluginJudge) Start(ctx context.Context) error {
}

//Stop is
func (j *PluginJudge) Stop(ctx context.Context) error {
}

//Destroy is
func (j *PluginJudge) Destroy(ctx context.Context) error {
}

//OnAction is
func (j *PluginJudge) OnAction(ctx context.Context, pname string, sname string, action string) error {
	switch action {
	case manager.StrategyActionEnable:
		if v := ctx.Value(pname); v != nil {
			if p, ok := v.(*manager.PluginStrategy); ok {
				j.Run(p, sname)
			} else {
				return errors.New("assert plugin interface error")
			}
		} else {
			return errors.New("get plugin from context failed")
		}
	case manager.StrategyActionDisable:
		logrus.Debug("disable stratetegy todo")
	default:
		logrus.Infof("action not found")
	}
}
