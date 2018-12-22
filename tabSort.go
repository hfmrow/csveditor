// tabSort.go

package main

import (
	"github.com/andlabs/ui"
)

// Prepare to display Tab
func tabSort() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	vbox.Disable() // Used to know if control has been already initialised

	return tabMakeSort(vbox)
}
