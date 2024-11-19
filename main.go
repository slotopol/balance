package main

import (
	"fyne.io/fyne/v2/app"

	cfg "github.com/slotopol/balance/config"
	"github.com/slotopol/balance/ui"
)

func main() {
	var a = app.NewWithID(cfg.AppID)
	ui.Lifecycle(a)
	var frame = &ui.Frame{}
	frame.CreateWindow(a)
	frame.ShowAndRun()
}
