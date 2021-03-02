package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/milgradesec/krypton/internal/installer"
	"github.com/milgradesec/krypton/internal/system"
	"github.com/milgradesec/krypton/internal/updater"
)

func main() {
	var (
		versionFlag     = flag.Bool("version", false, "Show version information.")
		installFlag     = flag.Bool("install", false, "Install Krypton.")
		updateFlag      = flag.Bool("update", false, "Apply latest system configuration from remote server.")
		forceUpdateFlag = flag.Bool("force-update", false, "Force system configuration update.")
		upgradeFlag     = flag.Bool("upgrade", false, "Updates Krypton to latest version available.")
		helpFlag        = flag.Bool("help", false, "Show all commands.")
	)
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println("extra command line arguments.")
		os.Exit(1)
	}

	if *versionFlag {
		fmt.Println("Krypton " + Version)
		fmt.Printf("%s/%s, %s, %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version(), Version)
	}

	if *installFlag {
		err := installer.Install()
		if err != nil {
			fmt.Printf("Error: install failed: %v\n", err)
		}
	}

	if *updateFlag {
		err := system.UpdateConfig(false)
		if err != nil {
			fmt.Printf("Error: system settings update failed: %v\n", err)
		}
	}

	if *forceUpdateFlag {
		err := system.UpdateConfig(true)
		if err != nil {
			fmt.Printf("Error: system settings update failed: %v\n", err)
		}
	}

	if *upgradeFlag {
		err := updater.Update(Version)
		if err != nil {
			fmt.Printf("Error: update failed: %v\n", err)
		}
	}

	if *helpFlag {
		flag.PrintDefaults()
	}
}

var (
	Version string
)
