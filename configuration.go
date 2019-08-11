package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

type powerShellScript struct {
	Args       string
	WorkingDir string
}

func (ps powerShellScript) Run() error {
	return nil
}

func updateConfig(force bool) {
	err := updateExploitMitigations()
	if err != nil {
		fmt.Printf("Error al actualizar la configuración contra exploits, %v\n", err)
	} else {
		fmt.Println("Actualizada configuracion contra exploits.")
	}

	url := "https://dl.paesacybersecurity.eu/krypton/config/stable/config.zip"
	path := "C:/Program Files/Krypton/Updates/config.zip"
	err = downloadToFile(url, path)
	if err != nil {
		log.Fatal("Error al descargar la configuracion de seguridad")
	}

	// Las actualizaciones semianuales de Windows modifican muchas
	// configuraciones y hay que volver a instalar la configuración
	// si cambia la versión de Windows
	if getWindowsVersion() != getLastUpdateWindowsVersion() {
		setLastUpdateWindowsVersion(getWindowsVersion())
		force = true
	}

	// Si se indica --force-update hay que aplicar la configuración
	// ignorando si ya se aplicó anteriormente
	if !force {
		configUpdateHash := computeFileSHA1(path)
		if configUpdateHash == getLastUpdateHash() {
			log.Println("No hay cambios de configuracion")
			os.Exit(0)
		}
		log.Println("Hay nueva configuracion disponible")
		setLastUpdateHash(configUpdateHash)
	}

	// Descomprimir configuración
	os.RemoveAll("C:\\Program Files\\Krypton\\Updates\\config")
	err = unzip(path, "C:\\Program Files\\Krypton\\Updates")
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir("C:/Program Files/Krypton/Updates/config")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".ps1") {
			err = runPowershellScript("./"+f.Name(),
				"C:/Program Files/Krypton/Updates/config")
			if err != nil {
				log.Println(err)
			}
		}
	}

	dir, err := os.Stat("C:/Program Files/Krypton/Settings")
	if err != nil {
		log.Fatal(err)
	}

	if dir.IsDir() {
		files, err := ioutil.ReadDir("C:/Program Files/Krypton/Settings")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".ps1") {
				err = runPowershellScript("./"+f.Name(),
					"C:/Program Files/Krypton/Settings")
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func isWoW64() (bool, error) {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return false, err
	}
	defer dll.Release()

	proc, err := dll.FindProc("IsWow64Process")
	if err != nil {
		return false, err
	}

	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return false, err
	}

	var result bool
	_, _, _ = proc.Call(uintptr(handle), uintptr(unsafe.Pointer(&result)))
	return result, nil
}

func updateExploitMitigations() error {
	err := downloadToFile("https://dl.paesacybersecurity.eu/krypton/Settings.xml",
		"C:/Program Files/Krypton/Updates/Settings.xml")
	if err != nil {
		return err
	}

	err = runPowershellScript("Set-ProcessMitigation -PolicyFilePath Settings.xml",
		"C:/Program Files/Krypton/Updates")
	if err != nil {
		return err
	}

	err = runPowershellScript("Set-ProcessMitigation -PolicyFilePath Settings.xml",
		"C:/Program Files/Krypton/Settings")
	if err != nil {
		return err
	}
	return nil
}

func runPowershellScript(flags string, workingDir string) error {
	var powershellPath string
	wow64, err := isWoW64()
	if err != nil {
		return err
	}

	if wow64 {
		powershellPath = "c:/windows/sysnative/WindowsPowerShell/v1.0/powershell.exe"
	} else {
		powershellPath = "powershell.exe"
	}
	cmd := exec.Command(powershellPath, "-ExecutionPolicy", "Bypass", flags)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}