package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func updateConfiguration(Force bool) {
	// Always update exploit mitigations
	updateExploitMitigations()

	path := "C:/Program Files/Krypton/Updates/config.zip"
	currentChannel := loadCurrentChannel()
	err := downloadToFile(currentChannel.configurationURL, path)
	if err != nil {
		log.Fatal("Error al descargar la configuracion de seguridad")
	}

	// Las actualizaciones semianuales de Windows modifican muchas
	// configuraciones y hay que volver a instalar la configuración
	// si cambia la versión de Windows
	if getWindowsVersion() != getLastUpdateWindowsVersion() {
		setLastUpdateWindowsVersion(getWindowsVersion())
		Force = true
	}

	// Si se indica --force-update hay que aplicar la configuración
	// ignorando si ya se aplicó anteriormente
	if !Force {
		configUpdateHash := getFileHash(path)
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
			runPowershellScript("./"+f.Name(),
				"C:/Program Files/Krypton/Updates/config")
		}
	}
}

func getLastUpdateHash() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Krypton", registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()

	hash, _, err := k.GetStringValue("lastUpdateHash")
	if err != nil {
		return ""
	}
	return hash
}

func setLastUpdateHash(hash string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Krypton", registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()
	k.SetStringValue("lastUpdateHash", hash)
}

func getWindowsVersion() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	buildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil {
		log.Fatal(err)
	}
	return buildNumber
}

func getWindowsPatchNumber() uint64 {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	patchNumber, _, err := k.GetIntegerValue("UBR")
	if err != nil {
		log.Fatal(err)
	}
	return patchNumber
}

func getLastUpdateWindowsVersion() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Krypton", registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()

	buildNumber, _, err := k.GetStringValue("lastBuildNumber")
	if err != nil {
		return ""
	}
	return buildNumber
}

func isWoW64() bool {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		log.Fatal("Error al cargar kernel32")
	}
	defer dll.Release()

	proc, err := dll.FindProc("IsWow64Process")
	if err != nil {
		log.Fatal(err)
	}

	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		log.Fatal(err)
	}

	var result bool
	_, _, _ = proc.Call(uintptr(handle), uintptr(unsafe.Pointer(&result)))
	return result
}

func setLastUpdateWindowsVersion(buildNumber string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		"SOFTWARE\\Krypton", registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()
	k.SetStringValue("lastBuildNumber", buildNumber)
}

func updateExploitMitigations() {
	path := "C:/Program Files/Krypton/Updates/Settings.xml"
	currentChannel := loadCurrentChannel()
	err := downloadToFile(currentChannel.exploitMitigationsURL, path)
	if err != nil {
		log.Println("Error al descargar la configuracion contra exploits")
		return
	}
	runPowershellScript("Set-ProcessMitigation -PolicyFilePath Settings.xml",
		"C:/Program Files/Krypton/Updates")
	log.Println("Actualizada configuracion contra exploits")
}

func runPowershellScript(flags string, workingDir string) {
	var powershellPath string
	if isWoW64() {
		powershellPath = "c:/windows/sysnative/WindowsPowerShell/v1.0/powershell.exe"
	} else {
		powershellPath = "powershell.exe"
	}
	cmd := exec.Command(powershellPath, "-ExecutionPolicy", "Bypass", flags)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}