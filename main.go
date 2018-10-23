package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
)

const (
	projectLabel = "com.docker.compose.project"
	serviceLabel = "com.docker.compose.service"
)

type ImgCandidate struct {
	serviceName, imageName, imageTag string
}

func main() {
	containers := listContainers()
	imgCandidates := make([]ImgCandidate, 0)
	for i := range containers {
		imgCandidates = append(imgCandidates, getCandidate(containers[i]))
	}
	spew.Dump(imgCandidates)

	yml, err := ioutil.ReadFile("/app/base/test.yml")
	if err != nil {
		panic(err)
	}

	var composeFile ComposeFile
	err = yaml.Unmarshal(yml, &composeFile)
	if err != nil {
		panic(err)
	}
	for service := range composeFile.Services {
	}
	spew.Dump(composeFile)
}

func listContainers() []types.Container {
	cid, exists := os.LookupEnv("HOSTNAME")
	if !exists {
		panic("No container id! Seems we are not running inside docker!")
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	insp, err := cli.ContainerInspect(context.Background(), cid)
	if err != nil {
		panic(err)
	}
	project, ok := insp.Config.Labels[projectLabel]
	if !ok {
		panic("Could not find required docker-compose meta information!")
	}

	labelFilter, err := filters.ParseFlag(fmt.Sprintf("label=%s=%s", projectLabel, project), filters.NewArgs())
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(),
		types.ContainerListOptions{Filters: labelFilter})
	if err != nil {
		panic(err)
	}
	return containers
}

func getCandidate(container types.Container) ImgCandidate {
	imgTags := strings.Split(container.Image, ":")
	if len(imgTags) < 2 {
		imgTags = append(imgTags, "latest")
	}
	candidate := ImgCandidate{
		serviceName: container.Labels[serviceLabel],
		imageName:   imgTags[0],
		imageTag:    imgTags[1],
	}
	return candidate
}
