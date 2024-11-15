package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/slotopol/balance/core"
)

func main() {
	var a = app.NewWithID("slotopol.balance")
	core.Lifecycle(a)
	var w = core.CreateMainWindow(a)
	w.ShowAndRun()
}
