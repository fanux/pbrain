package decider

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/fanux/pbrain/common"
)

const (
	COMMAND_APP_METRICAL = "app-metrical"

	DEFAULT_METRICAL = 50

	PLUGIN_NAME = "plugin_decider"
)

type AppMetrical struct {
	App      string
	Metrical int
}

type AppConf struct {
	App      string
	Image    string
	Priority int
	MinNum   int
	Spec     []struct {
		Metrical []int
		Number   int
	}
}

type AppInfo struct {
	AppConf         AppConf
	CurrentMetrical int
}

type Document struct {
	AppInfos map[string]AppInfo // key is app name
}

// this plugin only suport one strategy enable, many strategies
// enable at a same time may conflict
type Decider struct {
	*common.PluginBase

	Documents map[string]Document // key is strategy name
}

type Filter struct {
	Selected []common.AppScale
}

// Only one strategy is better, strategy name shows app belong to witch strategy
func (this *Decider) getAppInfo(appName string) (*AppInfo, string) {
	for s, v := range this.Documents {
		if this.GetStrategy(s).Status != common.STRATEGY_ENABLE {
			continue
		}

		for k, appInfo := range v.AppInfos {
			if k == appName {
				return &appInfo, s
			}
		}
	}

	return nil, ""
}

func (this *Decider) initDocuments() {
	this.Documents = make(map[string]Document)

	document := []AppConf{}

	for _, strategy := range this.Strategies {
		err := json.Unmarshal([]byte(strategy.Document), &document)
		if err != nil {
			errors.New("json umarshal error")
			return
		}

		documentMap := Document{make(map[string]AppInfo)}

		for _, appConf := range document {
			appInfo := AppInfo{appConf, DEFAULT_METRICAL}
			documentMap.AppInfos[appConf.App] = appInfo
		}

		this.Documents[strategy.Name] = documentMap
	}
}

func (this *Decider) Start() error {
	this.initDocuments()
	common.RegistCommand(COMMAND_APP_METRICAL, this.OnScale)
	return nil
}

func (this *Decider) Stop() error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
	return nil
}

func (this *Decider) DisableStrategy(strategyName string) error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
	return nil
}

func (this *Decider) EnableStrategy(strategyName string) error {
	this.initDocuments()
	common.RegistCommand(COMMAND_APP_METRICAL, this.OnScale)
	return nil
}

func (this *Decider) UpdateDocument(strategyName string) error {
	this.initDocuments()
	return nil
}

// the metrical oppsite scale app number
func getScaleNumber(metrical int, appConf AppConf) (int, error) {
	for _, v := range appConf.Spec {
		if metrical >= v.Metrical[0] && metrical < v.Metrical[1] {
			return v.Number, nil
		}
	}

	return 0, errors.New("get scale number failed")
}

func (this *Decider) getAppScales(appInfo *AppInfo, scaleNumber int) []common.AppScale {
	// TODO decide scale witch apps
	// filter := Filter{[]common.AppScale{}}

	// TODO filter witch metrical is below balance metrical

	// TODO filter witch priority is lower

	return []common.AppScale{}
}

// one strategy one scale, not one plugin one scale
// or select one enable strategy witch include the APP config
func (this *Decider) OnScale(command common.Command) error {
	appMetrical := AppMetrical{}
	err := json.Unmarshal([]byte(command.Body), &appMetrical)
	if err != nil {
		return err
	}

	appInfo, _ := this.getAppInfo(appMetrical.App)
	if appInfo == nil {
		log.Printf("can not get app info, may be the strategy is disabled.")
		return nil
	}

	scaleNumber, err := getScaleNumber(appMetrical.Metrical, appInfo.AppConf)
	if err != nil {
		log.Printf("get scale number error [%s] : %s", appMetrical.App, err)
		return err
	}

	if scaleNumber <= 0 {
		// Do nothing except set the new metrical
		return nil
	} else if scaleNumber > 0 {
		appScales := this.getAppScales(appInfo, scaleNumber)
		fmt.Println(appScales)

		// TODO must know current app instance number

		// this.Client.ScaleApps(appScales)
	}

	defer func() {
		// update current metrical
		(*appInfo).CurrentMetrical = appMetrical.Metrical
	}()

	return nil
}
