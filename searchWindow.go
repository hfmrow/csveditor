// searchWindow.go

/// +build OMIT
package main

import "github.com/andlabs/ui"

// Make dialog box with callback functions to allow or deny action. fyes give entry.text value
func DialogBoxEntrySearch(mainwin *ui.Window, title, yes, no string, fyes func(entry string, cs, ww, rx bool), fno func(), text ...string) {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	mainwin.Disable()

	dialogWindow := ui.NewWindow(title, 10, 10, false)
	dialogWindow.OnClosing(func(*ui.Window) bool {
		mainwin.Enable() // Make MAINWIN enaled
		return true
	})
	ui.OnShouldQuit(func() bool {
		dialogWindow.Destroy()
		return true
	})
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	// Display given message
	for _, txt := range text {
		vbox.Append(ui.NewLabel(txt), false) // Label with txt
	}

	// Add Option search
	hbox := ui.NewHorizontalBox()
	vbox.Append(hbox, true)
	chkWholeWord := ui.NewCheckbox("Whole word")
	chkCaseSensitive := ui.NewCheckbox("Case sensitive")
	chkRegex := ui.NewCheckbox("Use regex")
	chkWholeWord.OnToggled(func(*ui.Checkbox) {
		chkRegex.SetChecked(false)
	})
	chkRegex.OnToggled(func(*ui.Checkbox) {
		chkWholeWord.SetChecked(false)
	})
	hbox.Append(chkWholeWord, true)
	hbox.Append(chkCaseSensitive, true)
	hbox.Append(chkRegex, true)

	// Add entry control
	entry := ui.NewEntry()
	vbox.Append(entry, true)

	// Make a grid to set position of buttons
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	// Implement button "Yes"
	okButton := ui.NewButton(yes)
	okButton.OnClicked(func(*ui.Button) {
		toSearch := entry.Text() // Catch entry before destroying window ...
		cs := chkCaseSensitive.Checked()
		ww := chkWholeWord.Checked()
		rx := chkRegex.Checked()
		dialogWindow.Destroy()
		fyes(toSearch, cs, ww, rx)
		mainwin.Enable()
	})
	// Implement button "No"
	noButton := ui.NewButton(no)
	noButton.OnClicked(func(*ui.Button) {
		dialogWindow.Destroy()
		fno()
		mainwin.Enable()
	})
	// Add buttons
	grid.Append(ui.NewLabel(""), // Fake controle to align buttons at right end.
		1, 2, 1, 1,
		true, ui.AlignFill, false, ui.AlignEnd)
	grid.Append(okButton,
		2, 2, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	grid.Append(noButton,
		3, 2, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	// Let's show result
	dialogWindow.SetChild(vbox)
	dialogWindow.SetMargined(true)
	dialogWindow.Show()
}
