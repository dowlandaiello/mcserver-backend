package commands

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/xoreo/mcserver-backend/types"
)

func getTestServer() (*types.Server, error) {
	rand.Seed(time.Now().UnixNano())
	min := 25565
	max := 26000
	random := rand.Intn(max-min) + min

	server, err := types.NewServer("1.7.2", "server-"+strconv.Itoa(random), uint32(random), 1024)
	if err != nil {
		return nil, err
	}

	err = InitializeServer(server)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func TestInitializeServer(t *testing.T) {
	_, err := getTestServer()
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecute(t *testing.T) {
	server, err := getTestServer()
	if err != nil {
		t.Fatal(err)
	}

	output, err := Execute("start", *server)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)

}
