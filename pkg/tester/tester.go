package tester

import "github.com/Blockdaemon/bpm-sdk/pkg/node"

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

func runAllTests() error {
	return nil
}
