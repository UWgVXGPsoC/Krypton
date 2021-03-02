package system

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func downloadToFile(url, file string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("status != 200")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dest, err := os.Create(file)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = dest.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func powerShellRun(command, workingDir string) error {
	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", command)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
