package main

import (
	"time"

	tm "github.com/buger/goterm"

	app "github.com/chocnut/sentry-api/services"
)

func main() {
	tm.Clear()
	for {
		tm.MoveCursor(1, 1)
		tm.Flush()
		app.Run()

		time.Sleep(time.Second * 60)
	}

}
