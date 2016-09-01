package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Api struct {
	host string
	port string
}

func (this Api) GetPluginInfo(pluginName string) (Plugin, error) {
	//log.Printf("get plugin info:%s", pluginName)

	p := new(Plugin)

	/*
		p := Plugin{
			Name:        "plugin_pipeline",
			Kind:        "",
			Status:      "enable",
			Description: "",
			Spec:        "",
			Manual:      "",
		}
	*/

	client := &http.Client{}

	url := fmt.Sprintf("http://%s%s/plugins/%s", this.host, this.port, pluginName)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("get plugin info failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, p)

	log.Printf("get plugin info: %s", p.Description)

	return *p, err
}

func (this Api) GetPluginStrategies(pluginName string) ([]Strategy, error) {
	/*
		log.Printf("get plugin strategies:%s", pluginName)
		strategies := []Strategy{
			{
				PluginName: "plugin_pipeline",
				Name:       "scale-by-hour",
				Status:     "enable",
				Document:   `[{"Cron":"1","Apps":[{"App":"ats","Number":20},{"App":"hadoop","Number":10}]}]`,
			},
		}
	*/

	strategies := []Strategy{}

	client := &http.Client{}

	url := fmt.Sprintf("http://%s%s/plugins/%s/strategies", this.host, this.port, pluginName)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("get plugin strategies info failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &strategies)

	return strategies, err
}

func (this Api) ScaleApps(appscale []AppScale) error {
	s, _ := json.Marshal(appscale)
	log.Printf("scale apps: \n%s\n", string(s))

	client := &http.Client{}

	url := fmt.Sprintf("http://%s%s/plugins/scale", this.host, this.port)

	req, _ := http.NewRequest("POST", url, strings.NewReader(string(s)))

	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("get plugin strategies info failed: %s", err)
		fmt.Println(resp)
	}

	return err
}
