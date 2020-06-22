package tester_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TestGetContainer(t *testing.T) {
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
			if container.Names[0] == name {
				// ip := container.NetworkSettings.Networks[0].IpAddress
				fmt.Printf("Container IP: %s\n", container.ID)
			}
		}
	} else {
		fmt.Println("There are no containers running")
	}
}
