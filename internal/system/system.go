package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const baseURL = "https://dl.paesa.es/krypton/"

func updateExploitMitigations() error {
	err := downloadToFile(baseURL+"config/stable/Settings.xml",
		"C:/Program Files/Krypton/Updates/Settings.xml")
	if err != nil {
		return err
	}

	err = powerShellRun("Set-ProcessMitigation -PolicyFilePath Settings.xml",
		"C:/Program Files/Krypton/Updates")
	if err != nil {
		return err
	}

	_, err = os.Stat("C:/Program Files/Krypton/Settings/Settings.xml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = powerShellRun("Set-ProcessMitigation -PolicyFilePath Settings.xml",
		"C:/Program Files/Krypton/Settings")
	if err != nil {
		return err
	}

	return nil
}

func UpdateConfig(force bool) error {
	updateExploitMitigations()

	url := baseURL + "configurarWindows10.ps1"
	path := "C:/Program Files/Krypton/Updates/configurarWindows10.ps1"
	err := downloadToFile(url, path)
	if err != nil {
		return err
	}

	err = powerShellRun("./configurarWindows10.ps1",
		"C:/Program Files/Krypton/Updates")
	if err != nil {
		fmt.Println(err)
	}

	dir, err := os.Stat("C:/Program Files/Krypton/Settings")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if dir.IsDir() {
		files, err := ioutil.ReadDir("C:/Program Files/Krypton/Settings")
		if err != nil {
			return err
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".ps1") {
				err = powerShellRun("./"+f.Name(),
					"C:/Program Files/Krypton/Settings")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
}
