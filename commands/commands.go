package commands

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/xoreo/mcserver-backend/common"
	"github.com/xoreo/mcserver-backend/types"
)

var (
	// ErrUnsupportedVersion is thrown when an unsupported server version is given.
	ErrUnsupportedVersion = errors.New("that is not a supported version")

	// ErrServerHasNotBeenInitialized is thrown when a server's metadata exists but the server has not actually been initialized on the local machine.
	ErrServerHasNotBeenInitialized = errors.New("that server has not actually been initialized yet. Initialize it with InitializeServer()")
)

func newStartScript(path string, ram int) []byte {
	header := "#!/bin/bash\njava -Xms" + string(ram) + "MB -Xmx" + string(ram) + "MB -jar "
	body := path
	footer := " nogui\n"
	return []byte(header + body + footer)
}

func downloadServerJar(url, localPath, version string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the directory
	err = common.CreateDirIfDoesNotExist(localPath)
	if err != nil {
		return "", err
	}

	// Create the file
	zipPath := filepath.Join(localPath, version+".zip")
	out, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return zipPath, nil
}

// InitializeServer initializes a new server onto the local machine.
func InitializeServer(server *types.Server) error {
	var url string
	dServer := *server // Make a copy of the pointer (dereference for convenience)

	// Determine the pre-made server download url
	switch dServer.Version {
	case "1.12":
		url = common.ServerV112
	case "1.8":
		url = common.ServerV18
	case "1.7.2":
		url = common.ServerV172
	case "1.2.1":
		url = common.ServerV121
	default:
		return ErrUnsupportedVersion
	}

	// Download the pre-made server
	zipPath, err := downloadServerJar(url, dServer.Path, dServer.Version)
	if err != nil {
		return err
	}

	// Unzip the downloaded file
	_, err = common.Unzip(zipPath, dServer.Path)
	if err != nil {
		return err
	}

	// Create start script for the server
	startScriptPath := filepath.Join(dServer.Path, dServer.Version, "start.sh")
	serverJarPath := filepath.Join(dServer.Path, dServer.Version, dServer.Version+".jar")
	script := newStartScript(serverJarPath, dServer.RAM)

	// Install the script
	err = ioutil.WriteFile(startScriptPath, script, 0644)
	if err != nil {
		return err
	}

	server.Initialized = true // Set the server's initialized state to true
	return nil
}

// StartServer starts a server.
func StartServer(server *types.Server) error {
	// Make sure that the server has been initialized.
	if !(*server).Initialized {
		return ErrServerHasNotBeenInitialized
	}
	dServer := *server // Dereference for convenience

	launchPrefix := "cd " + dServer.Path + " && "
	launcher := launchPrefix + filepath.Join(dServer.Path, dServer.Version, "start.sh")
	cmd := exec.Command("/bin/sh", launcher)
	fmt.Println(cmd.Output())

	return nil
}

// RestartServer restarts a server.
func RestartServer(server *types.Server) {

}

// StopServer stops a server.
func StopServer(server *types.Server) {

}

// EnterServer enters the shell of the server.
func EnterServer(server *types.Server) {

}
