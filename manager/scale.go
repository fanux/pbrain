package manager

import (
	"fmt"
	"strings"

	"github.com/emicklei/go-restful"
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
