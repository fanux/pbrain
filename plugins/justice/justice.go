package justice

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/fanux/pbrain/common"
)

const (
	COMMAND_APP_METRICAL = "app-metrical"

	PLUGIN_NAME = "plugin_justice"
)

type Justice struct {
	*common.PluginBase

	Documents map[string]Document // the key is strategy name, value is strategy document
}

type Document struct {
	App  string
	Sepc []AppMetricalInterval
}

type AppMetricalInterval struct {
	Metrical []int
	Apps     []common.AppScale
}

type AppMetrical struct {
	App      string
	Metrical int
}

func (this *Justice) initDocuments() error {
	this.Documents = make(map[string]Document)

	for _, strategy := range this.Strategies {
		document := Document{}
		err := json.Unmarshal([]byte(strategy.Document), &document)
		if err != nil {
			return errors.New("init document error")
		}
		this.Documents[strategy.Name] = document
	}

	return nil
}

func (this *Justice) Start() error {
	this.initDocuments()
	common.RegistCommand(COMMAND_APP_METRICAL, this.OnScale)

	return nil
}

func (this *Justice) Stop() error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
	return nil
}

func (this *Justice) DisableStrategy(strategyName string) error {
	common.UnRegistCommand(COMMAND_APP_METRICAL)
	return nil
}

func (this *Justice) EnableStrategy(strategyName string) error {
	this.initDocuments()
	common.RegistCommand(COMMAND_APP_METRICAL, this.OnScale)
	return nil
}

func (this *Justice) UpdateDocument(strategyName string) error {
	this.initDocuments()
	return nil
}

func (this *Justice) getDocument(appName string) (Document, error) {
	for _, document := range this.Documents {
		if document.App == appName {
			return document, nil
		}
	}

	return Document{}, errors.New("can not find app document")
}

func (this *Justice) getAppScales(appMetrical AppMetrical) ([]common.AppScale, error) {
	document, err := this.getDocument(appMetrical.App)
	if err != nil {
		log.Printf("[%s] [%s]", appMetrical.App, err)
		return []common.AppScale{}, err
	}

	for _, appMetricalInterval := range document.Sepc {
		if appMetrical.Metrical > appMetricalInterval.Metrical[0] &&
			appMetrical.Metrical <= appMetricalInterval.Metrical[1] {
			return appMetricalInterval.Apps, nil
		}
	}

	return []common.AppScale{}, errors.New("get app scales failed")
}

func (this *Justice) OnScale(command common.Command) error {
	appMetrical := AppMetrical{}
	err := json.Unmarshal([]byte(command.Body), &appMetrical)

	if err != nil {
		log.Printf("Unmarshal command body failed")
		return err
	}

	appScales, e := this.getAppScales(appMetrical)
	if e != nil {
		log.Printf("[%s]", e)
		return e
	}

	this.Client.ScaleApps(appScales)
	return nil
}
