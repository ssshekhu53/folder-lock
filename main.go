package main

import (
	"runtime"
	"ssshekhu53/folder-lock/constants"
	"ssshekhu53/folder-lock/services/windows"

	"gofr.dev/pkg/gofr"

	"ssshekhu53/folder-lock/handlers"
	"ssshekhu53/folder-lock/services"
	cryptPkg "ssshekhu53/folder-lock/services/crypt"
	"ssshekhu53/folder-lock/services/unix"
)

func main() {
	var service services.FolderLock

	app := gofr.NewCMD()

	crypt, err := cryptPkg.New()
	if err != nil {
		app.Logger().Fatalf("Error occurred: %v", err)
	}

	switch runtime.GOOS {
	case constants.Darwin, constants.Linux:
		service = unix.New(crypt)
	case constants.Windows:
		service = windows.New(crypt)
	}

	handler := handlers.New(service)

	app.SubCommand("init", handler.Init)
	app.SubCommand("unlock", handler.Unlock)
	app.SubCommand("lock", handler.Lock)
	app.SubCommand("help", handler.Help)

	app.Run()
}
