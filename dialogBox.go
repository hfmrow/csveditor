// dialogBox.go

package main

import (
	"github.com/andlabs/ui"
)

// Make dialog box with callback functions to allow or deny action. fyes give entry.text value
func DialogBoxEntry(mainwin *ui.Window, title, yes, no string, fyes func(entry string), fno func(), text ...string) {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	dialogWindow := ui.NewWindow(title, 10, 10, false)
	dialogWindow.OnClosing(func(*ui.Window) bool {
		//	mainwin.Enable() // Make MAINWIN enaled
		return true
	})
	ui.OnShouldQuit(func() bool {
		dialogWindow.Destroy()
		//	mainwin.Enable() // Make MAINWIN enaled
		return true
	})
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewLabel(""), true) // Fake control, used as separator.
	for idx, txt := range text {
		if idx == len(text)-1 {
			vbox.Append(ui.NewLabel(txt+"\n"), false) // if its the last, add LF to make separation.
		} else {
			vbox.Append(ui.NewLabel(txt), false) // Label with txt
		}
	}

	// Add entry control
	entry := ui.NewEntry()
	vbox.Append(entry, true)

	// Make a grid to set position of buttons
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	// Implement Yes button
	okButton := ui.NewButton(yes)
	okButton.OnClicked(func(*ui.Button) {
		toSearch := entry.Text() // Catch entry before destroying window ...
		dialogWindow.Destroy()
		fyes(toSearch)
	})
	// Implement No button
	noButton := ui.NewButton(no)
	noButton.OnClicked(func(*ui.Button) {
		dialogWindow.Destroy()
		fno()
	})
	// Add bottom buttons
	grid.Append(ui.NewLabel(""), // Fake controle to aligne buttons at right end.
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

// Make dialog box with callback functions to allow or deny action
func DialogBoxAsk(mainwin *ui.Window, title, yes, no string, fyes, fno func(), text ...string) {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	dialogWindow := ui.NewWindow(title, 10, 10, false)
	dialogWindow.OnClosing(func(*ui.Window) bool {
		return true
	})
	ui.OnShouldQuit(func() bool {
		dialogWindow.Destroy()
		return true
	})
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewLabel(""), true) // Fake control, used as separator.
	for idx, txt := range text {
		if idx == len(text)-1 {
			vbox.Append(ui.NewLabel(txt+"\n"), false) // if its the last, add LF to make separation.
		} else {
			vbox.Append(ui.NewLabel(txt), false) // Label with txt
		}
	}
	// Make a grid to set position of buttons
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)
	// Implement Yes button
	okButton := ui.NewButton(yes)
	okButton.OnClicked(func(*ui.Button) {
		dialogWindow.Destroy()
		fyes()
	})
	// Implement No button
	noButton := ui.NewButton(no)
	noButton.OnClicked(func(*ui.Button) {
		dialogWindow.Destroy()
		fno()
	})
	// Add bottom buttons
	grid.Append(ui.NewLabel(""), // Fake controle to aligne buttons at right end.
		1, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignEnd)
	grid.Append(okButton,
		2, 1, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	grid.Append(noButton,
		3, 1, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	// Let's show result
	dialogWindow.SetChild(vbox)
	dialogWindow.SetMargined(true)
	dialogWindow.Show()
}

// Make dialog box with callback functions to allow or deny action
func DialogBoxAsk3(mainwin *ui.Window, title, yes, no, cancel string, fyes, fno, fcancel func(), text ...string) {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	dialogWindow := ui.NewWindow(title, 10, 10, false)
	dialogWindow.OnClosing(func(*ui.Window) bool {
		return true
	})
	ui.OnShouldQuit(func() bool {
		dialogWindow.Destroy()
		return true
	})
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewLabel(""), true) // Fake control, used as separator.
	for idx, txt := range text {
		if idx == len(text)-1 {
			vbox.Append(ui.NewLabel(txt+"\n"), false) // if its the last, add LF to make separation.
		} else {
			vbox.Append(ui.NewLabel(txt), false) // Label with txt
		}
	}
	// Make a grid to set position of buttons
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)
	// Implement Yes button
	okButton := ui.NewButton(yes)
	okButton.OnClicked(func(*ui.Button) {
		dialogWindow.Destroy()
		fyes()
	})
	// Implement No button
	noButton := ui.NewButton(no)
	noButton.OnClicked(func(*ui.Button) {
		dialogWindow.Destroy()
		fno()
	})
	// Implement Cancel button
	cancelButton := ui.NewButton(cancel)
	cancelButton.OnClicked(func(*ui.Button) {
		dialogWindow.Destroy()
		fcancel()
	})
	// Add bottom buttons
	grid.Append(ui.NewLabel(""), // Fake controle to align buttons at bottom right.
		0, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignEnd)
	grid.Append(cancelButton,
		1, 0, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	grid.Append(okButton,
		2, 0, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	grid.Append(noButton,
		3, 0, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	// Let's show result
	dialogWindow.SetChild(vbox)
	dialogWindow.SetMargined(true)
	dialogWindow.Show()
}

func DialogBoxMsg(mainwin *ui.Window, title, text string) {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	ui.MsgBox(mainwin,
		title,
		text)
}

func DialogBoxErr(mainwin *ui.Window, title, text string) {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	ui.MsgBoxError(mainwin,
		title,
		text)
}

func DialogBoxSve(mainwin *ui.Window) string {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	return ui.SaveFile(mainwin)
}

func DialogBoxOpn(mainwin *ui.Window) string {
	if mainwin == nil {
		panic("mainwin not present !!!")
	}
	return ui.OpenFile(mainwin)
}
