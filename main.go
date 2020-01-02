package main

import (
	"time"

	tm "github.com/buger/goterm"

	app "github.com/chocnut/sentry-api/services"
)

func runApp() {
	tm.MoveCursor(1, 1)
	tm.Flush()
	app.Run()
}

func main() {
	tm.Clear()
	for {
		go runApp()
		<-time.After(60 * time.Second)
	}

}
