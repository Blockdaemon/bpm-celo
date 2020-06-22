package tester_test

import (
    "testing"
    "context"
    "fmt"

    "github.com/davecgh/go-spew/spew"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
)

func TestGetContainer(t *testing.T){
    cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

    name := "/bpm-celo-test-fullnode"

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	if len(containers) > 0 {
		for _, container := range containers {
            spew.Dump(container.NetworkSettings.Networks["bpm"].IPAddress, "^ container.NetworkSettings.Networks[\"bpm\"].IPAddress")
            if container.Names[0] == name {
                // ip := container.NetworkSettings.Networks[0].IpAddress
                fmt.Printf("Container IP: %s\n", container.ID)
            }
		}
	} else {
		fmt.Println("There are no containers running")
	}
}