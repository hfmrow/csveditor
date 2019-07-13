// newTableWindow.go

package main

import (
	"encoding/csv"

	"io/ioutil"
	"strings"

	"github.com/andlabs/ui"
	"github.com/hfmrow/csveditor/genLib"
	"golang.org/x/net/html/charset"
)

func newTableWindow() {
	mainwin.Disable() // Make MAINWIN disabled

	editRowWin := ui.NewWindow("New CSV file", 640, 10, false)
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
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	group := ui.NewGroup("")

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	// Entry for columns list
	entryBox := ui.NewEntry()
	entryForm.Append("", ui.NewLabel(`Enter column names (i.e: column1, "column 2", column3 ...)`), false)
	entryForm.Append("", entryBox, false)

	group.SetChild(entryForm)
	vbox.Append(group, false)

	// Make a grid to set position of buttons
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	// OK button ****
	okButton := ui.NewButton("Ok")
	okButton.OnClicked(func(*ui.Button) {
		fieldNames := entryBox.Text()
		if fieldNames != "" {
			// Get charset and convert to utf-8 if needed
			charSET := genLib.DetectCharsetStr(fieldNames)
			if charSET != "utf-8" {
				r, _ := charset.NewReader(strings.NewReader(fieldNames), charSET) // convert to UTF-8
				textFileBytes, _ := ioutil.ReadAll(r)
				fieldNames = string(textFileBytes)
			}
			// Read as csv, this will detect it's the right csv format
			csvReader := csv.NewReader(strings.NewReader(fieldNames))
			csvReader.Comma = ','
			csvReader.FieldsPerRecord = len(strings.Split(fieldNames, ","))
			records, err := csvReader.ReadAll()
			if err != nil {
				DialogBoxErr(mainwin, "Error !", "Csv creation error, i.e: column1, "+`"column 2", column3 ...`+"\nTake care about multi-words name, they need to be double-quoted")
			} else {

				filename := DialogBoxSve(mainwin)
				if filename != "" {
// fixing issue when new csv file created. (
					if len(records) < 2 {
						var tmpRow []string
						for idx := len(records[0]); idx > 0; idx-- {
							tmpRow = append(tmpRow, " ")
						}
						records = append(records, tmpRow)
					}

					genLib.WriteCsv(filename, ",", records)
					// Remove previous opt file if exist
					optFilename := getOptFileName(filename)
					genLib.RemovIfExist(optFilename)
					// Initialising storeProfiles
					storeProfiles = NewStoreProfiles()
					// Load, Store and set CSV
					addAndSetCsv(filename, "Main")
					ChkRow.Reset()         // Reset stored checked rows history
					sortFieldsFlag = false // Reset sort state
					// Refresh tabs
					refreshTable()
					refreshOption()
					refreshSort()
				}
			}
			mainwin.Enable()     // Make MAINWIN enaled
			editRowWin.Destroy() // Closing edit window
		}
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
