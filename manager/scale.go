package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/fanux/pbrain/common"
	"github.com/samalba/dockerclient"
)

type ContainerNumberInfo struct {
	Current     int
	Need        int
	ContainerId string
}

type ScaleResult struct {
	Scaled []string
	Errors []string
}

func initScaleInfo(info map[string]ContainerNumberInfo, scaleApp []ScaleApp) {
	for _, v := range scaleApp {
		info[v.App] = ContainerNumberInfo{Current: 0, Need: v.Number, ContainerId: ""}
	}
}

func releaseContainers(info map[string]ContainerNumberInfo, client *dockerclient.DockerClient) {
	containers, err := client.ListContainers(true, false, "")

	if err != nil {
	}

	for _, c := range containers {
		fmt.Printf("container image name:%s\n", c.Image)
		containerNumberInfo, ok := info[c.Image]
		if !ok {
			fmt.Printf("out of scale:%s\n", c.Image)
			continue
		}
		containerNumberInfo.Current++

		if containerNumberInfo.ContainerId == "" {
			containerNumberInfo.ContainerId = c.Id
		}

		info[c.Image] = containerNumberInfo

		cid, ok2 := info[c.Image]
		if ok2 {
			fmt.Printf("image [%s] container id is:%s\n", c.Image, cid.ContainerId)
		}

		if containerNumberInfo.Current > containerNumberInfo.Need {
			// stop container with 5 seconds timeout
			client.StopContainer(c.Id, 5)
			// force remove, delete volume
			client.RemoveContainer(c.Id, true, true)
		}
	}
}

func deployContainers(info map[string]ContainerNumberInfo, client *dockerclient.DockerClient) {
	for _, v := range info {
		if v.Current < v.Need {
			n := v.Need - v.Current
			scaleContainer(v.ContainerId, n, client)
		}
	}
}

func scaleContainerByImageName(imageName string, numInstances int,
	client *dockerclient.DockerClient) ScaleResult {

	containers, err := client.ListContainers(true, false, "")
	if err != nil {
	}

	for _, c := range containers {
		if c.Image == imageName {
			return scaleContainer(c.Id, numInstances, client)
		}
	}
	return ScaleResult{Scaled: make([]string, 0), Errors: make([]string, 0)}
}

func scaleContainer(id string, numInstances int, client *dockerclient.DockerClient) ScaleResult {
	var (
		errChan = make(chan (error))
		resChan = make(chan (string))
		result  = ScaleResult{Scaled: make([]string, 0), Errors: make([]string, 0)}
	)

	// docker client get container info
	containerInfo, err := client.InspectContainer(id)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result
	}

	for i := 0; i < numInstances; i++ {
		go func(instance int) {
			config := containerInfo.Config
			// clear hostname to get a newly generated
			config.Hostname = ""
			hostConfig := containerInfo.HostConfig
			config.HostConfig = *hostConfig // sending hostconfig via the Start-endpoint is deprecated starting with docker-engine 1.12
			// using docker client create Container
			id, err := client.CreateContainer(config, "", nil)
			if err != nil {
				errChan <- err
				return
			}
			// using docker client start container
			if err := client.StartContainer(id, hostConfig); err != nil {
				errChan <- err
				return
			}
			resChan <- id
		}(i)
	}

	for i := 0; i < numInstances; i++ {
		select {
		case id := <-resChan:
			result.Scaled = append(result.Scaled, id)
		case err := <-errChan:
			result.Errors = append(result.Errors, strings.TrimSpace(err.Error()))
		}
	}

	return result
}

func (this PluginResource) scaleApp(request *restful.Request,
	response *restful.Response) {

	scaleApp := []ScaleApp{}

	dockerClient, err := dockerclient.NewDockerClient(DockerHost, nil)
	if err != nil {
		fmt.Printf("init docker client error:%s", err)
	}

	err = request.ReadEntity(&scaleApp)
	if err != nil {
	}
	fmt.Println("scale : ", scaleApp)

	/*
		{
			"ats:latest":{2, 20}
			"hadoop:v1.0":{20, 2}
		}
	*/

	scaleInfo := make(map[string]ContainerNumberInfo)

	initScaleInfo(scaleInfo, scaleApp)

	fmt.Println("map info: ", scaleInfo)

	releaseContainers(scaleInfo, dockerClient)
	deployContainers(scaleInfo, dockerClient)

	//fmt.Println("scale map: ", scaleInfo["ats"])
}

func showMap(hostsMap map[string]common.MetricalAppScaleHosts) {
	fmt.Println("+++++++++++++++++++++++++++++++++")
	for k, v := range hostsMap {
		fmt.Println("map key: ", k, " map value: ", v)
	}
	fmt.Println("+++++++++++++++++++++++++++++++++")
}

// return the one who want scale up
func initMetricalHostsMap(hostsMap map[string]common.MetricalAppScaleHosts,
	scaleApp []common.MetricalAppScale) common.MetricalAppScale {

	scaleUpApp := common.MetricalAppScale{}
	for _, v := range scaleApp {
		temp := common.MetricalAppScale{v.App, v.Number, v.MinNum}
		hostsMap[v.App] = common.MetricalAppScaleHosts{temp, []string{}}
		if v.Number > 0 {
			scaleUpApp = v
		}
	}

	fmt.Println("init hosts map: ", hostsMap)

	return scaleUpApp
}

func setAppNumber(appName string, hostsMap map[string]common.MetricalAppScaleHosts, length int) (int, error) {
	// set the scale up app number
	metricalAppScaleHosts, ok := hostsMap[appName]
	if !ok {
		log.Printf("get scale up app failed: [%s]", appName)
		return 0, errors.New("can't get scale up app")
	}
	// len(res.Scaled) apps is successed deployed
	metricalAppScaleHosts.Number -= length
	hostsMap[appName] = metricalAppScaleHosts

	return metricalAppScaleHosts.Number, nil
}

type ContainerNode struct {
	Id   string
	Node Node
}

func getContainerHostIp(Id string) string {
	// docker client don't support the node field!!!
	// can't using docker client, call swarm api directly!!!
	// cant't using container, err := client.InspectContainer(Id)
	// using DockerHost

	cNode := ContainerNode{}

	client := &http.Client{}

	url := fmt.Sprintf("%s/containers/%s/json", DockerHost, Id)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("get plugin info failed: %s", err)
		return ""
	}
	defer resp.Body.Close()

	/*
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &cNode)
	*/
	json.NewDecoder(resp.Body).Decode(&cNode)
	fmt.Println("get container node: ", cNode)

	//return cNode.Node.Addr
	return cNode.Node.IP
}

func updateScaleDownAppHosts(hostsMap map[string]common.MetricalAppScaleHosts,
	client *dockerclient.DockerClient) {

	containers, err := client.ListContainers(true, false, "")
	if err != nil {
	}
	for _, c := range containers {
		temp, ok := hostsMap[c.Image]
		if !ok || len(temp.Hosts) >= -temp.Number {
			continue
		}
		if temp.Number < 0 {
			ip := getContainerHostIp(c.Id)
			if ip != "" {
				temp.Hosts = append(temp.Hosts)
			}
		}
		hostsMap[c.Image] = temp
	}
}

func initHosts(hosts []common.MetricalAppScaleHosts,
	hostsMap map[string]common.MetricalAppScaleHosts) {

	for _, v := range hostsMap {
		hosts = append(hosts, v)
	}
}

func (this PluginResource) metricalScaleApp(request *restful.Request,
	response *restful.Response) {

	scaleApp := []common.MetricalAppScale{}

	hosts := []common.MetricalAppScaleHosts{}
	hostsMap := make(map[string]common.MetricalAppScaleHosts)

	dockerClient, err := dockerclient.NewDockerClient(DockerHost, nil)
	if err != nil {
		fmt.Printf("init docker client error:%s", err)
	}

	err = request.ReadEntity(&scaleApp)
	if err != nil {
		log.Printf("get metrical scale info failed")
	}
	fmt.Println("scale : ", scaleApp)

	scaleUpApp := initMetricalHostsMap(hostsMap, scaleApp)
	//showMap(hostsMap)

	// first round scale
	res := scaleContainerByImageName(scaleUpApp.App, scaleUpApp.Number, dockerClient)
	n, e := setAppNumber(scaleUpApp.App, hostsMap, len(res.Scaled))
	if e != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, hosts)
		return
	}
	if n <= 0 {
		//send empty list to plugin
		response.WriteHeaderAndEntity(http.StatusOK, hosts)
		return
	}

	// get container list and update host list
	updateScaleDownAppHosts(hostsMap, dockerClient)
	initHosts(hosts, hostsMap)
	fmt.Println(hosts)
	response.WriteHeaderAndEntity(http.StatusOK, hosts)
}
