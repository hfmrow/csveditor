// tabOptions.go

package main

import (
	"github.com/andlabs/ui"
)

// Prepare to display informations
func tabOptions() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	vbox.Disable() // Used to know if control has been already initialised
	optionsBox = tabMakeOptions(vbox)

	return optionsBox
}
