package main

import (
	"runtime"

	"gofr.dev/pkg/gofr"

	"ssshekhu53/folder-lock/handlers"
	"ssshekhu53/folder-lock/services"
	cryptPkg "ssshekhu53/folder-lock/services/crypt"
	"ssshekhu53/folder-lock/services/mac"
)

func main() {
	var service services.FolderLock

	app := gofr.NewCMD()

	crypt, err := cryptPkg.New()
	if err != nil {
		app.Logger().Fatalf("Error occurred: %v", err)
	}

	switch runtime.GOOS {
	case "darwin":
		service = mac.New(crypt)
	}

	handler := handlers.New(service)

	app.SubCommand("init", handler.Init)
	app.SubCommand("unlock", handler.Unlock)
	app.SubCommand("lock", handler.Lock)
	app.SubCommand("help", handler.Help)

	app.Run()
}
