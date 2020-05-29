package tester

import (
	"bytes"
	"context"
    "strconv"
	"fmt"
	"io/ioutil"
    "strings"
	// "log"
	// "os"

    // "github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"go.blockdaemon.com/bpm/sdk/pkg/docker"
	"go.blockdaemon.com/bpm/sdk/pkg/node"
)

// CeloTester Interface for running tests against node
type CeloTester struct {
	Cli *client.Client
    n node.Node
}

func New() *CeloTester {
	ct := &CeloTester{}

	return ct
}

// Test Method for calling tests against node
func (d CeloTester) Test(currentNode node.Node) (bool, error) {

	results, err := runAllTests(currentNode)
    if err != nil {
		return false, err
	}

    for i:=0; i<len(results.Tests); i++{
        fmt.Printf("    Test [%s]   => %s\n", results.Tests[i].name, string(results.Tests[i].result))
    }

    fmt.Printf("Total failed tests: %s\n", strconv.Itoa(results.failed))
    fmt.Printf("Total passed tests: %s\n", strconv.Itoa(results.succeeded))

	return true, nil
}

type testRunner struct {
	failed    int
	succeeded int
    bm *docker.BasicManager
    Tests []testRunnerTest
}
type testRunnerTest struct {
    name string
    result string
}

func (t *testRunner) test(testFunc func() (title string, result string, err error)) error{

    var err error
    testRes := testRunnerTest{}

    title, result, err := testFunc()
    if err!=nil {
        t.failed++
    } else {
        t.succeeded++
    }
    testRes.name = title
    testRes.result = result
    t.Tests = append(t.Tests, testRes)

    return err
}

func runAllTests(currentNode node.Node) (testRunner, error) {

    tr := testRunner{}

	containerName := "bpm-" + currentNode.ID + "-" + currentNode.StrParameters["subtype"]
	fmt.Printf("testing container: %s\n", containerName)

	bm, err := docker.NewBasicManager(currentNode)
	if err != nil {
		return tr, err
	}

    var testCase func() (string, string, error)

	// test is running
    testCase = func()(string, string, error){
        title := "Container is running"
        res, err := testIsRunning(bm, containerName)
        if err!= nil { return title, "false", err }
        return title, string(res), nil

    }
    if err := tr.test(testCase); err!=nil {
        return tr, err
    }

	// test peer count
    testCase = func()(string, string, error){
        title := "Peer Count"
        res, err := testPeerCount(bm, containerName)
        if err!= nil { return title, "false", err }
        return title, string(res), nil

    }
    if err := tr.test(testCase); err!=nil {
        return tr, err
    }

	return tr, nil
}

func testIsRunning(bm *docker.BasicManager, containerName string) (string, error) {

	ctx := context.Background()
	running, err := bm.IsContainerRunning(ctx, containerName)
	if err != nil {
        return "", err
	}

	return strconv.FormatBool(running), nil
}

func testPeerCount(bm *docker.BasicManager, containerName string) (string, error) {

	ctx := context.Background()
    cmdPeerCount := []string{
        "geth",
		"--exec",
        "net.peerCount",
        "attach",
    }
	id, _ := Exec(ctx, containerName, cmdPeerCount)
	res, _ := InspectExecResp(ctx, id.ID)

	return strings.Trim(string(res.StdOut), "\n"), nil
}

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func Exec(ctx context.Context, containerID string, command []string) (types.IDResponse, error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return types.IDResponse{}, err
	}
	defer docker.Close()

	config := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
	}

	return docker.ContainerExecCreate(ctx, containerID, config)
}

func InspectExecResp(ctx context.Context, id string) (ExecResult, error) {
	var execResult ExecResult
	docker, err := client.NewEnvClient()
	if err != nil {
		return execResult, err
	}
	defer docker.Close()

	resp, err := docker.ContainerExecAttach(ctx, id, types.ExecConfig{})
	if err != nil {
		return execResult, err
	}
	defer resp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return execResult, err
		}
		break

	case <-ctx.Done():
		return execResult, ctx.Err()
	}

	stdout, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		return execResult, err
	}
	stderr, err := ioutil.ReadAll(&errBuf)
	if err != nil {
		return execResult, err
	}

	res, err := docker.ContainerExecInspect(ctx, id)
	if err != nil {
		return execResult, err
	}

	execResult.ExitCode = res.ExitCode
	execResult.StdOut = string(stdout)
	execResult.StdErr = string(stderr)

    return execResult, nil
}
