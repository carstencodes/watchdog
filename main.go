package main

func main() {
	app := NewApp()
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
