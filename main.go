package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/slotopol/balance/core"
)

func main() {
	var a = app.NewWithID("slotopol.balance")
	core.Lifecycle(a)
	var frame = &core.Frame{}
	frame.CreateWindow(a)
	frame.ShowAndRun()
}
