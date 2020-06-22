package tester

import (
<<<<<<< HEAD
	"fmt"
	"log"
	"os"

	"github.com/Blockdaemon/bpm-sdk/pkg/node"
)

// CeloTester Interface for running tests against node
type CeloTester struct{}

// Test Method for calling tests against node
func (d CeloTester) Test(currentNode node.Node) (bool, error) {
	if err := runAllTests(); err != nil {
		return false, err
	}
=======
	"bytes"
	"context"
    "strconv"
	"fmt"
	"io/ioutil"
    "strings"
    "time"
    "net/http"
    "math/rand"
    "encoding/json"

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

>>>>>>> release-0.0.2
	return true, nil
}

type testRunner struct {
	failed    int
	succeeded int
<<<<<<< HEAD
}

func (t *testRunner) test(testFunc func() error) {
	if err := testFunc(); err != nil {
		t.failed++
	} else {
		t.succeeded++
	}
}

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
=======
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

    // test rpc, if no error then RPC is working.
    testCase = func()(string, string, error){

        title := "JSON RPC"
        rpcEndpoint, err := getContainerEndpoint("/"+containerName)
        if err != nil {
            return title, "false", err
        }
        fmt.Printf("RPC call to %s at %s\n", containerName, rpcEndpoint)

        _, _, _, err = rpcPost("eth_syncing", "", rpcEndpoint)
        if err!= nil { return title, "false", err }

        return title, "true", nil
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

func rpcPost(method string, params string, baseUrl string) (int, int, map[string]interface{}, error) {
	rand.Seed(time.Now().UnixNano())
	var messageID = rand.Int()
    var data map[string]interface{}

	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  method,
		"id":      messageID,
		"jsonrpc": "2.0",
	})
	if err != nil {
		return int(500), messageID, data, err
	}

	resp, err := http.Post(baseUrl+params, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
        return int(500), messageID, data, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
        return int(500), messageID, data, err
	}

	return resp.StatusCode, messageID, data, err
}

func getContainerEndpoint(name string) (string, error){

    cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

    containerJSON, _ := cli.ContainerInspect(context.Background(), name)
    if err != nil {
		return "", err
	}
    hostPort := containerJSON.NetworkSettings.NetworkSettingsBase.Ports["8545/tcp"][0].HostPort
    hostIP := containerJSON.NetworkSettings.NetworkSettingsBase.Ports["8545/tcp"][0].HostIP

    return "http://"+hostIP+":"+hostPort, nil
>>>>>>> release-0.0.2
}
