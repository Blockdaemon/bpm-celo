package tester

import (
	"fmt"
	"log"
	"os"

	"go.blockdaemon.com/bpm/sdk/pkg/node"
)

// CeloTester Interface for running tests against node
type CeloTester struct{}

// Test Method for calling tests against node
func (d CeloTester) Test(currentNode node.Node) (bool, error) {
	if err := runAllTests(); err != nil {
		return false, err
	}
	return true, nil
}

type testRunner struct {
	failed    int
	succeeded int
}

func (t *testRunner) test(testFunc func() error) {
	if err := testFunc(); err != nil {
		t.failed++
	} else {
		t.succeeded++
	}
}

/**
* Two options:
 - cmd.Exec(bpm nodes ... node.json)
 - decouple logic from main and call that (more professinoal)
* @type {[type]}
*/
func runAllTests() error {

	jsonfile := os.Args[2]

	n, err := node.Load(jsonfile)
	if err != nil {
		log.Fatalf("Unable to load node json: %s\n", err)
	}

	containerName := "bpm-" + n.ID + "-" + n.StrParameters["subtype"]
	fmt.Printf("testing container: %s\n", containerName)

	// 1. use sdk docker/BasicManager to get docker client
	// 2. docker exec container to get `eth.syncing` result.

	return nil
}
