package main

import "github.com/carstencodes/watchdog/internal/lib/app"

func main() {
	watchdogApp, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	err = watchdogApp.Run()
	if err != nil {
		panic(err)
	}
}
