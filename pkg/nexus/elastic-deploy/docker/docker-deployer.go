package docker

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	deployer_models "github.com/nuclio/nuclio/pkg/nexus/elastic-deploy/models"
	"strings"
	"time"
)

type DockerDeployer struct {
	*deployer_models.ProElasticDeployerConfig
	*docker.Client

	baseContainerName          string
	durationFunctionsContainer *map[string]time.Time
}

func NewDockerDeployer(baseContainerName string, config *deployer_models.ProElasticDeployerConfig, durationFunctionsContainer *map[string]time.Time) *DockerDeployer {
	return &DockerDeployer{
		baseContainerName:          baseContainerName,
		ProElasticDeployerConfig:   config,
		durationFunctionsContainer: durationFunctionsContainer,
	}
}

func (ds *DockerDeployer) Initialize() {
	ds.Client, _ = docker.NewClientFromEnv()

	container, err := ds.GetNuclioFunctionContainer()
	if err != nil {
		panic(err)
	}

	fmt.Println("The Nucleo function containers are:", container)
	for _, container := range *container {
		pauseTime := time.Now().Add(ds.MaxIdleTime)
		functionName := strings.TrimPrefix(container, "/")
		functionName = strings.TrimPrefix(functionName, ds.baseContainerName)
		(*ds.durationFunctionsContainer)[functionName] = pauseTime
	}

	fmt.Println("The durationFunctionContainer is:", ds.durationFunctionsContainer)
}

func (ds *DockerDeployer) GetNuclioFunctionContainer() (*[]string, error) {
	options := &docker.ListContainersOptions{
		Filters: map[string][]string{"name": {ds.baseContainerName}},
	}

	container, err := ds.ListContainers(*options)
	if err != nil {
		return nil, err
	}

	nexusContainer := make([]string, len(container))
	for i, c := range container {
		nexusContainer[i] = c.Names[0]
	}
	return &nexusContainer, nil
}

func (ds *DockerDeployer) Unpause(functionName string) error {
	container := ds.getFunctionContainer(functionName)
	if ds.IsRunning(functionName) {
		fmt.Printf("Container %s has been running already\n", ds.getContainerName(functionName))
		return nil
	}

	fmt.Println("Container state: ", container.State)
	if container.State == deployer_models.Paused {
		err := ds.UnpauseContainer(container.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Container %s unpaused\n", ds.getContainerName(functionName))
		(*ds.durationFunctionsContainer)[functionName] = time.Now().Add(ds.MaxIdleTime)

		return nil
	}

	fmt.Println("Container state: ", container.State, "does not match any of the expected states")
	return nil
}

func (ds *DockerDeployer) Pause(functionName string) error {
	container := ds.getFunctionContainer(functionName)
	if container.State == deployer_models.Paused {
		fmt.Printf("Container %s has been paused already\n", ds.getContainerName(functionName))
		return nil
	}

	if container.State == deployer_models.Running {
		err := ds.PauseContainer(container.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Container %s paused\n", ds.getContainerName(functionName))
		return nil
	}

	fmt.Println("Container state: ", container.State, "does not match any of the expected states")
	return nil
}

func (ds *DockerDeployer) IsRunning(functionName string) bool {
	container := ds.getFunctionContainer(functionName)
	return container.State == deployer_models.Running
}

func (ds *DockerDeployer) getFunctionContainer(functionName string) *docker.APIContainers {
	options := &docker.ListContainersOptions{
		Filters: map[string][]string{"name": {ds.getContainerName(functionName)}}, // Names
		Limit:   1,
	}

	container, err := ds.ListContainers(*options)
	if err != nil {
		fmt.Println("Error getting container: ", err)
		return nil
	}

	if len(container) == 0 {
		fmt.Println("No container existing: ", err)
		return nil
	}
	return &container[0]
}

func (ds *DockerDeployer) getContainerName(functionName string) string {
	return ds.baseContainerName + functionName
}
