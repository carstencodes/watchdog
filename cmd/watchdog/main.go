package main

import "github.com/carstencodes/watchdog/internal/lib/app"

func main() {
	watchdogApp := app.NewApp()
	err := watchdogApp.Run()
	if err != nil {
		panic(err)
	}
}
