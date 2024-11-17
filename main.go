package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/slotopol/balance/ui"
)

func main() {
	var a = app.NewWithID("slotopol.balance")
	ui.Lifecycle(a)
	var frame = &ui.Frame{}
	frame.CreateWindow(a)
	frame.ShowAndRun()
}
