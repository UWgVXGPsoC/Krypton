package installer

import (
	"io"
	"os"
	"os/exec"
	"time"
)

func Install() error {
	err := os.Mkdir("C:/Program Files/Krypton", os.ModeDir)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	os.Mkdir("C:/Program Files/Krypton/Updates", os.ModeDir) //nolint

	// workaround
	time.Sleep(3 * time.Second) //nolint

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	err = copyFile(exe, "C:/Program Files/Krypton/Krypton.exe")
	if err != nil {
		return err
	}

	err = createScheduledTasks()
	if err != nil {
		return err
	}
	return nil
}

// schtasks.exe /query /tn KryptonUpdate /xml
func createScheduledTasks() error {
	path := "C:/Program Files/Krypton/Krypton.exe"
	cmd := exec.Command("schtasks.exe", "/Create", "/SC", "HOURLY", "/TN",
		"KryptonUpdate", "/RU", "SYSTEM", "/F", "/TR", path+" --update")

	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("schtasks.exe", "/Create", "/SC", "DAILY", "/TN",
		"KryptonUpgrade", "/RU", "SYSTEM", "/F", "/TR", path+" --upgrade")

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
