// editRowWindow.go

package main

import (
	"fmt"

	"github.com/andlabs/ui"
)

func editRowWindow() {
	mainwin.Disable() // Make MAINWIN disabled
	var flagEdit, dupFlag bool
	var title string
	currentRowValues := make([]string, CsvProfileList.NumberCols)
	if ChkRow.NotEmpty() {
		currentRowValues = tableDatas[ChkRow.Last()] // Get entry for last checked row.
		flagEdit = true
		title = fmt.Sprintf("Entry editor, row :  %d", ChkRow.Last()+1)
	} else {
		title = "Entry editor :  New entry"
	}

	editRowWin := ui.NewWindow(title, 640, 10, false)
	editRowWin.OnClosing(func(*ui.Window) bool {
		mainwin.Enable() // Make MAINWIN enaled
		return true
	})
	ui.OnShouldQuit(func() bool {
		editRowWin.Destroy()
		return true
	})
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	// hbox := ui.NewHorizontalBox()
	// hbox.SetPadded(true)
	group := ui.NewGroup("")
	entryForm := ui.NewForm()
	entryForm.SetPadded(true)

	// Make spinBox (for insert) if it's a new row
	spinBox := ui.NewSpinbox(1, CsvProfileList.NumberRows+1)
	spinBox.SetValue(1)

	// CheckBox for duplicate entry when edited
	dupCheckBox := ui.NewCheckbox("Duplicate entry")
	dupCheckBox.OnToggled(func(*ui.Checkbox) { dupFlag = !dupFlag })
	if !flagEdit { // New entry
		entryForm.Append("", ui.NewLabel("Insert at:"), false)
		entryForm.Append("", spinBox, false)
	} else { // Edit entry
		entryForm.Append("", dupCheckBox, false)
	}

	tmpStoreEntry := make([]ui.Entry, CsvProfileList.NumberCols) // Array of ui.Entry to store struct
	for idx, colName := range CsvProfileList.FieldNames {
		tmpStoreEntry[idx] = *ui.NewEntry()
		tmpStoreEntry[idx].SetText(currentRowValues[idx])
		if CsvProfileList.FieldDisplay[idx] == getValChar("y", 0).(string) { // "Yes"to display this column
			entryForm.Append("", ui.NewLabel(colName), false) // Label with column name
			entryForm.Append("", &tmpStoreEntry[idx], false)  // Text entry
		}
	}

	group.SetChild(entryForm)
	vbox.Append(group, false)

	// Make a grid to set position of buttons
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	// OK button ****
	okButton := ui.NewButton("Ok")
	okButton.OnClicked(func(*ui.Button) {
		var verifModif string                                 // Used to check if row entry isn't empty
		tmpEntry := make([]string, CsvProfileList.NumberCols) // create new empty row with all columns
		for idx := 0; idx < len(tmpStoreEntry); idx++ {       // Check if row entry is empty or not

			tmpEntry[idx] = tmpStoreEntry[idx].Text()
			verifModif += tmpEntry[idx]

		}
		if verifModif != "" {
			tableAddEditRow(spinBox.Value(), flagEdit, dupFlag, tmpEntry)
		}
		mainwin.Enable()     // Make MAINWIN enaled
		editRowWin.Destroy() // Closing edit window
	})

	// CANCEL button ****
	noButton := ui.NewButton("Cancel")
	noButton.OnClicked(func(*ui.Button) {
		mainwin.Enable()     // Make MAINWIN enaled
		editRowWin.Destroy() // Closing edit window
	})

	grid.Append(ui.NewLabel(""), // Fake controle to align buttons at bottom right.
		0, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignEnd)
	grid.Append(okButton,
		1, 0, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)
	grid.Append(noButton,
		2, 0, 1, 1,
		false, ui.AlignEnd, false, ui.AlignEnd)

	// Display Edit window
	editRowWin.SetChild(vbox)
	editRowWin.SetMargined(true)
	editRowWin.Show()
}
