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
	COMMAND_APP_SCALE    = "app-scale"

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

/*
type Filter struct {
	Selected []common.AppScale
}
*/

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
	common.RegistCommand(COMMAND_APP_SCALE, this.ScaleAction)
	return nil
}

func (this *Decider) Stop() error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
	common.UnRegistCommand(COMMAND_APP_SCALE)
	return nil
}

func (this *Decider) DisableStrategy(strategyName string) error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
	common.UnRegistCommand(COMMAND_APP_SCALE)
	return nil
}

func (this *Decider) EnableStrategy(strategyName string) error {
	this.initDocuments()
	common.RegistCommand(COMMAND_APP_METRICAL, this.OnScale)
	common.RegistCommand(COMMAND_APP_SCALE, this.ScaleAction)
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

func makeScore(appInfo AppInfo) int {
	return 100 - appInfo.CurrentMetrical + 10*appInfo.AppConf.Priority
}

func getNeedFreeInstanceNum(score, totalScore, scaleNumber int) (n int) {
	n = scaleNumber * score / totalScore
	return (n + 1)
}

func (this *Decider) getAppScales(strategyName string, appInfo *AppInfo, scaleNumber int) []common.AppScale {
	// decide scale witch apps
	totalScore := 0
	filter := []common.AppScale{}
	appScore := make(map[string]int) // key is app name, and value is score

	document, ok := this.Documents[strategyName]
	if !ok {
		log.Printf("get strategy document failed: [%s]", strategyName)
		return nil
	}

	for _, v := range document.AppInfos {
		n, err := getScaleNumber(v.CurrentMetrical, v.AppConf)
		if err != nil {
			continue
		}

		if n < 0 {
			filter = append(filter, common.AppScale{v.AppConf.App, n})
		} else if n >= 0 && v.AppConf.Priority > appInfo.AppConf.Priority {
			// make a score
			s := makeScore(v)
			appScore[v.AppConf.App] = s

			totalScore += s
		}
	}

	var totalFree = 0

	for app, score := range appScore {
		num := getNeedFreeInstanceNum(score, totalScore, scaleNumber)
		log.Printf("get need free instance num: [%d] App: [%s]", num, app)
		if totalFree+num > scaleNumber {
			num = scaleNumber - totalFree
		}
		filter = append(filter, common.AppScale{app, -num})
		totalFree += num
	}

	return filter
}

func publishMessagesToApps(metricalAppScaleHosts []common.MetricalAppScaleHosts,
	scaleUpAppName string) {

	mq := common.NewMq()

	for _, v := range metricalAppScaleHosts {
		if v.Number < 0 {
			m := common.InformScaleDownAppMessage{v, scaleUpAppName}
			// each app subscribe itself
			mq.Publish(v.App, m)
			s, _ := json.Marshal(m)
			log.Println("Publish message to apps: ", string(s))
		}
	}
}

// one strategy one scale, not one plugin one scale
// or select one enable strategy witch include the APP config
func (this *Decider) OnScale(command common.Command) error {
	appMetrical := AppMetrical{}
	err := json.Unmarshal([]byte(command.Body), &appMetrical)
	if err != nil {
		return err
	}

	appInfo, strategyName := this.getAppInfo(appMetrical.App)
	if appInfo == nil {
		log.Printf("can not get app info, may be the strategy is disabled.")
		return nil
	}

	scaleNumber, err := getScaleNumber(appMetrical.Metrical, appInfo.AppConf)
	if err != nil {
		log.Printf("get scale number error [%s] : %s", appMetrical.App, err)
		return err
	}

	log.Printf("get scale number: %d", scaleNumber)

	if scaleNumber <= 0 {
		// Do nothing except set the new metrical
		return nil
	} else if scaleNumber > 0 {
		appScales := this.getAppScales(strategyName, appInfo, scaleNumber)
		fmt.Println("===need scale=====", appScales)

		metricalAppScales := []common.MetricalAppScale{}

		for _, app := range appScales {
			aInfo, _ := this.getAppInfo(app.App)
			metricalAppScales = append(metricalAppScales,
				common.MetricalAppScale{app.App, app.Number, aInfo.AppConf.MinNum})
		}

		metricalAppScales = append(metricalAppScales,
			common.MetricalAppScale{appInfo.AppConf.App,
				scaleNumber, appInfo.AppConf.MinNum})

		metricalAppScaleHosts, e := this.Client.MetricalScaleApps(metricalAppScales)

		if e != nil {
			log.Printf("get metrical app scale hosts failed [%s]", e)
			return e
		}
		fmt.Println("get metrical app scale hosts", metricalAppScaleHosts)
		// publish messages to apps
		publishMessagesToApps(metricalAppScaleHosts, appInfo.AppConf.App)
	}

	defer func() {
		// update current metrical
		(*appInfo).CurrentMetrical = appMetrical.Metrical
	}()

	return nil
}

// handler app send scale action message when app after stop itself
func (this *Decider) ScaleAction(command common.Command) error {
	message := common.InformScaleDownAppMessage{}
	err := json.Unmarshal([]byte(command.Body), &message)
	if err != nil {
	}

	this.Client.MetricalScaleAppsAction(message)
	return err
}
